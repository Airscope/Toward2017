package evaluator

import (
	"github.com/Airscope/Toward2017/utils"
	"github.com/sachaservan/bgn"
	"github.com/sachaservan/paillier"
	"math/big"
	"math"
	"fmt"
	gmp "github.com/ncw/gmp"
)

// Step 1: System Setup
func SystemSetup() (pkBGN *bgn.PublicKey, skBGN *bgn.SecretKey, pkPaillier *paillier.PublicKey, skPaillier *paillier.SecretKey) {

	// Generate BGN key pair
	keyBits := 512 // length of q1 and q2
	messageSpace := big.NewInt(1021)
	polyBase := 3 // base for the ciphertext polynomial
	fpScaleBase := 3
	fpPrecision := 0.0001

	pkBGN, skBGN, _ = bgn.NewKeyGen(keyBits, messageSpace, polyBase, fpScaleBase, fpPrecision, true)

	// Generate Paillier key pair
	skPaillier, pkPaillier = paillier.KeyGen(160)
	return
}

func SIMDEvaluate(randomizedSet [] *bgn.Ciphertext, pkBGN *bgn.PublicKey, skBGN *bgn.SecretKey,
	pkPaillier *paillier.PublicKey, originM, originK, numPacking, numInterval int) []*paillier.Ciphertext {
	randomizedSetNum := len(randomizedSet)
	// v := make([]*paillier.Ciphertext, originM)
	var v []*paillier.Ciphertext
	k10PowInterval := int(math.Pow(10, float64(numInterval)))// 10 ^ numInterval
	for i := 0; i < randomizedSetNum; i++ {
        // println("Decrypt", i, "-th bgn ctxt")
		packedPtxt := skBGN.Decrypt(randomizedSet[i], pkBGN)
		if (packedPtxt == nil) {
			panic("Error: BGN deryption failed.")
		}
		packedPtxtInt := 0
		_, err := fmt.Sscan(packedPtxt.String(), &packedPtxtInt)
		if (err != nil) {
			println("strings Atoi error, ", err.Error())
		}
		unpackedPtxts := make([] *bgn.Plaintext, numPacking)
		for j := 0; j < numPacking; j ++ {
			unpackedPtxts[j] = utils.BGNPlaintxt(pkBGN, packedPtxtInt % 10)
			if (j != numPacking - 1) {
				packedPtxtInt = packedPtxtInt / k10PowInterval
			}
		}
		for j:=0; j < numPacking; j++ {
			if (unpackedPtxts[j].String() == "0") {
				v = append(v, pkPaillier.Encrypt(gmp.NewInt(1)))
			} else {
				v = append(v, pkPaillier.Encrypt(gmp.NewInt(0)))
			}
		}
	}
	if (len(v) != originM + originK) {
		println( "len(v)=",len(v), ", M=", originM, ", K=", originK)
		panic("The length of v does not equal to M+K!")
	}

	return v
}

func Evaluate(randomizedSet [] *bgn.Ciphertext, pkBGN *bgn.PublicKey, skBGN *bgn.SecretKey,
	pkPaillier *paillier.PublicKey) []*paillier.Ciphertext {

	transNum := len(randomizedSet)
	v := make([]*paillier.Ciphertext, transNum)
	for i := 0; i < transNum; i++ {
        // println("Decrypt", i, "-th bgn ctxt")
		ptxt := skBGN.Decrypt(randomizedSet[i], pkBGN)
		if (ptxt == nil) {
			panic("Error: BGN deryption failed.")
		}
		if ptxt.String() == "0" {
			v[i] = pkPaillier.Encrypt(gmp.NewInt(1))
		} else {
			v[i] = pkPaillier.Encrypt(gmp.NewInt(0))
		}
	}
	return v
}

func Compare(maskedFlag *paillier.Ciphertext, sk *paillier.SecretKey) int {
	z:= sk.Decrypt(maskedFlag)
	halfN := gmp.NewInt(0)
	halfN.Div(sk.N, gmp.NewInt(2))
	return z.Cmp(halfN)
}
