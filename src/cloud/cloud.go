package cloud

import (
	// "fmt"
	"github.com/sachaservan/bgn"
	"github.com/sachaservan/paillier"
	"math/big"
	"math/rand"
	"math"
	"sync"
	"runtime"
	gmp "github.com/ncw/gmp"
)

type Callback func(ciphertexts [] *bgn.Ciphertext)

var shuffleMap [] int

func shuffle(encTrans [] *bgn.Ciphertext) [] *bgn.Ciphertext {
	transNum := len(encTrans)
	shuffleMap = make([]int, transNum)
	for i := 0; i < transNum; i++ {
		shuffleMap[i] = i
	}

	for i := 0; i < transNum; i++ {
		randIndex := rand.Intn(transNum)
		tmp := shuffleMap[i]
		shuffleMap[i] = shuffleMap[randIndex]
		shuffleMap[randIndex] = tmp
		tmpCtxt := encTrans[i]
		encTrans[i] = encTrans[randIndex]
		encTrans[randIndex] = tmpCtxt
	}
	return encTrans
}

func invShuffle(encTrans [] *paillier.Ciphertext) [] *paillier.Ciphertext {
	transNum := len(encTrans)
	duplicate := make([]*paillier.Ciphertext, transNum)
	copy(duplicate, encTrans)

	for i := 0; i < transNum; i++ {
		encTrans[shuffleMap[i]] = duplicate[i]
	}
	return encTrans
}

func makeDummyTrans(k, n int, pk *bgn.PublicKey) [][] *bgn.Ciphertext {
	var dummyTrans [][] *bgn.Ciphertext
	for i := 0; i < k; i++ {
		tmp := make([] *bgn.Ciphertext, n)
		for j := 0; j < n; j++ {
			ptxt := int64(rand.Intn(2))
			tmp[j] = pk.Encrypt(big.NewInt(ptxt))
		}
		dummyTrans = append(dummyTrans, tmp)
	}
	return dummyTrans
}

func computeInnerProd(encTrans [][] *bgn.Ciphertext, encQuery [] *bgn.Ciphertext, pk *bgn.PublicKey, sk *bgn.SecretKey) [] *bgn.Ciphertext {
	numTrans := len(encTrans)
	prods := make([] *bgn.Ciphertext, numTrans)
	for i := 0; i < numTrans; i++ {
		prods[i] = innerProCC(encTrans[i], encQuery, pk, sk)
	}
	return prods
}

// compute [start, end) transactions in encTrans
func asynComputeInnerProdRange(encTrans [][] *bgn.Ciphertext, start int, end int, encQuery [] *bgn.Ciphertext, pk *bgn.PublicKey, sk *bgn.SecretKey, callback Callback) {
	numTrans := len(encTrans)
	if (start < 0) {
		return;
	}
	if (end <= start) {
		return;
	}
	if end > numTrans {
		return;
	}
	var prods [] *bgn.Ciphertext
	for i := start; i < end; i++ {
		prods = append(prods, innerProCC(encTrans[i], encQuery, pk, sk))
	}
	callback(prods)
}

func innerProCC(encVec1, encVec2 [] *bgn.Ciphertext, pk *bgn.PublicKey, sk *bgn.SecretKey) *bgn.Ciphertext {
    if len(encVec1) != len(encVec2) {
        panic("Lengths of encrypted vectors are not the same")
    }
	ptxtZero := big.NewInt(0)
	temp := pk.Encrypt(ptxtZero)

	for i := 0; i < len(encVec1); i++ {
		prod := pk.Mult(encVec1[i], encVec2[i])
		temp = pk.Add(prod, temp)
		// _, err := sk.Decrypt(temp, pk)
		// if (err != nil) {
		// 	panic("Error: BGN deryption failed." + err.Error())
		// }
	}
	return temp
}

