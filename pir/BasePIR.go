package pir

import (
	"RandomLinearCodePIR/utils"
	"math/rand"
)

type BasePIR struct {
	Params BaseParams
}

type BaseParams struct {
	Rows           uint32
	Cols           uint32
	NumberOfBlocks uint32
	CodewordLength uint32
	PackedSize     uint32
}

type BasePIRQuery struct {
	Vector_1 []uint32
	Vector_2 []uint32
}

type BasePIRAux struct {
	FlipVector []uint32
	MaskValue  uint32
}
type BasePIRAnswer struct {
	Result_1 []uint32
	Result_2 []uint32
}

func (p *BasePIR) KeyGen(N, ell, lambda int, seed int64) SecretKey {
	return SecretKey{
		LinearCodeKey: seed,
		// Hardcode for now
		MaskKey: seed + 1,
		Lambda:  lambda,
		Ell:     ell,
	}
}

func (p *BasePIR) Encode(sk SecretKey, matrix Matrix) *Matrix {
	M := p.Params.CodewordLength

	encodedMatrix := SystematicEncoding(M, sk.LinearCodeKey, matrix)

	// Transpose the encoded matrix for more efficient access pattern in C for XOR of rows instead of columns.
	packedData := PackAndTransposeMatrix(encodedMatrix, matrix.Rows, M)

	// Update the packed Data Length
	p.Params.PackedSize = (matrix.Rows + 31) / 32

	// Mask the packed matrix column wisely
	for i := uint32(0); i < p.Params.PackedSize; i++ {
		rng := rand.New(rand.NewSource(sk.MaskKey + int64(i)))
		index := i
		for j := uint32(0); j < M; j++ {
			packedData[index] ^= rng.Uint32()
			index += p.Params.PackedSize
		}
	}

	return &Matrix{
		Rows:      M,
		Cols:      p.Params.PackedSize,
		EntryBits: 32,
		Data:      packedData,
	}
}

func (p *BasePIR) Query(sk SecretKey, queryIndex uint64) (*BasePIRQuery, *BasePIRAux) {
	queryVector := SampleVectorFromNullSpace(p.Params.Cols, p.Params.CodewordLength, sk.LinearCodeKey)

	// Add Unit Vector to retrieve the ith column
	queryVector[queryIndex%uint64(p.Params.Cols)] ^= 1

	// Calculate the mask for the final result
	rng := rand.New(rand.NewSource(int64(sk.MaskKey) + int64((queryIndex/uint64(p.Params.Cols))/32)))

	mask := uint32(0)

	for i := range queryVector {
		mask ^= queryVector[i] * rng.Uint32()
	}

	flipVector := utils.RandomizeFlipVector(p.Params.NumberOfBlocks)
	vec_1, vec_2 := makeVectors(queryVector, flipVector, len(queryVector), int(p.Params.NumberOfBlocks))

	return &BasePIRQuery{
			Vector_1: vec_1,
			Vector_2: vec_2,
		}, &BasePIRAux{
			FlipVector: flipVector,
			MaskValue:  mask,
		}
}

// Generate vec_1, vec_2 from queryVector and flipVector, globally packed.
func makeVectors(queryVector []uint32, flipVector []uint32,
	codewordLength, nBlocks int) (vec1, vec2 []uint32) {

	blockSize := codewordLength / nBlocks
	vec1 = utils.PackBinaryVector(queryVector)
	vec2 = make([]uint32, (codewordLength+31)>>5)

	for b := 0; b < nBlocks; b++ {
		blockStart := b * blockSize
		blockEnd := blockStart + blockSize
		if blockEnd > codewordLength {
			blockEnd = codewordLength
		}

		wStart := blockStart >> 5
		wEnd := (blockEnd - 1) >> 5

		for w := wStart; w <= wEnd; w++ {
			base := w << 5
			lo := 0
			if blockStart > base {
				lo = blockStart - base
			}
			hi := 32
			if blockEnd-base < 32 {
				hi = blockEnd - base
			}
			mask := utils.WordMask(lo, hi)

			switch flipVector[b] {
			case 0:
				// vec2 = random, vec1 unchanged
				vec2[w] = (vec2[w] & ^mask) | (rand.Uint32() & mask)

			case 1:
				// vec2 ^= vec1, then vec1 = random
				vec2[w] ^= vec1[w] & mask
				vec1[w] = (vec1[w] & ^mask) | (rand.Uint32() & mask)

			default:
				// vec2 = random, then vec1 ^= vec2
				r := rand.Uint32() & mask
				vec2[w] = (vec2[w] & ^mask) | r
				vec1[w] ^= r
			}
		}
	}
	return
}

func (p *BasePIR) Answer(matrix *Matrix, clientQuery *BasePIRQuery) *BasePIRAnswer {
	rows := matrix.Rows
	cols := matrix.Cols

	// TODO : Make sure it devides
	block_size := p.Params.CodewordLength / p.Params.NumberOfBlocks

	result1 := make([]uint32, cols*rows/block_size)
	result2 := make([]uint32, cols*rows/block_size)

	MatrixColXORByBlock2D(clientQuery.Vector_1, clientQuery.Vector_2, matrix.Data, result1, result2, rows, cols, block_size)

	return &BasePIRAnswer{
		Result_1: result1,
		Result_2: result2,
	}
}

func (p *BasePIR) Decode(sk SecretKey, index uint64, response *BasePIRAnswer, aux *BasePIRAux) uint32 {
	res := make([]uint32, p.Params.PackedSize)
	row := index / uint64(p.Params.Cols)
	wordIndex := row / 32
	bitOffset := row % 32

	start := 0
	for _, val := range aux.FlipVector {
		if val != 0 {
			for i := range res {
				res[i] ^= response.Result_2[start+i]
			}

		}

		if val != 1 {
			for i := range res {
				res[i] ^= response.Result_1[start+i]
			}
		}

		start += int(p.Params.PackedSize)
	}

	val := res[wordIndex] ^ aux.MaskValue

	return (val >> bitOffset) & 1
}
