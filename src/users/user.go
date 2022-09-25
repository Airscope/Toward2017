package users

import (
    "github.com/Airscope/Toward2017/utils"
	"github.com/sachaservan/bgn"
	"math/rand"
    "math"
    "os"
    "bufio"
    "io"
    "strings"
    "strconv"
    "fmt"
)

func DataPacking(pk *bgn.PublicKey, m, n, numPacking, numInterval int) [][] *bgn.Ciphertext {
	trans := readTransactions(m, n)
    // 垂直packing，将trans[i][j]至trans[i][j+numPacking]packing为一个大整数
    packedTrans := make([][] *bgn.Plaintext, m / numPacking)

    k10PowInterval := int(math.Pow(10, float64(numInterval))) // 10 ^ numInterval

    for i := 0; i < m / numPacking; i ++ {
        packedTranRow := make([] *bgn.Plaintext, n) // 1行packed的事务，长度为n，包含n个packed的数据
        for j := 0; j < n; j++ {
            packedInt := 0
            for k := 0; k < numPacking; k++ {
                packedInt = int(trans[i*numPacking+k][j]) + packedInt
                if (k != numPacking - 1) {
                    packedInt = k10PowInterval * packedInt
                }
            }
            packedTranRow[j] = utils.BGNPlaintxt(pk, packedInt)
         }
        packedTrans[i] = packedTranRow
    }


	encTrans := make([][] *bgn.Ciphertext, len(packedTrans))
	for i := 0; i < len(packedTrans); i++ {
		tmp := make([] *bgn.Ciphertext, n)
		for j := 0; j < n; j++ {
			tmp[j] = pk.Encrypt(packedTrans[i][j])
		}
		encTrans[i] = tmp
	}
	return encTrans
}

func DataProcess(pk *bgn.PublicKey, m, n int) [][] *bgn.Ciphertext {
	trans := readTransactions(m, n)
	encTrans := make([][] *bgn.Ciphertext, m)
	for i := 0; i < m; i++ {
		tmp := make([] *bgn.Ciphertext, n)
		for j := 0; j < n; j++ {
			tmp[j] = pk.Encrypt(utils.BGNPlaintxt(pk, int(trans[i][j]), ))
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
