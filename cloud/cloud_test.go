package cloud

import (
	"github.com/Airscope/evaluator"
	"github.com/sachaservan/bgn"
	"github.com/sachaservan/paillier"
	"math/big"
	"math/rand"
	"testing"
)

func TestCompare(t *testing.T) {
	skPaillier, pkPaillier := paillier.CreateKeyPair(128)
	for i := 0; i < 100; i++ {
		msg1 := rand.Int63n(1000)
		msg2 := rand.Int63n(1000)
		ctxt1 := pkPaillier.Encrypt(big.NewInt(msg1))
		ctxt2 := pkPaillier.Encrypt(big.NewInt(msg2))
		maskedFlag := Compare(ctxt1, ctxt2, pkPaillier)
		output := evaluator.Compare(maskedFlag, skPaillier)
		cmp := msg1 < msg2
		cmp2 := output == 1
		if cmp != cmp2 {
			t.Errorf("wrong at compare %d and %d", msg1, msg2)
		}
	}
}

func TestInnerProCC(t *testing.T) {
	cases := []struct {
		in1 [] float64
		in2 [] float64
		want string
	}{
		{[]float64{1.0, 2.0, 3.0}, []float64{1.0, 2.0, 3.0}, "14"},
		{[]float64{1.0, 0.0, 1.0, 0.0, 1.0}, []float64{0.0, 1.0, 1.0, 0.0, 1.0}, "2"},
		{[]float64{11.0, 22.0, 33.0, 44.0}, []float64{55.0, 66.0, 77.0, 88.0}, "8470"},
	}
	for _, c:= range cases {
		got := InnerProCCForTesting(c.in1, c.in2)
		if got != c.want {
			t.Errorf("InnerProCC == E[%q], want %q", got, c.want)
		}
	}
}

func InnerProCCForTesting(vec1, vec2 [] float64) string {
	keyBits := 512 // length of q1 and q2
	messageSpace := big.NewInt(1021)
	polyBase := 3 // base for the ciphertext polynomial
	fpScaleBase := 3
	fpPrecision := 0.0001

	pk, sk, _ := bgn.NewKeyGen(keyBits, messageSpace, polyBase, fpScaleBase, fpPrecision, true)

	genG1 := pk.P.NewFieldElement()
	genG1.PowBig(pk.P, sk.Key)

	genGT := pk.Pairing.NewGT().Pair(pk.P, pk.P)
	genGT.PowBig(genGT, sk.Key)
	pk.PrecomputeTables(genG1, genGT)

	var encVec1 [] *bgn.Ciphertext
	var encVec2 [] *bgn.Ciphertext

	for i:= 0; i < len(vec1) ; i++ {
		m1 := pk.NewPlaintext(big.NewFloat(vec1[i]))
		m2 := pk.NewPlaintext(big.NewFloat(vec2[i]))
		c1 := pk.Encrypt(m1)
		c2 := pk.Encrypt(m2)
		encVec1 = append(encVec1, c1)
		encVec2 = append(encVec2, c2)
	}

	prod := innerProCC(encVec1, encVec2, pk)
	ptxt := sk.Decrypt(prod, pk).String()
	// val, err = strconv.ParseFloat(ptxt, 64)
	return ptxt
}


