package main

import (
    "github.com/Airscope/Toward2017/protocols"
)

func main() {
	const TESTNUM = 1
    for i := 0; i < TESTNUM; i++ {
        println("Test at", i+1, "times \n=================================")
        protocols.RunParallel()
        println("==================================")
    }
}
