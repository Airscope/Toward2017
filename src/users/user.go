package users

import (
	"github.com/sachaservan/bgn"
	"math/big"
	"math/rand"
    "os"
    "bufio"
    "io"
    "strings"
    "strconv"
    "fmt"
)

func DataProcess(pk *bgn.PublicKey, m, n int) [][] *bgn.Ciphertext {
	trans := readTransactions(m, n)
	encTrans := make([][] *bgn.Ciphertext, m)
	for i := 0; i < m; i++ {
		tmp := make([] *bgn.Ciphertext, n)
		for j := 0; j < n; j++ {
			tmp[j] = pk.Encrypt(big.NewInt(int64(trans[i][j])))
		}
		encTrans[i] = tmp
	}
	return encTrans
}

func readTransactions(m, n int) [][] float64 {
	var trans [][] float64
	fromDisk := true // false
	if !fromDisk {
		for i := 0; i < m; i++ {
			tmp := make([] float64, n)
			for j := 0; j < n; j++ {
				tmp[j] = float64(rand.Intn(2))
			}
			trans = append(trans, tmp)
		}
	} else {
		for i := 0; i < m; i++ {
			tmp := make([] float64, n)
			for j := 0; j < n; j++ {
				tmp[j] = float64(0)
			}
			trans = append(trans, tmp)
		}
        fmt.Printf("Read a %d * %d boolean matrix from disk\n", m, n)
        wd, _ := os.Getwd()
        filePath := (wd + "/../data/dataset_chess.txt")
        println("Dataset filepath: " + filePath)
        file, err := os.OpenFile(filePath, os.O_RDONLY, 0600)
        if err != nil {
            println("Open file error.", err)
            panic(err)
        }
        defer file.Close()
        //stat, err := file.Stat()
        //if err != nil {
        //    panic(err)
        //}
        // size := stat.Size()
        // println("file size=", size)
        buf := bufio.NewReader(file)
        for i := 0; i < m; i++{
            line, err := buf.ReadString('\n')
            if err != nil {
                if err == io.EOF {
                    break
                } else {
                    panic(err)
                }
            }
            labels := strings.Split(line, " ")
            for _, label := range labels {
                j, _ := strconv.Atoi(label)
                if j > n {
                    break
                }
                trans[i][j-1] = float64(1)
            }
        }
	}
    //for i := range trans {
    //    fmt.Println(trans[i])
    //}
	return trans
}
