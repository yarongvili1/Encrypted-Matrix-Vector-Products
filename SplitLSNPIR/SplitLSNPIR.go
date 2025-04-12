package splitlsnpir

/*
#include "PrimeFieldMVP.h"
#include "RandomLinearCode.h"
*/
import "C"
import (
	"RandomLinearCodePIR/pir"
	"unsafe"
)

type SlsnPIR struct {
	Params SlsnParams
}

type SecretKey struct {
	LinearCodeKey   int64
	TDMKey          int64
	PreLoadedMatrix []uint32
}

// N = K + L denotes the length of the codeword
// Encoding Matrix C with dimension N x K
// Original Data Matrix has dimension M x N
type SlsnParams struct {
	PrimeField uint32
	NumBlocks  uint32
	K          uint32
	L          uint32
	N          uint32
	M          uint32
}

type SPIRQuery struct {
	vec []uint32
}

type SPIRAux struct {
	Coeff []uint32
}

// C in dimension K x N with form (I | P), P with dimension K x L
// D in dimension L x N with form (-P//I)
func (spir *SlsnPIR) KeyGen(seed int64) SecretKey {
	params := spir.Params
	return SecretKey{
		LinearCodeKey:   seed,
		TDMKey:          seed + 1,
		PreLoadedMatrix: Generate1DDualSpaceRandomMatrix(params.K, params.L, params.PrimeField, seed),
	}
}

func (spir *SlsnPIR) Encode(sk SecretKey, matrix pir.Matrix) *pir.Matrix {
	params := spir.Params
	rlcMatrix := Generate1DRLCMatrix(params.K, params.L, params.PrimeField, sk.LinearCodeKey)

	N := params.K + params.L
	data := make([]uint32, matrix.Rows*N)

	for i := uint32(0); i < matrix.Rows; i++ {
		for j := uint32(0); j < params.K; j++ {
			data[i*N+j] = matrix.Data[i*params.K+j]
		}

		C.BlockMatrixVectorProduct(
			(*C.uint32_t)(unsafe.Pointer(&rlcMatrix[0])),
			(*C.uint32_t)(unsafe.Pointer(&matrix.Data[i*matrix.Cols])),
			(*C.uint32_t)(unsafe.Pointer(&data[i*N+params.K])),
			C.uint32_t(params.K), C.uint32_t(params.L), C.uint32_t(1), C.uint32_t(params.PrimeField))
	}

	// Transpose to be easy for Row linear combination
	transposedData := make([]uint32, N*matrix.Rows)
	idx := 0
	for i := uint32(0); i < N; i++ {
		for j := uint32(0); j < matrix.Rows; j++ {
			transposedData[idx] = data[j*N+i]
			idx += 1
		}
	}

	return &pir.Matrix{
		Rows: N,
		Cols: matrix.Rows,
		Data: transposedData,
	}
}

func (spir *SlsnPIR) Query(sk SecretKey, queryIndex uint64) ([]uint32, []uint32) {
	params := spir.Params

	// Sample Vector From Nullspace
	PofDual := sk.PreLoadedMatrix
	if len(PofDual) == 0 {
		PofDual = Generate1DDualSpaceRandomMatrix(params.K, params.L, params.PrimeField, sk.LinearCodeKey)
	}

	nullspaceCoeff := pir.RandomPrimeFieldVector(params.L, params.PrimeField)

	queryVector := make([]uint32, params.K+params.L)

	C.BlockMatrixVectorProduct(
		(*C.uint32_t)(unsafe.Pointer(&PofDual[0])),
		(*C.uint32_t)(unsafe.Pointer(&nullspaceCoeff[0])),
		(*C.uint32_t)(unsafe.Pointer(&queryVector[0])),
		C.uint32_t(params.L), C.uint32_t(params.K), C.uint32_t(1), C.uint32_t(params.PrimeField))

	for i := uint32(0); i < params.L; i++ {
		queryVector[params.K+i] = nullspaceCoeff[i]
	}

	// Add queryIndex bit by 1
	queryVector[queryIndex%uint64(params.K)] = (queryVector[queryIndex%uint64(params.K)] + 1) % params.PrimeField

	// For Each Block multiply by a non-zero scalar
	coeff := pir.RandomSplitLSNNoiseCoeff(params.NumBlocks, params.PrimeField)

	blkSize := (params.K + params.L) / params.NumBlocks

	for i := uint32(0); i < params.NumBlocks; i++ {
		for j := uint32(0); j < blkSize; j++ {
			queryVector[i*blkSize+j] = uint32((uint64(queryVector[i*blkSize+j]) * uint64(coeff[i])) % uint64(params.PrimeField))
		}
	}

	return queryVector, coeff
}

func (spir *SlsnPIR) Answer(matrix pir.Matrix, clientQuery []uint32) []uint32 {
	n := matrix.Rows
	m := matrix.Cols
	s := spir.Params.NumBlocks
	p := spir.Params.PrimeField
	result := make([]uint32, s*m)

	C.BlockMatrixVectorProduct(
		(*C.uint32_t)(unsafe.Pointer(&matrix.Data[0])),
		(*C.uint32_t)(unsafe.Pointer(&clientQuery[0])),
		(*C.uint32_t)(unsafe.Pointer(&result[0])),
		C.uint32_t(n), C.uint32_t(m), C.uint32_t(s), C.uint32_t(p))

	return result
}

func (spir *SlsnPIR) Decode(sk SecretKey, index uint64, response []uint32, aux []uint32) uint32 {
	params := spir.Params

	vec := make([]uint32, spir.Params.NumBlocks)

	for i := range aux {
		vec[i] = pir.ModInverse(aux[i], spir.Params.PrimeField)
	}

	result := make([]uint32, spir.Params.M)

	C.BlockMatrixVectorProduct(
		(*C.uint32_t)(unsafe.Pointer(&response[0])),
		(*C.uint32_t)(unsafe.Pointer(&vec[0])),
		(*C.uint32_t)(unsafe.Pointer(&result[0])),
		C.uint32_t(spir.Params.NumBlocks), C.uint32_t(params.M), C.uint32_t(1), C.uint32_t(params.PrimeField))

	return result[index/uint64(params.K)]
}
