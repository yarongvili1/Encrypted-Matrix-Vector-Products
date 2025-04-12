package splitlsnpir

import (
	"RandomLinearCodePIR/pir"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestFullFunctionOfSLSNPIR(t *testing.T) {
	row := uint32(1 << 14)
	l := uint32(1 << 14)
	s := uint32(16)
	k := uint32(16)
	n := k + l
	// Largest Prime in 32 bits
	p := uint32(1<<32 - 5)
	seed := int64(1)

	pi := &SlsnPIR{Params: SlsnParams{
		PrimeField: p,
		NumBlocks:  s,
		K:          k,
		N:          n,
		M:          row,
		L:          l,
	}}

	matrix := pir.GeneratePrimeFieldMatrix(pi.Params.M, pi.Params.L, p, seed)

	fmt.Printf("Running PIR with Database %d * %d \n", pi.Params.M, pi.Params.L)

	queryIndex := rand.Uint64() % (uint64(pi.Params.M) * uint64(pi.Params.L))

	fmt.Println("Generate Key...")
	start := time.Now()
	sk := pi.KeyGen(seed)
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
	serverResponse := pi.Answer(*encodedMatrix, clientQuery)
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
