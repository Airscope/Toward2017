package utils

import (
	"github.com/sachaservan/bgn"
	"math/big"
)

func BGNPlaintxt(pk *bgn.PublicKey, num int) * bgn.Plaintext{
	return pk.NewPlaintext(big.NewFloat(float64(num)))
}

func BGNPlaintxtOne(pk *bgn.PublicKey) * bgn.Plaintext {
	return BGNPlaintxt(pk, 1)
}

func BGNPlaintxtZero(pk *bgn.PublicKey) * bgn.Plaintext {
	return BGNPlaintxt(pk, 0)
}