package main

import (
    "github.com/Airscope/Toward2017/protocols"
)

func main() {
	const TESTNUM = 1

    println(">>>>>>>>>>>>>>>SIMD Protocol Start <<<<<<<<<<<<<<<<")
    for i := 0; i < TESTNUM; i++ {
        println("Test at", i+1, "times \n=================================")
        protocols.RunSIMD()
        println("==================================")
    }
    println(">>>>>>>>>>>>>>>SIMD Protocol End.<<<<<<<<<<<<<<<<")

    println("\n>>>>>>>>>>>>>>>Parallel Protocol Start <<<<<<<<<<<<<<<<")
    for i := 0; i < TESTNUM; i++ {
        println("Test at", i+1, "times \n=================================")
        protocols.RunParallel()
        println("==================================")
    }
    println(">>>>>>>>>>>>>>>Parallel Protocol End.<<<<<<<<<<<<<<<<")
}
