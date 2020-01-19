package protocols

import (
	"github.com/Airscope/cloud"
	"github.com/Airscope/evaluator"
	"github.com/Airscope/miner"
	"github.com/Airscope/users"
	"math/big"
    "time"
)

func Run() {
<<<<<<< HEAD
	const M = 1000
//	const N = 20
	const K = M / 2
	const MINSUPP = M * 4 / 5

    // const ITER = 5
    nRange := []int{10,20,30,40,50}


    for i := range nRange {
        N := nRange[i]
        var cloudTime int64
        var evaTime int64
        cloudTime = 0
        evaTime = 0

	    println("Step 1 System Setup")
	    pkBGN, skBGN, pkPaillier, skPaillier := evaluator.SystemSetup()

	    println("Step 2 Data Processing")
	    encTrans := users.DataProcess(pkBGN, M, N)

	    println("Step 3 Computation")
    
        originStart := time.Now().UnixNano()
	    println("at miner...")
	    encQuery, negL1Norm := miner.Compute(pkBGN, N)
        negL1Norm = pkBGN.EMult(negL1Norm, pkBGN.Encrypt(pkBGN.NewPlaintext(big.NewFloat(1.0))))
	    println("at cloud...")
        startTime := time.Now().UnixNano()
	    randomizedSet := cloud.Compute(encTrans, encQuery, negL1Norm, pkBGN, K, N)
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
	    minSupp := pkPaillier.Encrypt(big.NewInt(MINSUPP))
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
=======
	const M = 100
	const N = 20
	const K = M / 2
	const MINSUPP = M * 4 / 5

    var cloudTime int64
    var evaTime int64
    cloudTime = 0
    evaTime = 0

	println("Step 1 System Setup")
	pkBGN, skBGN, pkPaillier, skPaillier := evaluator.SystemSetup()

	println("Step 2 Data Processing")
	encTrans := users.DataProcess(pkBGN, M, N)

	println("Step 3 Computation")
    
    originStart := time.Now().UnixNano()
	println("at miner...")
	encQuery, negL1Norm := miner.Compute(pkBGN, N)
	println("at cloud...")
    startTime := time.Now().UnixNano()
	randomizedSet := cloud.Compute(encTrans, encQuery, negL1Norm, pkBGN, K, N)
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
	minSupp := pkPaillier.Encrypt(big.NewInt(MINSUPP))
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

>>>>>>> 776741279ca80cebcd71cf7a6900206b7bfa3665
}
