package lpnpir

import (
	"math/rand"
)

// The S in this file is S^T in the paper
// S has dimension L x K
func GenerateS(L, K, p uint32, seed int64) [][]uint32 {
	S := make([][]uint32, L)

	rng := rand.New(rand.NewSource(seed))

	for i := uint32(0); i < L; i++ {
		S[i] = make([]uint32, K)
		for j := uint32(0); j < K; j++ {
			S[i][j] = uint32(rng.Intn(int(p)))
		}
	}

	return S
}

// Assume encoding matrix D = (I | -S)
// We can encode message use m x D^T where D^T = (I // -S^T)
// Here we return -S^T by row, which is -S by column
func Generate1DEncodingMatrix(L, K, p uint32, seed int64) []uint32 {
	S := GenerateS(L, K, p, seed)

	vS := make([]uint32, K*L)

	idx := 0

	// Negate the values in S and
	for j := uint32(0); j < K; j++ {
		for i := uint32(0); i < L; i++ {
			vS[idx] = (p - S[i][j]) % p
			idx += 1
		}
	}

	return vS
}

// C = (S // I), S has dimension L x K
func Generate1DLSNMatrix(L, K, p uint32, seed int64) []uint32 {
	S := GenerateS(L, K, p, seed)

	vS := make([]uint32, K*L)

	idx := 0

	for i := uint32(0); i < L; i++ {
		for j := uint32(0); j < K; j++ {
			vS[idx] = S[i][j]
			idx += 1
		}
	}

	return vS
}
