package mvp

import (
	"RandomLinearCodePIR/dataobjects"
	"RandomLinearCodePIR/ecc"
	"RandomLinearCodePIR/linearcode"
	"RandomLinearCodePIR/utils"
	"fmt"
	"testing"
	"time"
)

func TestSlsnMVPComplete(t *testing.T) {
	m := uint32(1 << 10)
	l := uint32(1 << 10)
	k := uint32(1 << 4)
	s := uint32(2)
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

	fmt.Printf("Running MVP with Database %d * %d \n", pi.Params.M, pi.Params.L)

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
	fmt.Println("    Include Calculate Mask Time: ", aux.Dur)

	fmt.Println("Answer...")
	start = time.Now()
	serverResponse := pi.Answer(*encodedMatrix, *clientQuery)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Decode...")
	start = time.Now()
	val := pi.Decode(sk, serverResponse, *aux)
	fmt.Println("    Elapsed: ", time.Since(start))

	target := make([]uint32, m)
	BlockMatVecProduct(matrix.Data, query, target, m, l, 1, p)

	for i := range target {
		if target[i] != val[i] {
			panic("Vec doesn't match ! ")
		}
	}

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

	fmt.Printf("Running MVP with Database %d * %d \n", pi.Params.M, pi.Params.L)

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

	for i := range target {
		if target[i] != val[i] {
			panic("Vec doesn't match ! ")
		}
	}
}

func TestRingSlsnMVPComplete(t *testing.T) {
	m := uint32(1 << 10)
	l := uint32(1 << 10)
	k := uint32(1 << 4)
	s := uint32(2)
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

	code := linearcode.GetLinearCode(
		linearcode.LinearCodeConfig{
			Name:  linearcode.Vandermonde,
			K:     k,
			L:     l,
			Field: dataobjects.NewPrimeField(p),
		},
	)

	ring := &RingSlsnMVP{
		SlsnMVP:           *pi,
		LinearCodeEncoder: code,
	}

	matrix := utils.GeneratePrimeFieldMatrix(pi.Params.M, pi.Params.L, p, seed)
	query := utils.RandomPrimeFieldVector(pi.Params.L, pi.Params.P)

	target := make([]uint32, m)
	MatVecProduct(matrix.Data, query, target, m, l, p)

	fmt.Printf("Running MVP with Database %d * %d \n", pi.Params.M, pi.Params.L)

	fmt.Println("Generate Key...")
	start := time.Now()
	sk := ring.KeyGen(seed)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Generate Trapdoored Matrix...")
	start = time.Now()
	TDM := ring.GenerateTDM(sk)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Encode Message...")
	start = time.Now()
	encodedMatrix := ring.Encode(sk, matrix, TDM)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Generate Query...")
	start = time.Now()
	clientQuery, aux := ring.Query(sk, query)
	fmt.Println("    Elapsed: ", time.Since(start))
	fmt.Println("    Include Calculate Mask Time: ", aux.Dur)

	fmt.Println("Answer...")
	start = time.Now()
	serverResponse := ring.Answer(*encodedMatrix, *clientQuery)
	fmt.Println("    Elapsed: ", time.Since(start))

	fmt.Println("Decode...")
	start = time.Now()
	val := ring.Decode(sk, serverResponse, *aux)
	fmt.Println("    Elapsed: ", time.Since(start))

	for i := range target {
		if target[i] != val[i] {
			panic("Vec doesn't match ! ")
		}
	}
}

func BenchmarkCleartextServerExecution(b *testing.B) {
	m := uint32(1) << 10
	l := uint32(1) << 10
	p := uint32(65537)
	seed := int64(1)
	matrix := utils.GeneratePrimeFieldMatrix(m, l, p, seed)
	result := make([]uint32, m)

	var totalDuration time.Duration
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		query := utils.RandomPrimeFieldVector(l, p)
		start := time.Now()
		MatVecProduct(matrix.Data, query, result, m, l, p)
		duration := time.Since(start)
		totalDuration += duration
	}

	b.StopTimer()
	fmt.Printf("Average server execution time for m = %d, l = %d : %v\n", m, l, totalDuration/time.Duration(b.N))
}

