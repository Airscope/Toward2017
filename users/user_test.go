package users

import (
    "testing"
    "fmt"
)

func TestReadTranscations(t *testing.T) {

    read := readTransactions(12, 14)

    fmt.Printf("Read a %d * %d matrix", len(read), len(read[0]))
    fmt.Println(read)
}
