package pir

import (
	"RandomLinearCodePIR/dataobjects"
	"RandomLinearCodePIR/utils"
	"math/rand"
)

// The functions in this file are hardcoded for F_2

// Generate Sysmtematic Random Linear Code which has form G=(I_N | P) where P has dimension N * (M-N)
// Random each column of P using seed = seed + col where col := range 1 to M-N
func GenerateRandomLinearCode(N, M uint32, seed int64) {

}

func LinearizeMatrixByRows(rows, cols uint32, matrix [][]uint32) []uint32 {
	matrix1D := make([]uint32, rows*cols)
	for i := uint32(0); i < rows; i++ {
		copy(matrix1D[i*cols:(i+1)*cols], matrix[i])
	}

	return matrix1D
}

func LinearizeMatrixByCols(rows, cols uint32, matrix [][]uint32) []uint32 {
	matrix1D := make([]uint32, rows*cols)
	index := 0

	for col := uint32(0); col < cols; col++ {
		for row := uint32(0); row < rows; row++ {
			matrix1D[index] = matrix[row][col]
			index++
		}
	}

	return matrix1D
}

func GenerateRandomColsOfRLC(N, M uint32, seed int64) [][]uint32 {
	if M < N {
		panic("Codeword length should be longer than Message length.")
	}

	matrix := make([][]uint32, N)

	for i := uint32(0); i < N; i++ {
		matrix[i] = dataobjects.AlignedMake[uint32](uint64(M - N))
	}

	for j := uint32(0); j < M-N; j++ {
		rng := rand.New(rand.NewSource(seed + int64(j)))
		for i := uint32(0); i < N; i++ {
			matrix[i][j] = rng.Uint32() % 2
		}
	}

	return matrix
}

func SystematicEncoding(M uint32, seed int64, matrix Matrix) [][]uint32 {
	N := matrix.Cols
	RandomColsOfRLC := GenerateRandomColsOfRLC(N, M, seed)
	RLC1D := LinearizeMatrixByRows(N, M-N, RandomColsOfRLC)

	encodedMatrix := make([][]uint32, matrix.Rows)
	for i := range encodedMatrix {
		encodedMatrix[i] = dataobjects.AlignedMake[uint32](uint64(M))
	}

	for row := uint32(0); row < matrix.Rows; row++ {
		start := uint64(row) * uint64(matrix.Cols)
		end := start + uint64(matrix.Cols)
		originalRow := matrix.Data[start:end]

		// Copy systematic part (first N bits)
		copy(encodedMatrix[row][:N], originalRow[:N])

		VecMatrixMulF2(encodedMatrix[row][N:M], RLC1D, originalRow[:N], N, M-N)
	}

	return encodedMatrix
}

func SystematicEncodingF4(M uint32, seed int64, matrix MatrixF4) ([][]uint32, [][]uint32) {
	N := matrix.Cols
	RandomColsOfRLCBit1 := GenerateRandomColsOfRLC(N, M, seed)
	RLC1DBit1 := LinearizeMatrixByRows(N, M-N, RandomColsOfRLCBit1)

	RandomColsOfRLCBitP := GenerateRandomColsOfRLC(N, M, seed+10)
	RLC1DBitP := LinearizeMatrixByRows(N, M-N, RandomColsOfRLCBitP)

	encodedMatrixBit1 := make([][]uint32, matrix.Rows)
	encodedMatrixBitP := make([][]uint32, matrix.Rows)

	for i := range encodedMatrixBit1 {
		encodedMatrixBit1[i] = make([]uint32, M)
		encodedMatrixBitP[i] = make([]uint32, M)
	}

	for row := uint32(0); row < matrix.Rows; row++ {
		start := uint64(row) * uint64(matrix.Cols)
		end := start + uint64(matrix.Cols)
		originalRowBit1 := matrix.Bit1[start:end]
		originalRowBitP := matrix.BitP[start:end]

		// Copy systematic part (first N bits)
		copy(encodedMatrixBit1[row][:N], originalRowBit1[:N])
		copy(encodedMatrixBitP[row][:N], originalRowBitP[:N])

		//(a + bp) * (x + yp) = (a * x + b * y) + (b * x + a * y + b * y)p

		bit1 := make([]uint32, M-N)
		bitP := make([]uint32, M-N)

		VecMatMulF4(bit1, bitP, RLC1DBit1, RLC1DBitP, originalRowBit1[:N], originalRowBitP[:N], N, M-N)

		copy(encodedMatrixBit1[row][N:], bit1)
		copy(encodedMatrixBitP[row][N:], bitP)

	}
	return encodedMatrixBit1, encodedMatrixBitP
}

func SampleVectorFromNullSpaceF4(N, M uint32, seed int64) VectorF4 {
	coeffBit1 := utils.RandomizeBinaryVector(M - N)
	coeffBitP := utils.RandomizeBinaryVector(M - N)

	RandomColsOfRLCBit1 := GenerateRandomColsOfRLC(N, M, seed)
	RLC1DBit1 := LinearizeMatrixByCols(N, M-N, RandomColsOfRLCBit1)

	RandomColsOfRLCBitP := GenerateRandomColsOfRLC(N, M, seed+10)
	RLC1DBitP := LinearizeMatrixByCols(N, M-N, RandomColsOfRLCBitP)

	bit1 := make([]uint32, M)
	bitP := make([]uint32, M)
	bitSum := make([]uint32, M)

	VecMatMulF4(bit1, bitP, RLC1DBit1, RLC1DBitP, coeffBit1, coeffBitP, M-N, N)
	copy(bit1[N:M], coeffBit1)
	copy(bitP[N:M], coeffBitP)

	for i := range bit1 {
		bitSum[i] = bit1[i] ^ bitP[i]
	}

	return VectorF4{
		Cols:   M,
		Bit1:   bit1,
		BitP:   bitP,
		BitSum: bitSum,
	}

}

// The Parity check matrix has the form H = vcat(P, I_(M-N))
// We sample a vector of length M-N in F2 to be the coefficients of the linear combination of the columns
// We can do XOR of the columns while we know the column i is composed by the ith column of P and the ith unit vector
func SampleVectorFromNullSpace(N, M uint32, seed int64) []uint32 {
	coeff := utils.RandomizeBinaryVector(M - N)
	res := dataobjects.AlignedMake[uint32](uint64(M))

	for i := uint32(0); i < M-N; i++ {
		if coeff[i] == 1 {
			rng := rand.New(rand.NewSource(seed + int64(i)))
			for j := uint32(0); j < N; j++ {
				res[j] ^= rng.Uint32() % 2
			}
			res[N+i] ^= 1
		}
	}

	return res
}