func BenchmarkRingSLSNQuery(b *testing.B) {
	m := uint32(1 << 10)
	l := uint32(1 << 10)
	k := uint32(1 << 4)
	s := uint32(2)
	n := k + l
	block := n / s
	p := uint32(65537)
	seed := int64(1)

	pi := &SlsnMVP{Params: SlsnParams{
		Field: dataobjects.NewPrimeField(p),
		S:     s,
		K:     k,
		N:     n,
		M:     m,
		L:     l,
		B:     block,
		P:     p,
	}}

	code := linearcode.GetLinearCode(
		linearcode.LinearCodeConfig{
			Name:  linearcode.Vandermonde,
			K:     k,
			L:     l,
			Field: dataobjects.NewPrimeField(p),
		},
	)

	ring := &RingSlsnMVP{
		SlsnMVP:           *pi,
		LinearCodeEncoder: code,
	}

	sk := ring.KeyGen(seed)
	ring.GenerateTDM(sk)

	var totalDuration time.Duration
	var unmaskDuration time.Duration

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		query := utils.RandomPrimeFieldVector(pi.Params.L, pi.Params.P)
		start := time.Now()
		_, aux := ring.Query(sk, query)
		duration := time.Since(start)
		totalDuration += duration
		unmaskDuration += aux.Dur
	}

	b.StopTimer()

	fmt.Printf("Ring SLSN For m = %d, l = %d, k = %d \n", m, l, k)
	fmt.Printf("Average Query time: %v\n", totalDuration/time.Duration(b.N))
	fmt.Printf("Average Calculate Mask time: %v\n", unmaskDuration/time.Duration(b.N))
}

func BenchmarkSLSNQuery(b *testing.B) {
	m := uint32(1 << 10)
	l := uint32(1 << 10)
	k := uint32(1 << 4)
	s := uint32(2)
	n := k + l
	block := n / s
	p := uint32(65537)
	seed := int64(1)

	pi := &SlsnMVP{Params: SlsnParams{
		Field: dataobjects.NewPrimeField(p),
		S:     s,
		K:     k,
		N:     n,
		M:     m,
		L:     l,
		B:     block,
		P:     p,
	}}

	sk := pi.KeyGen(seed)
	pi.GenerateTDM(sk)

	var totalDuration time.Duration
	var unmaskDuration time.Duration

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		query := utils.RandomPrimeFieldVector(pi.Params.L, pi.Params.P)
		start := time.Now()
		_, aux := pi.Query(sk, query)
		duration := time.Since(start)
		totalDuration += duration
		unmaskDuration += aux.Dur
	}

	b.StopTimer()

	fmt.Printf("SLSN For m = %d, l = %d, k = %d \n", m, l, k)
	fmt.Printf("Average Query time: %v\n", totalDuration/time.Duration(b.N))
	fmt.Printf("Average Calculate Mask time: %v\n", unmaskDuration/time.Duration(b.N))
}

func BenchmarkSLSNAnswer(b *testing.B) {
	m := uint32(1 << 10)
	l := uint32(1 << 10)
	k := uint32(1 << 4)
	s := uint32(2)
	n := k + l
	block := n / s
	p := uint32(65537)
	seed := int64(1)

	pi := &SlsnMVP{Params: SlsnParams{
		Field: dataobjects.NewPrimeField(p),
		S:     s,
		K:     k,
		N:     n,
		M:     m,
		L:     l,
		B:     block,
		P:     p,
	}}

	matrix := utils.GeneratePrimeFieldMatrix(pi.Params.M, pi.Params.L, p, seed)

	sk := pi.KeyGen(seed)
	TDM := pi.GenerateTDM(sk)
	encodedMatrix := pi.Encode(sk, matrix, TDM)

	var totalDuration time.Duration

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		query := utils.RandomPrimeFieldVector(pi.Params.L, pi.Params.P)
		clientQuery, _ := pi.Query(sk, query)

		start := time.Now()
		pi.Answer(*encodedMatrix, *clientQuery)
		totalDuration += time.Since(start)
	}

	b.StopTimer()

	fmt.Printf("Benchmark of SLSN Answer For m = %d, l = %d, k = %d \n", m, l, k)
	fmt.Printf("Average Answer time: %v\n", totalDuration/time.Duration(b.N))
}
