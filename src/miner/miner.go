package miner

import (
	"github.com/sachaservan/bgn"
	"math"
	"math/big"
)

func Compute(pk *bgn.PublicKey, n int) (encQuery [] *bgn.Ciphertext, negL1Norm *bgn.Ciphertext) {
	ptxtQuery := makeQuery(n)
	sum := 0.0
	for i:=0; i<n;i++  {
		sum += math.Abs(ptxtQuery[i])
		encQuery = append(encQuery, pk.Encrypt(big.NewInt(int64(ptxtQuery[i]))))
	}
	negL1Norm = pk.Encrypt((big.NewInt(int64(-sum))))
	return
}

func makeQuery(n int) [] float64 {
	query := make([] float64, n)
	for i := 0; i < n; i++ {
		// query[i] = float64(rand.Intn(2))
        query[i] = 0.0
	}
    query[0] = 1.0
	return query
}
