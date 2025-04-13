package lpnpir

/*
#cgo CFLAGS: -I../LPNPIR
#cgo LDFLAGS: -L../LPNPIR -L/opt/homebrew/lib -lMVP -lntl -lgmp -lstdc++
#include "MVP.h"
*/
import "C"
import (
	"RandomLinearCodePIR/pir"
	"unsafe"
)

type LpnPIR struct {
	Params LpnParams
}

type LpnParams struct {
	PrimeField uint32
	Lambda     uint32
	K          uint32
	Epsi       float32
	N          uint32
	M          uint32
	L          uint32
	M_1        uint32
	ECCLength  uint32
}

type SecretKey struct {
	LinearCodeKey   int64
	TDMKey          int64
	PreLoadedMatrix []uint32
}

// TODO: Check if it is stable to have queries as [][]uint32
type LpnPIRQuery struct {
	Queries      []uint32
	QueryLen     uint32
	NumOfQueries uint32
}

// TODO: Check if it is stable to have queries as [][]uint32
type LpnPIRAux struct {
	EC            []uint32
	ECLen         uint32
	NoisyQueryInd []bool
}

type LpnPIRAnswer struct {
	// Concatinate of ECCLength Answer, each with length M/M_1
	Answers []uint32
	AnsLen  uint32
}

func (pi *LpnPIR) KeyGen(seed int64) SecretKey {
	// TODO Add the calculations for the params from protocol
	params := pi.Params
	return SecretKey{
		LinearCodeKey:   seed,
		TDMKey:          seed + 1,
		PreLoadedMatrix: Generate1DLSNMatrix(params.L, params.K, params.PrimeField, seed),
	}
}

func (pi *LpnPIR) Encode(sk SecretKey, matrix pir.Matrix) *pir.Matrix {
	params := pi.Params

	encMatrix := Generate1DEncodingMatrix(params.L, params.K, params.PrimeField, sk.LinearCodeKey)

	sliceRows := params.M / params.M_1
	sliceSize := sliceRows * params.N

	// Partition the Rows into blocks with each block size M_1
	// matrix.Rows should equal to params.M
	numBlock := (matrix.Rows + params.M_1 - 1) / params.M_1

	data := make([]uint32, params.N*numBlock*uint32(params.ECCLength))

	// Reused Slot for ECC encoding
	message := make([]uint32, params.ECCLength)

	rsGeneratorMatrix := GetRSGeneratorMatrix(params.M_1, params.ECCLength, params.PrimeField)

	for i := uint32(0); i < numBlock; i++ {
		for j := uint32(0); j < params.M_1; j++ {
			// Each block has size params.L*params.M_1, Each Row has size params.L
			matrixStart := i*params.L*params.M_1 + j*params.L

			// Put into the jth slice, ith row, each row with size params.N
			dataStart := j*sliceSize + i*params.N

			// Copy the matrix row to the first L places
			copy(data[dataStart:dataStart+params.L], matrix.Data[matrixStart:matrixStart+params.L])

			// Encode S x message into the [L:L+K] places
			C.MatVecProduct(
				(*C.uint32_t)(unsafe.Pointer(&encMatrix[0])),
				(*C.uint32_t)(unsafe.Pointer(&matrix.Data[matrixStart])),
				(*C.uint32_t)(unsafe.Pointer(&data[dataStart+params.L])),
				C.uint32_t(params.K), C.uint32_t(params.L), C.uint32_t(params.PrimeField))

		}

		// Encode each M_1 length slice with ECC to length ECCLength
		for j := uint32(0); j < params.N; j++ {
			// Get the row i, col j of each block
			for t := uint32(0); t < params.M_1; t++ {
				message[t] = data[t*sliceSize+i*params.N+j]
			}

			C.MatVecProduct(
				(*C.uint32_t)(unsafe.Pointer(&rsGeneratorMatrix[0])),
				(*C.uint32_t)(unsafe.Pointer(&message[0])),
				(*C.uint32_t)(unsafe.Pointer(&message[params.M_1])),
				C.uint32_t(params.ECCLength), C.uint32_t(params.M_1), C.uint32_t(params.PrimeField))

			// Put to the M_1:ECCLength slice
			for t := params.M_1; t < params.ECCLength; t++ {
				data[t*sliceSize+i*params.N+j] = message[t]
			}
		}
	}

	return &pir.Matrix{
		Rows: numBlock * uint32(params.ECCLength),
		Cols: pi.Params.N,
		Data: data,
	}
}

