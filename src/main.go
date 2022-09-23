package main

import (
    "github.com/Airscope/Toward2017/protocols"
)

func main() {
	const TESTNUM = 1

    println(">>>>>>>>>>>>>>>Prallel SIMD Protocol 2 Start <<<<<<<<<<<<<<<<")
    for i := 0; i < TESTNUM; i++ {
        println("Test at", i+1, "times \n=================================")
        protocols.RunParallelSIMD()
        println("==================================")
    }
    println(">>>>>>>>>>>>>>>Prallel SIMD Protocol2 End.<<<<<<<<<<<<<<<<")
}
