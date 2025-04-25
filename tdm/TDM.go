package tdm

/*
#cgo CFLAGS: -I../TDM
#cgo LDFLAGS: -L../TDM -L/opt/homebrew/lib -lNTT -lntl -lgmp -lstdc++
#include "NTT.h"
*/
import "C"
import (
	"RandomLinearCodePIR/utils"
	"math"
	"math/rand"
)

type TDM struct {
	M     uint32
	N     uint32
	Q     uint32
	SeedL int64
	SeedP int64
	SeedR int64
	// Internal Use
	m      uint32
	n      uint32
	rootL  uint32
	rootR  uint32
	blockL uint32
	blockR uint32
	permD  uint32
}

func (td *TDM) GenerateTrapDooredMatrix(seedL, seedP, seedR int64) [][]uint32 {
	td.updateInternalUseParams()
	Q_R := GetQuasiCyclicMatrix(td.permD, td.n, td.blockR, td.Q, seedR)

	perm := GetPermutation(td.permD, seedP)
	PermuteRowsInPlace(Q_R, perm)

	// R = Q_L x perm(Q_R)
	R := make([][]uint32, td.m)
	for i := range R {
		R[i] = make([]uint32, td.n)
	}

	for j := uint32(0); j < td.n; j++ {
		v := make([]uint32, td.permD)
		for i := range Q_R {
			v[i] = Q_R[i][j]
		}
		res := QuasiCyclicVectorMul(td.m, td.permD, td.blockL, td.Q, td.rootL, seedL, v)
		for i := range res {
			R[i][j] = res[i]
		}
	}

	return R
}

func (td *TDM) GenerateFlattenedTrapDooredMatrix() []uint32 {
	result := make([]uint32, td.M*td.N)
	R := td.GenerateTrapDooredMatrix(td.SeedL, td.SeedP, td.SeedR)

	// Only return the upper-left cornor of the TDM
	for i := uint32(0); i < td.M; i++ {
		copy(result[i*td.N:(i+1)*td.N], R[i])
	}
	return result
}

func (td *TDM) GenerateFlattenedTrapDooredMatrixPerSlice(sliceNum int64) []uint32 {
	result := make([]uint32, td.M*td.N)
	R := td.GenerateTrapDooredMatrix(td.SeedL+sliceNum, td.SeedP+sliceNum, td.SeedR+sliceNum)

	// Only return the upper-left cornor of the TDM
	for i := uint32(0); i < td.M; i++ {
		copy(result[i*td.N:(i+1)*td.N], R[i])
	}
	return result
}

func (td *TDM) EvaluationCircuit(v []uint32) []uint32 {
	if td.m == 0 {
		td.updateInternalUseParams()
	}

	if int(td.n) > len(v) {
		padded := make([]uint32, td.n)
		copy(padded, v)
		v = padded
	}

	vec := QuasiCyclicVectorMul(td.permD, td.n, td.blockR, td.Q, td.rootR, td.SeedR, v)
	perm := GetPermutation(td.permD, td.SeedP)
	PermuteVectorInPlace(vec, perm)
	masks := QuasiCyclicVectorMul(td.m, td.permD, td.blockL, td.Q, td.rootL, td.SeedL, vec)

	return masks[0:td.M]
}

func (td *TDM) EvaluationCircuitPerSlice(v []uint32, i int64) []uint32 {
	// Pad with 0's if TDM has more columns
	if int(td.n) > len(v) {
		padded := make([]uint32, td.n)
		copy(padded, v)
		v = padded
	}

	vec := QuasiCyclicVectorMul(td.permD, td.n, td.blockR, td.Q, td.rootR, td.SeedR+i, v)
	perm := GetPermutation(td.permD, td.SeedP+i)
	PermuteVectorInPlace(vec, perm)
	masks := QuasiCyclicVectorMul(td.m, td.permD, td.blockL, td.Q, td.rootL, td.SeedL+i, vec)

	// Trim the padded rows of TDM
	return masks[0:td.M]
}

