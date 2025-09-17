package pir

import (
	"RandomLinearCodePIR/dataobjects"
	"math"
	"math/rand"
)

type Matrix struct {
	Rows      uint32
	Cols      uint32
	EntryBits uint32
	Data      []uint32
}

type MatrixF4 struct {
	Rows      uint32
	Cols      uint32
	EntryBits uint32
	Bit1      []uint32 // length: Rows * Cols, coeff of 1
	BitP      []uint32 // length: Rows * Cols, coeff of p
	BitSum    []uint32 // Bit1 + Bitp
}

type VectorF4 struct {
	Cols   uint32
	Bit1   []uint32
	BitP   []uint32
	BitSum []uint32 // Bit1 + Bitp
}

func GenerateMatrix(rows, cols, entryBits uint32, seed int64) Matrix {
	rng := rand.New(rand.NewSource(seed))

	dataSize := uint64(rows) * uint64(cols)

	data := dataobjects.AlignedMake[uint32](uint64(dataSize))

	bitmask := uint32((1 << entryBits) - 1)

	for i := range data {
		data[i] = rng.Uint32() & bitmask
	}

	return Matrix{
		Rows:      rows,
		Cols:      cols,
		EntryBits: entryBits,
		Data:      data,
	}
}

func GenerateMatrixF4(rows, cols, entryBits uint32, seed int64) MatrixF4 {
	rng := rand.New(rand.NewSource(seed))

	dataSize := uint64(rows) * uint64(cols)

	databit1 := make([]uint32, dataSize)
	databitp := make([]uint32, dataSize)
	databitsum := make([]uint32, dataSize)

	bitmask := uint32((1 << entryBits) - 1)

	for i := range databit1 {
		databit1[i] = rng.Uint32() & bitmask
		databitp[i] = rng.Uint32() & bitmask
		databitsum[i] = databit1[i] ^ databitp[i]
	}

	return MatrixF4{
		Rows:      rows,
		Cols:      cols,
		EntryBits: entryBits,
		Bit1:      databit1,
		BitP:      databitp,
		BitSum:    databitsum,
	}
}

func GenerateRandomMatrix(rows, cols uint32, seed int64) [][]uint32 {
	rng := rand.New(rand.NewSource(seed)) // PRNG with fixed seed
	matrix := make([][]uint32, rows)

	for i := uint32(0); i < rows; i++ {
		matrix[i] = dataobjects.AlignedMake[uint32](uint64(cols))
		for j := uint32(0); j < cols; j++ {
			matrix[i][j] = rng.Uint32() % 2 // Random F₂ value (0 or 1)
		}
	}

	return matrix
}

func GeneratePrimeFieldMatrix(rows, cols, p uint32, seed int64) Matrix {
	rng := rand.New(rand.NewSource(seed))

	dataSize := uint64(rows) * uint64(cols)

	data := dataobjects.AlignedMake[uint32](uint64(dataSize))

	for i := range data {
		data[i] = uint32(rng.Intn(int(p)))
	}

	return Matrix{
		Rows:      rows,
		Cols:      cols,
		EntryBits: uint32(math.Ceil(math.Log2(float64(p)))),
		Data:      data,
	}
}

func PackAndTransposeMatrix(matrix [][]uint32, rows, cols uint32) []uint32 {
	packedSize := (rows + 31) / 32 * cols // Each column stores 32 rows
	packedData := dataobjects.AlignedMake[uint32](uint64(packedSize))

	// Column-major packing (transpose)
	for col := uint32(0); col < cols; col++ {
		for row := uint32(0); row < rows; row++ {
			index := col*((rows+31)/32) + row/32 // Compute index in packed data, which is stored column by column
			offset := row % 32                   // Bit position in uint32
			if matrix[row][col] == 1 {
				packedData[index] |= (1 << offset) // Set bit
			}
		}
	}

	return packedData
}
