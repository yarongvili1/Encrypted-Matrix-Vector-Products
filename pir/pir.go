package pir

type PIR interface {
	KeyGen(N, Ell, Lambda int, seed int64) SecretKey
	Encode(sk SecretKey, db Matrix) *Matrix
	Query(sk SecretKey, index uint64) (*ClientQuery, *AuxiliaryInfo)
	Answer(encodedDB *Matrix, query *ClientQuery) *ServerResponse
	Decode(sk SecretKey, index uint64, ans *ServerResponse, aux *AuxiliaryInfo) uint32
}

type SecretKey struct {
	LinearCodeKey int64
	MaskKey       int64
	Lambda        int
	N             int
	Ell           int
}

type ClientQuery interface {
}

type AuxiliaryInfo interface {
}

type ServerResponse interface {
}
