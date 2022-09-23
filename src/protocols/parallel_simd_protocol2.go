package protocols

import (
	"github.com/Airscope/Toward2017/cloud"
	"github.com/Airscope/Toward2017/evaluator"
	"github.com/Airscope/Toward2017/miner"
	"github.com/Airscope/Toward2017/users"
	"math/big"
	"math"
    "time"
	gmp "github.com/ncw/gmp"

)

func RunParallelSIMD() {
	const M = 3000
	const N = 20
	const K = M / 2
	const MINSUPP = M * 4 / 5
	const numPacking = 2 // packing的trans个数
	numsCPU := []int{1,2,3,4,5,6,7,8,9,10,11,12}
	numInterval := int(math.Floor(math.Log10(float64(N)))) + 1 // 保证不会溢出的最小packing间隔F

	if (M % numPacking != 0) {
		panic("M must be divided by numPacking.")
	}

    for i := range numsCPU {
        numCPU := numsCPU[i]

        var cloudTime int64
        var evaTime int64
        cloudTime = 0
        evaTime = 0

		println("----------------------\n[CPUs number: ", numCPU, "]")
		println("[number of Packing is ", numPacking, ", the interval is ", numInterval, "]")

		println("Step 1 System Setup")
	    pkBGN, skBGN, pkPaillier, skPaillier := evaluator.SystemSetup()

	    println("Step 2 Data Processing")
	    encTrans := users.DataProcess(pkBGN, M, N)

		println("Step 3 Parallel SIMD Computation")

        originStart := time.Now().UnixNano()
	    println("at miner...")
	    encQuery, negL1Norm := miner.Compute(pkBGN, N)
        negL1Norm = pkBGN.Mult(negL1Norm, pkBGN.Encrypt(big.NewInt(1)))
	    println("at cloud...")
        startTime := time.Now().UnixNano()
	    randomizedSet := cloud.ParallelSIMDCompute(encTrans, encQuery, negL1Norm, pkBGN, skBGN, K, N, numCPU, numPacking, numInterval)
        cloudTime += time.Now().UnixNano() - startTime

	    println("Step 4 Evaluation")
	    println("at evaluator...")
        startTime = time.Now().UnixNano()
	    v := evaluator.SIMDEvaluate(randomizedSet, pkBGN, skBGN, pkPaillier, M, K, numPacking, numInterval)
        evaTime += time.Now().UnixNano() - startTime
	    println("at cloud...")
        startTime = time.Now().UnixNano()
	    support := cloud.ParallelSIMDEvaluate(v, pkPaillier, M, K)
        evaTime += time.Now().UnixNano() - startTime

	    println("Step 5 Comparison")
	    minSupp := pkPaillier.Encrypt(gmp.NewInt(MINSUPP))
	    println("at cloud...")
        startTime = time.Now().UnixNano()
	    maskedFlag := cloud.Compare(support, minSupp, pkPaillier) // return support < minSupp
	    cloudTime += time.Now().UnixNano() - startTime
        println("at evaluator...")
        startTime = time.Now().UnixNano()
	    output := evaluator.Compare(maskedFlag, skPaillier)
        evaTime += time.Now().UnixNano() - startTime

        println("Mining output:", output)
        println("Cloud time:", float64(cloudTime/int64(time.Millisecond)), "milliseconds")
        println("Evaluator time:", float64(evaTime/int64(time.Millisecond)), "milliseconds")
        println("Total time:", float64((time.Now().UnixNano() - originStart)/int64(time.Millisecond)), "milliseconds")
        println()

    }
}
