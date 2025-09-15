package mvp

/*
#cgo CXXFLAGS: -std=c++17 -O3 -march=native  -I.
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lMVP -L/opt/homebrew/lib -lNTT -lntl -lgmp -lstdc++
#include "mvp.h"
#include "matrixShapeTransform.h"
*/
import "C"
import "unsafe"

func BlockMatVecProduct(mat, vec, out []uint32, row, col, numBlock, p uint32) {
	C.BlockMatVecProduct(
		(*C.uint32_t)(unsafe.Pointer(&mat[0])),
		(*C.uint32_t)(unsafe.Pointer(&vec[0])),
		(*C.uint32_t)(unsafe.Pointer(&out[0])),
		C.uint32_t(row), C.uint32_t(col), C.uint32_t(numBlock), C.uint32_t(p),
	)
}

func MatVecProduct(mat, vec, out []uint32, row, col, p uint32) {
	C.MatVecProduct(
		(*C.uint32_t)(unsafe.Pointer(&mat[0])),
		(*C.uint32_t)(unsafe.Pointer(&vec[0])),
		(*C.uint32_t)(unsafe.Pointer(&out[0])),
		C.uint32_t(row), C.uint32_t(col), C.uint32_t(p),
	)
}

func BlockVecMatProduct(mat, vec, out []uint32, row, col, numBlock, p uint32) {
	C.BlockVecMatProduct(
		(*C.uint32_t)(unsafe.Pointer(&mat[0])),
		(*C.uint32_t)(unsafe.Pointer(&vec[0])),
		(*C.uint32_t)(unsafe.Pointer(&out[0])),
		C.uint32_t(row), C.uint32_t(col), C.uint32_t(numBlock), C.uint32_t(p),
	)
}

func TransformToBlockwise(mat, matBlocked []uint32, n, m, s uint32) {
	C.TransformRowMajorToBlockRowMajor(
		(*C.uint32_t)(unsafe.Pointer(&mat[0])),
		(*C.uint32_t)(unsafe.Pointer(&matBlocked[0])),
		C.uint32_t(n), C.uint32_t(m), C.uint32_t(s),
	)
}
