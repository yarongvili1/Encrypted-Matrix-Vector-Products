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
