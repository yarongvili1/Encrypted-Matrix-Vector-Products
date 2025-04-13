package lpnpir

import (
	"RandomLinearCodePIR/pir"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestLPNBasedPIRComplete(t *testing.T) {
	row := uint32(1 << 10)
	l := uint32(1 << 10)
	k := uint32(1 << 4)
	n := k + l
	// Largest Prime in 32 bits
	p := uint32(2 ^ 32 - 5)
	seed := int64(1)
	m_1 := uint32(4)

	pi := &LpnPIR{Params: LpnParams{
		PrimeField: p,
		K:          k,
		N:          n,
		M:          row,
		L:          l,
		M_1:        m_1,
		ECCLength:  7,
		Epsi:       0.01,
	}}

	rng := rand.New(rand.NewSource(1))

	originalData := make([]uint32, pi.Params.M*pi.Params.L)
	for i := range originalData {
		originalData[i] = rng.Uint32() % p
	}

	// fmt.Println("Original Data: ", originalData)
	matrix := pir.Matrix{
		Rows: pi.Params.M,
		Cols: pi.Params.L,
		Data: originalData,
	}

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
	serverResponse := pi.Answer(encodedMatrix, clientQuery)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Decode...")
	start = time.Now()
	val, err := pi.Decode(sk, queryIndex, serverResponse, aux)
	fmt.Println("    Elapsed: ", time.Since(start))

	if err == nil && val != matrix.Data[queryIndex] {
		fmt.Println(val)
		panic("INCORRECT RESULT!")
	}
}
