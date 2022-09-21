package protocols

import (
	"github.com/Airscope/Toward2017/cloud"
	"github.com/Airscope/Toward2017/evaluator"
	"github.com/Airscope/Toward2017/miner"
	"github.com/Airscope/Toward2017/users"
	"math/big"
    "time"
	gmp "github.com/ncw/gmp"

)

func RunParallel() {
	const M = 20
	const N = 20
	const K = M / 2
	const MINSUPP = M * 4 / 5
    numsCPU := []int{3,5,7,9,11}

    for i := range numsCPU {
		numCPU := numsCPU[i]
        var cloudTime int64
        var evaTime int64
        cloudTime = 0
        evaTime = 0

	    println("---------------------\nStep 1 System Setup")
	    pkBGN, skBGN, pkPaillier, skPaillier := evaluator.SystemSetup()

	    println("Step 2 Data Processing")
	    encTrans := users.DataProcess(pkBGN, M, N)

	    println("Step 3 Parallel Computation")
		println("[CPUs number: ", numCPU, "]")

        originStart := time.Now().UnixNano()
	    println("at miner...")
	    encQuery, negL1Norm := miner.Compute(pkBGN, N)
        negL1Norm = pkBGN.Mult(negL1Norm, pkBGN.Encrypt(big.NewInt(1)))
	    println("at cloud...")
        startTime := time.Now().UnixNano()
	    randomizedSet := cloud.ParallelCompute(encTrans, encQuery, negL1Norm, pkBGN, K, N, numCPU)
        cloudTime += time.Now().UnixNano() - startTime

	    println("Step 4 Evaluation")
	    println("at evaluator...")
        startTime = time.Now().UnixNano()
	    v := evaluator.Evaluate(randomizedSet, pkBGN, skBGN, pkPaillier)
        evaTime += time.Now().UnixNano() - startTime
	    println("at cloud...")
        startTime = time.Now().UnixNano()
	    support := cloud.Evaluate(v, pkPaillier, M, K)
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
