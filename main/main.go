package main

import (
    "github.com/Airscope/protocols"
)

func main() {
	const TESTNUM = 5
    for i := 0; i < TESTNUM; i++ {
        println("Test at", i+1, "--------------------")
        protocols.Run()
        println("-----------------------")
    }
}
