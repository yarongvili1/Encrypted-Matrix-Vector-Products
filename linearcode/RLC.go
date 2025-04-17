package linearcode

import (
	"RandomLinearCodePIR/dataobjects"
	"math/rand"
)

// // Generate D = (I | P) with dimension L x (L + K)
// func GenerateRLC(L, K, p uint32, seed int64) [][]uint32 {
// 	randomP := GenerateRandomColsOfD(L, K, p, seed)
// 	matrix := make([][]uint32, L)

// 	for i := uint32(0); i < L; i++ {
// 		matrix[i] = make([]uint32, L+K)
// 		matrix[i][i] = 1
// 		for j := uint32(0); j < K; j++ {
// 			matrix[i][L+j] = randomP[i][j]
// 		}
// 	}

// 	return matrix
// }

// func GenerateDualCode(L, K, p uint32, seed int64) [][]uint32 {
// 	randomP := GenerateRandomColsOfC(L, K, p, seed)
// 	matrix := make([][]uint32, K+L)

// 	for i := uint32(0); i < L; i++ {
// 		matrix[i] = make([]uint32, K)
// 		matrix[i] = randomP[i]
// 	}

// 	for i := uint32(0); i < K; i++ {
// 		matrix[L+i] = make([]uint32, K)
// 		matrix[L+i][i] = 1
// 	}

// 	return matrix
// }

// // D = (I | P) and C = (-P // I) where P has dimension L x K
// func GenerateRandomColsOfD(L, K, p uint32, seed int64) [][]uint32 {
// 	matrix := make([][]uint32, L)

// 	for i := uint32(0); i < L; i++ {
// 		matrix[i] = make([]uint32, K)
// 	}

// 	rng := rand.New(rand.NewSource(seed))

// 	for j := uint32(0); j < K; j++ {
// 		for i := uint32(0); i < L; i++ {
// 			matrix[i][j] = uint32(rng.Intn(int(p)))
// 		}
// 	}

// 	return matrix
// }

// // C = (-P // I) Generate -P from C
// func GenerateRandomColsOfC(L, K, p uint32, seed int64) [][]uint32 {
// 	matrix := GenerateRandomColsOfD(L, K, p, seed)

// 	for j := uint32(0); j < K; j++ {
// 		for i := uint32(0); i < L; i++ {
// 			matrix[i][j] = (p - matrix[i][j]) % p
// 		}
// 	}

// 	return matrix
// }

// Return -P of C = (-P // I) flattened
func Generate1DDualMatrix(L, K uint32, field dataobjects.Field, seed int64) []uint32 {
	P := GenerateP(L, K, field, seed)

	vmatrix := make([]uint32, K*L)

	idx := 0

	for i := uint32(0); i < L; i++ {
		for j := uint32(0); j < K; j++ {
			vmatrix[idx] = field.Neg(P[i][j])
			idx += 1
		}
	}

	return vmatrix
}

func GenerateP(L, K uint32, field dataobjects.Field, seed int64) [][]uint32 {
	P := make([][]uint32, L)

	rng := rand.New(rand.NewSource(seed))

	for i := uint32(0); i < L; i++ {
		P[i] = make([]uint32, K)
		for j := uint32(0); j < K; j++ {
			P[i][j] = field.SampleElement(rng)
		}
	}

	return P
}

// Generate D = (I | P), transpose to D' = (I // P^T) and flatten
func Generate1DRLCMatrix(L, K uint32, p dataobjects.Field, seed int64) []uint32 {
	P := GenerateP(L, K, p, seed)

	vmatrix := make([]uint32, K*L)

	idx := 0

	for j := uint32(0); j < K; j++ {
		for i := uint32(0); i < L; i++ {
			vmatrix[idx] = P[i][j]
			idx += 1
		}
	}

	return vmatrix
}

// // D has dimension N x L
// func SampleVectorFromNullSpace(L, K, p uint32, seed int64) []uint32 {
// 	coeff := utils.RandomPrimeFieldVector(K, p)

// 	matrix := GenerateRandomColsOfC(L, K, p, seed)

// 	res := make([]uint32, K+L)

// 	for i := uint32(0); i < K; i++ {

// 		for j := uint32(0); j < L; j++ {
// 			res[j] = uint32((uint64(res[j]) + uint64(matrix[j][i])*uint64(coeff[i])) % uint64(p))
// 		}
// 		res[L+i] = uint32((uint64(res[L+i]) + uint64(coeff[i])) % uint64(p))

// 	}

// 	return res
// }
