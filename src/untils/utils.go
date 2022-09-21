package utils

import (
	"github.com/sachaservan/bgn"
	"github.com/sachaservan/paillier"
)

func

func Int2GMPInt(i int) *gmp.Int {
	return gmp.NewInt(int64(i))
}

func GMPInt2Int(i *gmp.Int) int {
	return int(i.Int64())
}
