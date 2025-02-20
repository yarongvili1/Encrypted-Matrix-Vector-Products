package main

import (
	"RandomLinearCodePIR/pir"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func RunBasePIR() {
	lambda := uint32(32)
	row := uint32(math.Pow(2, 8))
	col := uint32(math.Pow(2, 6))

	pi := &pir.BasePIR{
		Params: pir.BaseParams{
			Rows:           row,
			Cols:           col,
			NumberOfBlocks: uint32(16),
			CodewordLength: col + lambda,
		},
	}

	matrix := pir.GenerateMatrix(pi.Params.Rows, pi.Params.Cols, 1, 1)
	queryIndex := uint64(rand.Intn(int(row) * int(col)))

	fmt.Printf("Running PIR with Database %d * %d \n", pi.Params.Rows, pi.Params.Cols)

	fmt.Println("Generate Key...")
	start := time.Now()
	sk := pi.KeyGen(1, 2, 32, 1)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Encode Message...")
	start = time.Now()
	encodedMatrix := pi.Encode(sk, matrix)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Generate Query...")
	start = time.Now()
	clientQuery, aux := pi.Query(sk, queryIndex)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Answer...")
	start = time.Now()
	serverResponse := pi.Answer(encodedMatrix, clientQuery)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Decode...")
	start = time.Now()
	val := pi.Decode(sk, queryIndex, serverResponse, aux)
	fmt.Println("    Elapsed: ", time.Since(start))

	if val != matrix.Data[queryIndex] {
		fmt.Println(val)
		panic("INCORRECT RESULT!")
	}
}

func main() {
	RunBasePIR()
}
