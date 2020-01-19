package evaluator

import (
	"github.com/sachaservan/bgn"
	"github.com/sachaservan/paillier"
	"math/big"
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
	skPaillier, pkPaillier = paillier.CreateKeyPair(160)
	return
}

func Evaluate(randomizedSet [] *bgn.Ciphertext, pkBGN *bgn.PublicKey, skBGN *bgn.SecretKey,
	pkPaillier *paillier.PublicKey) []*paillier.Ciphertext {

	transNum := len(randomizedSet)
	v := make([]*paillier.Ciphertext, transNum)
	for i := 0; i < transNum; i++ {
<<<<<<< HEAD
        // println("Decrypt", i, "-th bgn ctxt")
=======
>>>>>>> 776741279ca80cebcd71cf7a6900206b7bfa3665
		ptxt := skBGN.Decrypt(randomizedSet[i], pkBGN)
		if ptxt.String() == "0" {
			v[i] = pkPaillier.Encrypt(big.NewInt(1))
		} else {
			v[i] = pkPaillier.Encrypt(big.NewInt(0))
		}
	}
	return v
}

func Compare(maskedFlag *paillier.Ciphertext, sk *paillier.SecretKey) int {
	z := sk.Decrypt(maskedFlag)
	halfN := big.NewInt(0)
	halfN.Div(sk.N, big.NewInt(2))
	return z.Cmp(halfN)
}
