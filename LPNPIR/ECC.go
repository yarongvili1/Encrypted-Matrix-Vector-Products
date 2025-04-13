package lpnpir

/*
#include "ReedSolomon.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

// Only return the evaluation part
func GetRSGeneratorMatrix(M_1, ECCLength, p uint32) []uint32 {
	alphas := getAlphas(ECCLength)
	rsGeneratorMatrix := make([]uint32, M_1*ECCLength)
	C.GenerateSystematicRSMatrix_uint32(
		C.uint32_t(ECCLength), C.uint32_t(M_1), C.uint32_t(p),
		(*C.uint32_t)(unsafe.Pointer(&alphas[0])),
		(*C.uint32_t)(unsafe.Pointer(&rsGeneratorMatrix[0])))

	return rsGeneratorMatrix[M_1*M_1:]

}

func DecodeECC(ecc []uint32, noisyQuery []bool, index, M_1 uint32) (uint32, error) {
	if !noisyQuery[index] {
		return ecc[index], nil
	} else {
		x_in := make([]uint32, M_1)
		y_in := make([]uint32, M_1)
		idx := uint32(0)
		for i := range noisyQuery {
			if !noisyQuery[i] && idx < M_1 {
				x_in[idx] = uint32(i)
				y_in[idx] = ecc[i]
				idx += 1
			}
		}

		if idx < M_1 {
			return 0, errors.New("Decoding Failed Due To Not Enough Data.")
		}

		t := C.LagrangeInterpEval(
			(*C.uint32_t)(unsafe.Pointer(&x_in[0])),
			(*C.uint32_t)(unsafe.Pointer(&y_in[0])),
			C.uint32_t(M_1), C.uint32_t(index), C.uint32_t(7))

		return uint32(t), nil
	}
}

func getAlphas(ECCLength uint32) []uint32 {
	alphas := make([]uint32, ECCLength)
	for i := range alphas {
		alphas[i] = uint32(i)
	}
	return alphas
}
