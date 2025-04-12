package splitlsnpir

import (
	"RandomLinearCodePIR/pir"
	"math/rand"
)

func GenerateRLC(K, L, p uint32, seed int64) [][]uint32 {
	randomP := GenerateRandomColsOfC(K, L, p, seed)
	matrix := make([][]uint32, K)

	for i := uint32(0); i < K; i++ {
		matrix[i] = make([]uint32, L+K)
		matrix[i][i] = 1
		for j := uint32(0); j < L; j++ {
			matrix[i][K+j] = randomP[i][j]
		}
	}

	return matrix
}

func GenerateDualCode(K, L, p uint32, seed int64) [][]uint32 {
	randomP := GenerateRandomColsOfD(K, L, p, seed)
	matrix := make([][]uint32, K+L)

	for i := uint32(0); i < K; i++ {
		matrix[i] = make([]uint32, L)
		matrix[i] = randomP[i]
	}

	for i := uint32(0); i < L; i++ {
		matrix[K+i] = make([]uint32, L)
		matrix[K+i][i] = 1
	}

	return matrix
}

// C = (I | P) and D = (-P // I) where P has dimension K x L
func GenerateRandomColsOfC(K, L, p uint32, seed int64) [][]uint32 {
	matrix := make([][]uint32, K)

	for i := uint32(0); i < K; i++ {
		matrix[i] = make([]uint32, L)
	}

	rng := rand.New(rand.NewSource(seed))

	for j := uint32(0); j < L; j++ {
		for i := uint32(0); i < K; i++ {
			matrix[i][j] = uint32(rng.Intn(int(p)))
		}
	}

	return matrix
}

// D = (-P // I) Generate -P from C
func GenerateRandomColsOfD(K, L, p uint32, seed int64) [][]uint32 {
	matrix := GenerateRandomColsOfC(K, L, p, seed)

	for j := uint32(0); j < L; j++ {
		for i := uint32(0); i < K; i++ {
			matrix[i][j] = (p - matrix[i][j]) % p
		}
	}

	return matrix
}

// Transpose D = (-P // I) to D' = (-P^T | I) for easy row combination
func Generate1DDualSpaceRandomMatrix(K, L, p uint32, seed int64) []uint32 {
	matrix := GenerateRandomColsOfD(K, L, p, seed)

	vmatrix := make([]uint32, K*L)

	idx := 0

	for i := uint32(0); i < L; i++ {
		for j := uint32(0); j < K; j++ {
			vmatrix[idx] = matrix[j][i]
			idx += 1
		}
	}

	return vmatrix
}

// Transpose D = (-P // I) to D' = (-P^T | I) for easy row combination
func Generate1DRLCMatrix(K, L, p uint32, seed int64) []uint32 {
	matrix := GenerateRandomColsOfC(K, L, p, seed)

	vmatrix := make([]uint32, K*L)

	idx := 0

	for i := uint32(0); i < K; i++ {
		for j := uint32(0); j < L; j++ {
			vmatrix[idx] = matrix[i][j]
			idx += 1
		}
	}

	return vmatrix
}

// D has dimension N x L
func SampleVectorFromNullSpace(K, L, p uint32, seed int64) []uint32 {
	coeff := pir.RandomPrimeFieldVector(L, p)

	matrix := GenerateRandomColsOfD(K, L, p, seed)

	res := make([]uint32, K+L)

	for i := uint32(0); i < L; i++ {

		for j := uint32(0); j < K; j++ {
			res[j] = uint32((uint64(res[j]) + uint64(matrix[j][i])*uint64(coeff[i])) % uint64(p))
		}
		res[K+i] = uint32((uint64(res[K+i]) + uint64(coeff[i])) % uint64(p))

	}

	return res
}