func (pi *LpnPIR) Query(sk SecretKey, queryIndex uint64) (*LpnPIRQuery, *LpnPIRAux) {
	params := pi.Params

	// Sample Vector From Nullspace
	SofLSN := sk.PreLoadedMatrix
	if len(SofLSN) == 0 {
		SofLSN = Generate1DLSNMatrix(params.L, params.K, params.PrimeField, sk.LinearCodeKey)
	}

	queryVector := make([]uint32, params.N*params.ECCLength)
	noisyQuery := make([]bool, params.ECCLength)

	for t := uint32(0); t < params.ECCLength; t++ {
		// r \in F^k
		r := pir.RandomPrimeFieldVector(params.K, params.PrimeField)

		// e \in Ber(epsi)^L
		e := pir.RandomNoiseVector(params.L, params.Epsi, params.PrimeField)

		C.MatVecProduct(
			(*C.uint32_t)(unsafe.Pointer(&SofLSN[0])),
			(*C.uint32_t)(unsafe.Pointer(&r[0])),
			(*C.uint32_t)(unsafe.Pointer(&queryVector[t*params.N])),
			C.uint32_t(params.L), C.uint32_t(params.K), C.uint32_t(params.PrimeField))

		// Add Noise If non-zero and flag the prediction vector
		if !pir.IsZeroVector(e) {
			noisyQuery[t] = true
			for i := uint32(0); i < params.L; i++ {
				queryVector[t*params.N+i] = uint32((uint64(queryVector[t*params.N+i]) + uint64(e[i])) % uint64(params.PrimeField))
			}
		}

		copy(queryVector[t*params.N+params.L:t*params.N+params.N], r[:params.K])

		queryVector[uint64(t*params.N)+queryIndex%uint64(params.L)] = (queryVector[uint64(t*params.N)+queryIndex%uint64(params.L)] + 1) % params.PrimeField
	}

	return &LpnPIRQuery{
			Queries:      queryVector,
			QueryLen:     params.N,
			NumOfQueries: params.ECCLength,
		}, &LpnPIRAux{
			NoisyQueryInd: noisyQuery,
		}
}

func (pi *LpnPIR) Answer(matrix *pir.Matrix, clientQuery *LpnPIRQuery) *LpnPIRAnswer {
	params := pi.Params
	sliceRows := params.M / params.M_1
	sliceSize := sliceRows * params.N

	answers := make([]uint32, sliceRows*params.ECCLength)

	// Query Slice i with Query Vector i And Put into Answer i (Length = sliceRows)
	for i := uint32(0); i < params.ECCLength; i++ {
		C.MatVecProduct(
			(*C.uint32_t)(unsafe.Pointer(&matrix.Data[i*sliceSize])),
			(*C.uint32_t)(unsafe.Pointer(&clientQuery.Queries[i*clientQuery.QueryLen])),
			(*C.uint32_t)(unsafe.Pointer(&answers[i*sliceRows])),
			C.uint32_t(sliceRows), C.uint32_t(params.N), C.uint32_t(params.PrimeField))
	}

	return &LpnPIRAnswer{
		Answers: answers,
		AnsLen:  params.M / params.M_1,
	}
}

func (pi *LpnPIR) Decode(sk SecretKey, index uint64, response *LpnPIRAnswer, aux *LpnPIRAux) (uint32, error) {
	params := pi.Params
	// Each Block has size L x M_1, Same Block will become same row of each slice
	// The row number inside of each block becomes slice number
	// Data i falls into (Col: i % L, Row: i / (L * M_1), Slice: (i % (L * M_1))/L)
	indexRow := index / (uint64(params.L) * uint64(params.M_1))
	sliceNum := (index % (uint64(params.L) * uint64(params.M_1))) / uint64(params.L)

	// TODO: UnMask the response.Answers

	ecc := make([]uint32, params.ECCLength)
	// Extract the indexRow element from each answer
	for i := uint32(0); i < params.ECCLength; i++ {
		ecc[i] = response.Answers[i*response.AnsLen+uint32(indexRow)]
	}

	return DecodeECC(ecc, aux.NoisyQueryInd, uint32(sliceNum), params.M_1)
}
