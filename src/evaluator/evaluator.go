package evaluator

import (
	"github.com/sachaservan/bgn"
	"github.com/sachaservan/paillier"
	"math/big"
	"math"
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

	genG1 := pkBGN.P.NewFieldElement()
	genG1.PowBig(pkBGN.P, skBGN.Key)

	genGT := pkBGN.Pairing.NewGT().Pair(pkBGN.P, pkBGN.P)
	genGT.PowBig(genGT, skBGN.Key)
	pkBGN.PrecomputeTables(genG1, genGT)

	// Generate Paillier key pair
	skPaillier, pkPaillier = paillier.KeyGen(160)
	return
}

func SIMDEvaluate(randomizedSet [] *bgn.Ciphertext, pkBGN *bgn.PublicKey, skBGN *bgn.SecretKey,
	pkPaillier *paillier.PublicKey, originM, originK, numPacking, numInterval int) []*paillier.Ciphertext {
	randomizedSetNum := len(randomizedSet)
	// v := make([]*paillier.Ciphertext, originM)
	var v []*paillier.Ciphertext
	k10PowInterval := big.NewInt(int64(math.Pow(10, float64(numInterval)))) // 10 ^ numInterval

	for i := 0; i < randomizedSetNum; i++ {
        // println("Decrypt", i, "-th bgn ctxt")
		packedPtxt, err := skBGN.Decrypt(randomizedSet[i], pkBGN)
		if (err != nil) {
			panic("Error: BGN deryption failed." + err.Error())
		}
		unpackedPtxts := make([] *big.Int, numPacking)
		for j := 0; j < numPacking; j ++ {
			ten := big.NewInt(int64(10))
			unpackedPtxts[j] = packedPtxt.Mod(packedPtxt, ten)
			if (j != numPacking - 1) {
				packedPtxt = packedPtxt.Div(packedPtxt, k10PowInterval)
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
		ptxt, err := skBGN.Decrypt(randomizedSet[i], pkBGN)
		if (err != nil) {
			panic("Error: BGN deryption failed." + err.Error())
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
