package mvp

import (
	"RandomLinearCodePIR/dataobjects"
	"RandomLinearCodePIR/linearcode"
	"RandomLinearCodePIR/tdm"
	"time"
)

type SlsnMVP struct {
	Params SlsnParams
}

type SecretKey struct {
	LinearCodeKey   int64
	TDMKey          int64
	PreLoadedMatrix []uint32
	TDM             *tdm.TDM
}

// N = K + L denotes the length of the codeword
// Encoding Matrix D with dimension N x L
// Original Data Matrix has dimension M x N
// S denotes the number of blocks
// B denotes the block size
// We assume N = S x B
type SlsnParams struct {
	Field dataobjects.Field
	// Temporarily add P here
	P uint32
	S uint32
	B uint32
	K uint32
	L uint32
	N uint32
	M uint32
}

type SlsnQuery struct {
	Vec []uint32
}

type SlsnAux struct {
	Coeff []uint32
	Masks []uint32
	Dur   time.Duration
}

func (slsn *SlsnMVP) KeyGen(seed int64) SecretKey {
	params := slsn.Params
	return SecretKey{
		LinearCodeKey:   seed,
		PreLoadedMatrix: linearcode.Generate1DDualMatrix(params.L, params.K, params.Field, seed),
		TDM: &tdm.TDM{
			M: params.M,
			N: params.N,
			// NOTE: Now TDM only support Q = 2^x + 1, Change this to Field later
			Q:     params.P,
			SeedL: seed + 1,
			SeedP: seed + 500,
			SeedR: seed + 1000,
		},
	}
}

func (slsn *SlsnMVP) GenerateTDM(sk SecretKey) []uint32 {
	return sk.TDM.GenerateFlattenedTrapDooredMatrix()
}

func (slsn *SlsnMVP) Encode(sk SecretKey, input dataobjects.Matrix, mask []uint32) *dataobjects.Matrix {
	params := slsn.Params
	rlcMatrix := linearcode.Generate1DRLCMatrix(params.L, params.K, params.Field, sk.LinearCodeKey)
	encoded := make([]uint32, input.Rows*params.N)

	for i := uint32(0); i < input.Rows; i++ {
		copy(encoded[i*params.N:i*params.N+params.L], input.Data[i*params.L:(i+1)*params.L])

		MatVecProduct(rlcMatrix, input.Data[i*input.Cols:(i+1)*input.Cols], encoded[i*params.N+params.L:(i+1)*params.N],
			params.K, params.L, params.P)
	}

	// Add Masks
	for j := uint64(0); j < uint64(len(encoded)); j++ {
		encoded[j] = params.Field.Add(encoded[j], mask[j])
	}

	return &dataobjects.Matrix{
		Rows: params.M,
		Cols: params.N,
		Data: encoded,
	}
}

func (slsn *SlsnMVP) Query(sk SecretKey, vec []uint32) (*SlsnQuery, *SlsnAux) {
	params := slsn.Params

	PofDual := sk.PreLoadedMatrix
	if len(PofDual) == 0 {
		PofDual = linearcode.Generate1DDualMatrix(params.L, params.K, params.Field, sk.LinearCodeKey)
	}

	// Sample codeword c From NullSpace
	nullspaceCoeff := params.Field.SampleVector(params.K)

	queryVector := make([]uint32, params.N)

	MatVecProduct(PofDual, nullspaceCoeff, queryVector, params.L, params.K, params.P)

	copy(queryVector[params.L:params.N], nullspaceCoeff[:params.K])

	// Add Vector v to c
	for i := uint32(0); i < params.L; i++ {
		queryVector[i] = params.Field.Add(queryVector[i], vec[i])
	}

	// The time is just for benchmark
	start := time.Now()
	// Calculate The Mask
	masks := sk.TDM.EvaluationCircuit(queryVector)
	dur := time.Since(start)

	// Generate Non-zero coefficient
	coeff := params.Field.SampleInvertibleVec(params.S)

	for i := uint32(0); i < params.S; i++ {
		for j := uint32(0); j < params.B; j++ {
			queryVector[i*params.B+j] = params.Field.Mul(queryVector[i*params.B+j], coeff[i])
		}
	}

	return &SlsnQuery{
			Vec: queryVector,
		}, &SlsnAux{
			Coeff: coeff,
			Masks: masks,
			Dur:   dur,
		}
}

func (slsn *SlsnMVP) Answer(encodedMatrix dataobjects.Matrix, clientQuery SlsnQuery) []uint32 {
	params := slsn.Params
	result := make([]uint32, params.S*params.M)

	BlockMatVecProduct(encodedMatrix.Data, clientQuery.Vec, result, params.M, params.N, params.S, params.P)
	return result
}

func (slsn *SlsnMVP) Decode(sk SecretKey, response []uint32, aux SlsnAux) []uint32 {
	params := slsn.Params

	vec := params.Field.InvertVector(aux.Coeff)

	result := make([]uint32, params.M)

	BlockVecMatProduct(response, vec, result, params.S, params.M, 1, params.P)
	// Unmask
	for i := uint32(0); i < params.M; i++ {
		result[i] = params.Field.Sub(result[i], aux.Masks[i])
	}

	return result
}
