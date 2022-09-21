package cloud

import (
	"github.com/sachaservan/bgn"
	"github.com/sachaservan/paillier"
	"math/big"
	"math/rand"
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

func computeInnerProd(encTrans [][] *bgn.Ciphertext, encQuery [] *bgn.Ciphertext, pk *bgn.PublicKey) [] *bgn.Ciphertext {
	numTrans := len(encTrans)
	prods := make([] *bgn.Ciphertext, numTrans)
	for i := 0; i < numTrans; i++ {
		prods[i] = innerProCC(encTrans[i], encQuery, pk)
	}
	return prods
}

// compute [start, end) transactions in encTrans
func asynComputeInnerProdRange(encTrans [][] *bgn.Ciphertext, start int, end int, encQuery [] *bgn.Ciphertext, pk *bgn.PublicKey, callback Callback) {
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
		prods = append(prods, innerProCC(encTrans[i], encQuery, pk))
	}
	callback(prods)
}

func innerProCC(encVec1, encVec2 [] *bgn.Ciphertext, pk *bgn.PublicKey) *bgn.Ciphertext {
    if len(encVec1) != len(encVec2) {
        panic("length of vectors is not same")
    }
	ptxtZero := big.NewInt(0)
	temp := pk.Encrypt(ptxtZero)

	for i := 0; i < len(encVec1); i++ {
		prod := pk.Mult(encVec1[i], encVec2[i])
		temp = pk.Add(prod, temp)
	}
	return temp
}


func ParallelCompute(encTrans [][] *bgn.Ciphertext, encQuery [] *bgn.Ciphertext, negL1Norm *bgn.Ciphertext,
	pk *bgn.PublicKey, k, n, numCPU int) [] *bgn.Ciphertext {
		runtime.GOMAXPROCS(numCPU + 1)
		dummyTrans := makeDummyTrans(k, n, pk)
		mergedTrans := append(encTrans, dummyTrans...)
		numTrans := len(mergedTrans)
		numTransPerCPU := numTrans / numCPU
		remainder := numTrans % numCPU

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

		// Parallelism
		var mergedProds [] *bgn.Ciphertext
		var wg sync.WaitGroup
		wg.Add(numCPU)
		for i := range starts {
			f := func(subProds [] *bgn.Ciphertext) {
				mergedProds = append(mergedProds, subProds...)
				wg.Done()
			}
			go asynComputeInnerProdRange(mergedTrans, starts[i], ends[i], encQuery, pk, f)
		}
		wg.Wait()

		var randomizedSet [] *bgn.Ciphertext
		for i := range mergedProds {
			sum := pk.Add(mergedProds[i], negL1Norm)
			randCoef := rand.Int() + 1
			w := pk.MultConst(sum, big.NewInt(int64(randCoef)))
			randomizedSet = append(randomizedSet, w)
		}

		randomizedSet = shuffle(randomizedSet)
		return randomizedSet
	}

func Compute(encTrans [][] *bgn.Ciphertext, encQuery [] *bgn.Ciphertext, negL1Norm *bgn.Ciphertext,
	pk *bgn.PublicKey, k, n int) [] *bgn.Ciphertext {

	transProds := computeInnerProd(encTrans, encQuery, pk)
	dummyProds := computeInnerProd(makeDummyTrans(k, n, pk), encQuery, pk)
	mergedProds := append(transProds, dummyProds...)
	var randomizedSet [] *bgn.Ciphertext
	for i := range mergedProds {
		sum := pk.Add(mergedProds[i], negL1Norm)
		randCoef := rand.Int() + 1
		w := pk.MultConst(sum, big.NewInt(int64(randCoef)))
		randomizedSet = append(randomizedSet, w)
	}

	randomizedSet = shuffle(randomizedSet)
	return randomizedSet
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
