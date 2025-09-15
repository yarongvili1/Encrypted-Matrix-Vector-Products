package pir

import (
	"RandomLinearCodePIR/utils"
	"math/rand"
)

type MixedSLSNPIR struct {
	Params MixedSLSNParams
}

type MixedSLSNParams struct {
	Rows           uint32
	Cols           uint32
	NumberOfBlocks uint32
	CodewordLength uint32
	PackedSize     uint32
}

type MixedSLSNPIRQuery struct {
	vec VectorF4
}

type MixedSLSNPIRAux struct {
	inv           VectorF4
	MaskValueBit1 uint32
	MaskValueBitP uint32
}

type MixedSLSNPIRAnswer struct {
	vec VectorF4
}

func (p *MixedSLSNPIR) KeyGen(N, ell, lambda int, seed int64) SecretKey {
	return SecretKey{
		LinearCodeKey: seed,
		// Hardcode for now
		MaskKey: seed + 1,
		Lambda:  lambda,
		Ell:     ell,
	}
}

func (p *MixedSLSNPIR) Encode(sk SecretKey, matrix MatrixF4) *MatrixF4 {
	M := p.Params.CodewordLength
	N := p.Params.Cols

	encodedMatrixBit1, encodedMatrixBitP := SystematicEncodingF4(M, sk.LinearCodeKey, matrix)

	// Transpose the encoded matrix for more efficient access pattern in C for XOR of rows instead of columns.
	packedMatrixBit1 := PackAndTransposeMatrix(encodedMatrixBit1, matrix.Rows, M)
	packedMatrixBitP := PackAndTransposeMatrix(encodedMatrixBitP, matrix.Rows, M)

	// Update the packed Data Length
	p.Params.PackedSize = (matrix.Rows + 31) / 32

	// Mask the matrix
	for i := uint32(0); i < p.Params.PackedSize; i++ {
		rng := rand.New(rand.NewSource(sk.MaskKey + int64(i)))
		index := i + N*p.Params.PackedSize
		for j := uint32(0); j < M-N; j++ {
			a := rng.Uint32()
			packedMatrixBit1[index] ^= a
			index += p.Params.PackedSize
		}
	}

	for i := uint32(0); i < p.Params.PackedSize; i++ {
		rng := rand.New(rand.NewSource(sk.MaskKey + int64(i+p.Params.PackedSize)))
		index := i + N*p.Params.PackedSize
		for j := uint32(0); j < M-N; j++ {
			b := rng.Uint32()
			packedMatrixBitP[index] ^= b
			index += p.Params.PackedSize
		}
	}

	packedMatrixBitSum := make([]uint32, len(packedMatrixBit1))
	for i := range packedMatrixBit1 {
		packedMatrixBitSum[i] = packedMatrixBit1[i] ^ packedMatrixBitP[i]
	}

	return &MatrixF4{
		Rows:      M,
		Cols:      p.Params.PackedSize,
		EntryBits: 32,
		Bit1:      packedMatrixBit1,
		BitP:      packedMatrixBitP,
		BitSum:    packedMatrixBitSum,
	}
}

