package cloud

import (
	"github.com/sachaservan/bgn"
	"github.com/sachaservan/paillier"
	"math/big"
	"math/rand"
)

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
			ptxt := float64(rand.Intn(2))
			tmp[j] = pk.Encrypt(pk.NewPlaintext(big.NewFloat(ptxt)))
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

func innerProCC(encVec1, encVec2 [] *bgn.Ciphertext, pk *bgn.PublicKey) *bgn.Ciphertext {
    if len(encVec1) != len(encVec2) {
        panic("length of vectors is not same")
    }
	ptxtZero := pk.NewPlaintext(big.NewFloat(0.0))
	temp := pk.Encrypt(ptxtZero)

	for i := 0; i < len(encVec1); i++ {
		prod := pk.EMult(encVec1[i], encVec2[i])
		temp = pk.EAdd(prod, temp)
	}
	return temp
}

func Compute(encTrans [][] *bgn.Ciphertext, encQuery [] *bgn.Ciphertext, negL1Norm *bgn.Ciphertext,
	pk *bgn.PublicKey, k, n int) [] *bgn.Ciphertext {

	transProds := computeInnerProd(encTrans, encQuery, pk)
	dummyProds := computeInnerProd(makeDummyTrans(k, n, pk), encQuery, pk)
	mergedProds := append(transProds, dummyProds...)
	var randomizedSet [] *bgn.Ciphertext
	for i := range mergedProds {
		sum := pk.EAdd(mergedProds[i], negL1Norm)
		randCoef := rand.Int() + 1
		w := pk.EMultC(sum, big.NewFloat(float64(randCoef)))
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
	sum := pk.Encrypt(big.NewInt(0))
	for i := 0; i < m; i++ {
		sum = pk.EAdd(sum, originV[i])
	}
	return sum
}

func Compare(supp, minSupp *paillier.Ciphertext, pk *paillier.PublicKey) *paillier.Ciphertext {
	diff := pk.ESub(supp, minSupp)
	tmp := pk.ECMult(diff, big.NewInt(rand.Int63n(1024)))
	return tmp
}
