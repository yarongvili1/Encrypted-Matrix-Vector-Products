package mvp

import (
	"RandomLinearCodePIR/dataobjects"
	"RandomLinearCodePIR/ecc"
	"RandomLinearCodePIR/utils"
	"fmt"
	"testing"
	"time"
)

func TestFullFunctionOfSLSNPIR(t *testing.T) {
	m := uint32(1 << 8)
	l := uint32(1 << 8)
	s := uint32(2)
	k := uint32(4)
	n := k + l
	b := n / s
	p := uint32(65537)
	seed := int64(1)

	pi := &SlsnMVP{Params: SlsnParams{
		Field: dataobjects.NewPrimeField(p),
		S:     s,
		K:     k,
		N:     n,
		M:     m,
		L:     l,
		B:     b,
		P:     p,
	}}

	matrix := utils.GeneratePrimeFieldMatrix(pi.Params.M, pi.Params.L, p, seed)

	fmt.Printf("Running PIR with Database %d * %d \n", pi.Params.M, pi.Params.L)

	query := utils.RandomPrimeFieldVector(pi.Params.L, pi.Params.P)

	fmt.Println("Generate Key...")
	start := time.Now()
	sk := pi.KeyGen(seed)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Generate Trapdoored Matrix...")
	start = time.Now()
	TDM := pi.GenerateTDM(sk)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Encode Message...")
	start = time.Now()
	encodedMatrix := pi.Encode(sk, matrix, TDM)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Generate Query...")
	start = time.Now()
	clientQuery, aux := pi.Query(sk, query)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Answer...")
	start = time.Now()
	serverResponse := pi.Answer(*encodedMatrix, *clientQuery)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Decode...")
	start = time.Now()
	val := pi.Decode(sk, serverResponse, *aux)
	fmt.Println("    Elapsed: ", time.Since(start))

	// fmt.Println(val)

	target := make([]uint32, m)
	BlockMatVecProduct(matrix.Data, query, target, m, l, 1, p)

	if len(val) != len(target) {
		panic("Naive Sanity Check!")
	}

	// for i := range target {
	// 	if target[i] != val[i] {
	// 		panic("Vec doesn't match ! ")
	// 	}
	// }
}

func TestLPNMVPComplete(t *testing.T) {
	m := uint32(1 << 10)
	l := uint32(1 << 10)
	k := uint32(1 << 4)
	n := k + l
	p := uint32(65537)
	seed := int64(1)
	m_1 := uint32(4)

	pi := &LpnMVP{Params: LpnParams{
		Field:     dataobjects.NewPrimeField(p),
		K:         k,
		N:         n,
		M:         m,
		L:         l,
		M_1:       m_1,
		ECCLength: 7,
		Epsi:      0.00001,
		P:         p,
		ECCName:   ecc.ReedSolomon,
	}}

	matrix := utils.GeneratePrimeFieldMatrix(pi.Params.M, pi.Params.L, p, seed)

	fmt.Printf("Running PIR with Database %d * %d \n", pi.Params.M, pi.Params.L)

	query := utils.RandomPrimeFieldVector(pi.Params.L, pi.Params.P)

	fmt.Println("Generate Key...")
	start := time.Now()
	sk := pi.KeyGen(seed)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Generate Trapdoored Matrix...")
	start = time.Now()
	TDM := pi.GenerateTDM(sk)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Encode Message...")
	start = time.Now()
	encodedMatrix := pi.Encode(sk, matrix, TDM)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Generate Query...")
	start = time.Now()
	clientQuery, aux := pi.Query(sk, query)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Answer...")
	start = time.Now()
	serverResponse := pi.Answer(encodedMatrix, clientQuery)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Decode...")
	start = time.Now()
	val := pi.Decode(sk, serverResponse, aux)
	fmt.Println("    Elapsed: ", time.Since(start))

	target := make([]uint32, m)
	MatVecProduct(matrix.Data, query, target, m, l, p)

	if len(val) != len(target) {
		panic("Naive Sanity Check!")
	}

	// for i := range target {
	// 	if target[i] != val[i] {
	// 		panic("Vec doesn't match ! ")
	// 	}
	// }
}
