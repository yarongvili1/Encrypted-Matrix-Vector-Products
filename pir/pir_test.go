package pir

import (
	"RandomLinearCodePIR/utils"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestBasePIR(t *testing.T) {
	lambda := uint32(32)
	row := uint32(1 << 8)
	col := uint32(1 << 6)

	pi := &BasePIR{
		Params: BaseParams{
			Rows:           row,
			Cols:           col,
			NumberOfBlocks: uint32(16),
			CodewordLength: col + lambda,
		},
	}

	matrix := GenerateMatrix(pi.Params.Rows, pi.Params.Cols, 1, 1)
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

func BenchmarkQueryGeneration(b *testing.B) {
	row := uint32(1 << 19)
	col := uint32(1 << 14)
	k := uint32(4176)

	pi := &BasePIR{
		Params: BaseParams{
			Rows:           row,
			Cols:           col,
			NumberOfBlocks: uint32(80),
			CodewordLength: col + k,
		},
	}

	sk := pi.KeyGen(1, 2, 32, 1)

	var totalDuration time.Duration
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		queryIndex := uint64(rand.Intn(int(row) * int(col)))
		start := time.Now()
		pi.Query(sk, queryIndex)
		totalDuration += time.Since(start)
	}

	b.StopTimer()

	fmt.Printf("\nBenchmark Query Generation For m = %d, l = %d, k = %d\n", row, col, k)
	fmt.Printf("\tAverage Query Time with repeat times: %d: %v \n\n", b.N, totalDuration/time.Duration(b.N))
}

func BenchmarkAnswer(b *testing.B) {
	row := uint32(1 << 19)
	col := uint32(1 << 14)
	k := uint32(4176)

	pi := &BasePIR{
		Params: BaseParams{
			Rows:           row,
			Cols:           col,
			NumberOfBlocks: uint32(80),
			CodewordLength: col + k,
			PackedSize:     row / 32,
		},
	}

	encodedMatrix := GenerateMatrix(pi.Params.CodewordLength, pi.Params.PackedSize, 32, 1)

	var totalDuration time.Duration
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		vec_1 := utils.RandomizeBinaryVector(pi.Params.CodewordLength)
		vec_2 := utils.RandomizeBinaryVector(pi.Params.CodewordLength)
		start := time.Now()
		pi.Answer(&encodedMatrix, &BasePIRQuery{Vector_1: vec_1, Vector_2: vec_2})
		totalDuration += time.Since(start)
	}

	b.StopTimer()

	fmt.Printf("\nBenchmark Server Response For Encoded Database with %d x %d entires of size ~%fMB. \n",
		encodedMatrix.Rows, encodedMatrix.Cols, float64(encodedMatrix.Rows/1024.0)*float64(encodedMatrix.Cols)/1024.0*32/8.0)
	fmt.Printf("\tAverage Answer Time with repeat times: %d: %v \n\n", b.N, totalDuration/time.Duration(b.N))
}

func BenchmarkDecode(b *testing.B) {
	row := uint32(1 << 19)
	col := uint32(1 << 14)
	k := uint32(4176)

	pi := &BasePIR{
		Params: BaseParams{
			Rows:           row,
			Cols:           col,
			NumberOfBlocks: uint32(80),
			CodewordLength: col + k,
			PackedSize:     row / 32,
		},
	}

	sk := pi.KeyGen(1, 2, 32, 1)

	var totalDuration time.Duration
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		index := uint64(rand.Intn(int(row)))
		res_1 := utils.RandomPrimeFieldVector(pi.Params.PackedSize*pi.Params.NumberOfBlocks, 2^32)
		res_2 := utils.RandomPrimeFieldVector(pi.Params.PackedSize*pi.Params.NumberOfBlocks, 2^32)
		flip := utils.RandomizeFlipVector(pi.Params.NumberOfBlocks)
		mask := rand.Uint32()

		start := time.Now()
		pi.Decode(sk, index, &BasePIRAnswer{Result_1: res_1, Result_2: res_2},
			&BasePIRAux{FlipVector: flip, MaskValue: mask})
		totalDuration += time.Since(start)
	}

	b.StopTimer()

	fmt.Printf("\nBenchmark Decoding For m = %d, l = %d, k = %d, with repeat times: %d\n", row, col, k, b.N)
	fmt.Printf("\tAverage Decoding Time with repeat times: %d: %v \n\n", b.N, totalDuration/time.Duration(b.N))
}