func ParallelSIMDCompute(encTrans [][] *bgn.Ciphertext, encQuery [] *bgn.Ciphertext, negL1Norm *bgn.Ciphertext,
	pk *bgn.PublicKey, sk *bgn.SecretKey, k, n, numCPU, numPacking, numInterval int) [] *bgn.Ciphertext {
		runtime.GOMAXPROCS(numCPU + 1)
		dummyTrans := makeDummyTrans(k, n, pk)
		mergedTrans := append(encTrans, dummyTrans...)
		numTrans := len(mergedTrans)

		/** SIMD **/

		// Packing
		numPackedTrans := numTrans / numPacking
		packedTrans := make([][] *bgn.Ciphertext, numPackedTrans)
		for i := 0; i < numPackedTrans; i ++ {
			packedTranRow := make([] *bgn.Ciphertext, n) // 1行packed的事务，长度为n，包含n个packed的数据
			for j := 0; j < n; j++ {
				packedCtxt := pk.Encrypt(big.NewInt(0))
				k10PowInterval := big.NewInt(int64(math.Pow(10, float64(numInterval)))) // 10 ^ numInterval
				for k := 0; k < numPacking; k++ {
					packedCtxt = pk.Add(packedCtxt, mergedTrans[i*numPacking+k][j])
					if (k != numPacking - 1) {
						packedCtxt = pk.MultConst(packedCtxt, k10PowInterval)
					}
				}
				packedTranRow[j] = packedCtxt
			 }
			packedTrans[i] = packedTranRow
		}

		/** Parallelism **/
		numTransPerCPU := numPackedTrans / numCPU
		remainder := numPackedTrans % numCPU

		// 为使划分均匀，其中remainder个CPU持有numTransPerCPU+1个trans
		// numCPU-remainder个CPU持有numTransPerCPU个trans
		var starts [] int;
		var ends [] int;
		for i := 0; i < remainder; i++ {
			starts = append(starts, i*(numTransPerCPU+1))
			ends = append(ends, (i+1)*(numTransPerCPU+1))
		}
		for i := 0; i < numCPU-remainder; i++ {
			offset := remainder * (numTransPerCPU+1)
			starts = append(starts, offset + i * numTransPerCPU)
			ends = append(ends, offset + (i+1)*numTransPerCPU)
		}
		if len(starts) != len(ends) {
			panic("The number of array \"starts\" != \"end\" !")
		}

		// Multi-cpus
		var packedProds [] *bgn.Ciphertext
		var wg sync.WaitGroup
		wg.Add(numCPU)
		for i := range starts {
			callback := func(subProds [] *bgn.Ciphertext) {
				packedProds = append(packedProds, subProds...)
				wg.Done()
			}
			go asynComputeInnerProdRange(packedTrans, starts[i], ends[i], encQuery, pk, sk, callback)
		}
		wg.Wait()

		/** Randomization **/
		var randomizedSet [] *bgn.Ciphertext
		for i := range packedProds {
			sum := pk.Add(packedProds[i], negL1Norm)
			randCoef := rand.Intn(8) + 1
			w := pk.MultConst(sum, big.NewInt(int64(randCoef)))
			randomizedSet = append(randomizedSet, w)
		}
		// CAUTION: DO NOT Randomization-shuffling, only DO multiply random coef.
		return randomizedSet
	}

func Compute(encTrans [][] *bgn.Ciphertext, encQuery [] *bgn.Ciphertext, negL1Norm *bgn.Ciphertext,
	pk *bgn.PublicKey, sk *bgn.SecretKey, k, n int) [] *bgn.Ciphertext {

	transProds := computeInnerProd(encTrans, encQuery, pk, sk)
	// _, err := sk.Decrypt(transProds[0], pk)
	// if (err != nil) {
	// 	panic("Error: BGN deryption failed." + err.Error())
	// }

	dummyProds := computeInnerProd(makeDummyTrans(k, n, pk), encQuery, pk, sk)
	mergedProds := append(transProds, dummyProds...)
	var randomizedSet [] *bgn.Ciphertext
	for i := range mergedProds {
		sum := pk.Add(mergedProds[i], negL1Norm)
		randCoef := rand.Intn(16) + 1
		w := pk.MultConst(sum, big.NewInt(int64(randCoef)))
		randomizedSet = append(randomizedSet, w)
	}

	randomizedSet = shuffle(randomizedSet)
	return randomizedSet
}

func ParallelSIMDEvaluate(v [] *paillier.Ciphertext, pk *paillier.PublicKey, m, k int) *paillier.Ciphertext {
	if len(v) != m+k {
		panic("The number of transactions in v's != m+k !")
	}
	// DO NOT need to invShuffling
	sum := pk.Encrypt(gmp.NewInt(0))
	for i := 0; i < m; i++ {
		sum = pk.Add(sum, v[i])
	}
	return sum
}

func Evaluate(v [] *paillier.Ciphertext, pk *paillier.PublicKey, m, k int) *paillier.Ciphertext {
	if len(v) != m+k {
		panic("The number of transactions in v's != m+k !")
	}
	originV := invShuffle(v)[:m]
	sum := pk.Encrypt(gmp.NewInt(0))
	for i := 0; i < m; i++ {
		sum = pk.Add(sum, originV[i])
	}
	return sum
}

func Compare(supp, minSupp *paillier.Ciphertext, pk *paillier.PublicKey) *paillier.Ciphertext {
	diff := pk.Sub(supp, minSupp)
	tmp := pk.ConstMult(diff, gmp.NewInt(int64(rand.Int63n(1024))))
	return tmp
}
