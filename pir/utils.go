package pir

import "math/rand"

// =========== Random Vectors ===========

func RandomizeQueryVector(N uint32, i uint64) []uint32 {
	return RandomizeBinaryVector(N)
}

func RandomizeBinaryVector(N uint32) []uint32 {
	vector := make([]uint32, N)
	for i := range vector {
		vector[i] = uint32(rand.Intn(2)) // Generates either 0 or 1
	}
	return vector
}

// Generator vector that contains values in {0,1,2}
func RandomizeFlipVector(N uint32) []uint32 {
	vector := make([]uint32, N)
	for i := range vector {
		vector[i] = uint32(rand.Intn(3)) // Generates either 0 or 1
	}
	return vector
}

func RandomSplitLSNNoiseCoeff(s uint32, p uint32) []uint32 {
	vector := make([]uint32, s)
	for i := range vector {
		vector[i] = uint32(rand.Intn(int(p)-1) + 1) // Generates non-zero values in F_p
	}
	return vector
}

func RandomPrimeFieldVector(n uint32, p uint32) []uint32 {
	vector := make([]uint32, n)
	for i := range vector {
		vector[i] = uint32(rand.Intn(int(p))) // Generates non-zero values in F_p
	}
	return vector
}

func RandomNoiseVector(n uint32, epsi float32, p uint32) []uint32 {
	vector := make([]uint32, n)
	for i := range vector {
		if rand.Float32() <= epsi {
			vector[i] = uint32(rand.Intn(int(p-1))) + 1 // Generates non-zero values in F_p
		}
	}
	return vector
}

func IsZeroVector(v []uint32) bool {
	for _, val := range v {
		if val != 0 {
			return false
		}
	}
	return true
}

func ModInverse(a, p uint32) uint32 {
	var t, newT int64 = 0, 1
	var r, newR int64 = int64(p), int64(a)

	for newR != 0 {
		quotient := r / newR
		t, newT = newT, t-quotient*newT
		r, newR = newR, r-quotient*newR
	}

	if r > 1 {
		panic("a is not invertible")
	}

	if t < 0 {
		t += int64(p)
	}

	return uint32(t)
}