func QuasiCyclicVectorMul(row, col, blockSize, q, root uint32, seed int64, v []uint32) []uint32 {
	result := make([]uint32, row)
	tmp_result := make([]uint32, blockSize)

	polyQC := make([]uint32, blockSize)

	for i := uint32(0); i < row/blockSize; i++ {
		for j := uint32(0); j < col/blockSize; j++ {
			seed_ij := seed + int64(i*row/blockSize+j)
			rng := rand.New(rand.NewSource(seed_ij))
			polyQC[0] = uint32(rng.Intn(int(q)))
			for t := uint32(1); t < blockSize; t++ {
				polyQC[blockSize-t] = uint32(rng.Intn(int(q)))
			}
			NTT_Convolution(polyQC, v[j*blockSize:(j+1)*blockSize], tmp_result, blockSize, root, q)

			for t := uint32(0); t < blockSize; t++ {
				result[i*blockSize+t] = uint32((uint64(result[i*blockSize+t]) + uint64(tmp_result[t])) % uint64(q))
			}
		}
	}

	return result
}

func PermuteRowsInPlace(matrix [][]uint32, perm []uint32) {
	n := uint32(len(matrix))
	visited := make([]bool, n)

	for i := uint32(0); i < n; i++ {
		if visited[i] || perm[i] == i {
			continue
		}

		j := i
		temp := matrix[i]
		for {
			next := perm[j]
			visited[j] = true

			if next == i {
				matrix[j] = temp
				break
			}

			matrix[j] = matrix[next]
			j = next
		}
	}
}

func PermuteVectorInPlace(vec []uint32, perm []uint32) {
	n := uint32(len(vec))
	visited := make([]bool, n)

	for i := uint32(0); i < n; i++ {
		if visited[i] || perm[i] == i {
			continue
		}

		j := i
		temp := vec[i]
		for {
			next := perm[j]
			visited[j] = true

			if next == i {
				vec[j] = temp
				break
			}

			vec[j] = vec[next]
			j = next
		}
	}
}

func GetPermutation(n uint32, seed int64) []uint32 {
	rng := rand.New(rand.NewSource(seed))
	perm := make([]uint32, n)
	for i := uint32(0); i < n; i++ {
		perm[i] = i
	}
	rng.Shuffle(int(n), func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})

	return perm
}

func GetQuasiCyclicMatrix(row, col, blockSize, q uint32, seed int64) [][]uint32 {
	Q := make([][]uint32, row)
	for i := uint32(0); i < row; i++ {
		Q[i] = make([]uint32, col)
	}

	for i := uint32(0); i < row/blockSize; i++ {
		for j := uint32(0); j < col/blockSize; j++ {
			seed_ij := seed + int64(i*row/blockSize+j)
			rng := rand.New(rand.NewSource(seed_ij))
			poly := make([]uint32, blockSize)
			for t := uint32(0); t < blockSize; t++ {
				poly[t] = uint32(rng.Intn(int(q)))
			}

			for t := uint32(0); t < blockSize; t++ {
				copy(Q[i*blockSize+t][j*blockSize+t:(j+1)*blockSize], poly[0:blockSize-t])
				copy(Q[i*blockSize+t][j*blockSize:j*blockSize+t], poly[blockSize-t:blockSize])
			}

		}
	}

	return Q
}

func (td *TDM) updateInternalUseParams() {
	td.permD = max(roundUpToPowerOf2(td.M), roundUpToPowerOf2(td.N))

	td.blockL = td.determineBlockSize(td.M, td.permD)
	td.blockR = td.determineBlockSize(td.N, td.permD)

	td.m = utils.RoundUp(td.M, td.blockL)
	td.n = utils.RoundUp(td.N, td.blockR)
	td.rootL = NthRootOfUnity(td.Q, td.blockL)
	td.rootR = NthRootOfUnity(td.Q, td.blockR)
}

func roundUpToPowerOf2(m uint32) uint32 {
	return uint32(1) << uint32(math.Ceil(math.Log2(float64(m))))
}

// Currently hardcode it to be max(2^(ceil(log2(min(m,n)))), (q-1)/2) for m x n matrix
// TODO: update for m,n,q in general
func (td *TDM) determineBlockSize(m, n uint32) uint32 {
	minOfMN := min(m, n)
	if minOfMN >= (td.Q-1)/2 {
		return (td.Q - 1) / 2
	}

	return uint32(1) << uint32(math.Ceil(math.Log2(float64(minOfMN))))
}