func (p *MixedSLSNPIR) Query(sk SecretKey, queryIndex uint64) (*MixedSLSNPIRQuery, *MixedSLSNPIRAux) {
	queryVector := SampleVectorFromNullSpaceF4(p.Params.Cols, p.Params.CodewordLength, sk.LinearCodeKey)
	queryVector.Bit1[queryIndex%uint64(p.Params.Cols)] ^= 1

	// Calculate the mask for the final result
	rng := rand.New(rand.NewSource(int64(sk.MaskKey) + int64((queryIndex/uint64(p.Params.Cols))/32)))
	rng2 := rand.New(rand.NewSource(int64(sk.MaskKey) + int64((queryIndex/uint64(p.Params.Cols))/32) + int64(p.Params.PackedSize)))

	maskBit1 := uint32(0)
	maskBitP := uint32(0)

	for i := p.Params.Cols; i < p.Params.CodewordLength; i++ {
		a := rng.Uint32()        // random F4 element: a1 = a & 1, ap = (a >> 1) & 1
		b := rng2.Uint32()       // another independent element
		x := queryVector.Bit1[i] // x1
		y := queryVector.BitP[i] // xp

		ax := a * x // a1 * x1
		by := b * y // bp * xp
		ay := a * y // a1 * xp
		bx := b * x // bp * x1

		maskBit1 ^= (ax ^ by)      // low bit
		maskBitP ^= (bx ^ ay ^ by) // high bit
	}

	nonZeroCoeffBit1, nonZeroCoeffBitP, invBit1, invBitP := utils.RandomSplitLSNNoiseCoeffF4(p.Params.NumberOfBlocks)

	blockSize := p.Params.CodewordLength / p.Params.NumberOfBlocks

	for i := uint32(0); i < p.Params.NumberOfBlocks; i++ {
		blockStart := i * blockSize
		for j := blockStart; j < blockStart+blockSize; j++ {
			bit1 := (queryVector.Bit1[j] & nonZeroCoeffBit1[i]) ^ (queryVector.BitP[j] & nonZeroCoeffBitP[i])
			bitp := (queryVector.BitP[j] & nonZeroCoeffBit1[i]) ^ (queryVector.Bit1[j] & nonZeroCoeffBitP[i]) ^ (queryVector.BitP[j] & nonZeroCoeffBitP[i])
			queryVector.Bit1[j] = bit1
			queryVector.BitP[j] = bitp
			queryVector.BitSum[j] = queryVector.Bit1[j] ^ queryVector.BitP[j]
		}
	}

	return &MixedSLSNPIRQuery{
			vec: queryVector,
		}, &MixedSLSNPIRAux{
			inv:           VectorF4{Bit1: invBit1, BitP: invBitP},
			MaskValueBit1: maskBit1,
			MaskValueBitP: maskBitP,
		}
}

func (p *MixedSLSNPIR) Answer(matrix *MatrixF4, clientQuery *MixedSLSNPIRQuery) *MixedSLSNPIRAnswer {
	rows := matrix.Rows
	cols := matrix.Cols

	block_size := p.Params.CodewordLength / p.Params.NumberOfBlocks

	// (x + yp) * (a + bp) = (a * x + b * y) + [(a + b) * (x + y) - a * x]
	vec_ax := make([]uint32, cols*rows/block_size)
	vec_by := make([]uint32, cols*rows/block_size)
	vec_abxy := make([]uint32, cols*rows/block_size)

	MatrixColXORByBlock(clientQuery.vec.Bit1, matrix.Bit1, vec_ax, rows, cols, block_size)
	MatrixColXORByBlock(clientQuery.vec.BitP, matrix.BitP, vec_by, rows, cols, block_size)
	MatrixColXORByBlock(clientQuery.vec.BitSum, matrix.BitSum, vec_abxy, rows, cols, block_size)

	for i := range vec_ax {
		vec_by[i] ^= vec_ax[i]
		vec_abxy[i] ^= vec_ax[i]
	}

	return &MixedSLSNPIRAnswer{vec: VectorF4{Cols: cols * rows / block_size, Bit1: vec_by, BitP: vec_abxy}}
}

func (p *MixedSLSNPIR) Decode(sk SecretKey, index uint64, response *MixedSLSNPIRAnswer, aux *MixedSLSNPIRAux) (uint32, uint32) {
	row := index / uint64(p.Params.Cols)
	wordIndex := row / 32
	bitOffset := row % 32

	bit1 := uint32(0)
	bitP := uint32(0)

	for i := uint32(0); i < p.Params.NumberOfBlocks; i++ {
		offset := i*((p.Params.Rows+31)/32) + uint32(wordIndex)
		res_bit1 := response.vec.Bit1[offset] >> bitOffset & 1
		res_bitP := response.vec.BitP[offset] >> bitOffset & 1
		bit1 ^= (res_bit1 & aux.inv.Bit1[i]) ^ (res_bitP & aux.inv.BitP[i])
		bitP ^= (res_bit1 & aux.inv.BitP[i]) ^ (res_bitP & aux.inv.Bit1[i]) ^ (res_bitP & aux.inv.BitP[i])
	}

	valBit1 := bit1 ^ (aux.MaskValueBit1 >> bitOffset & 1)
	valBitP := bitP ^ (aux.MaskValueBitP >> bitOffset & 1)

	return valBit1, valBitP
}
