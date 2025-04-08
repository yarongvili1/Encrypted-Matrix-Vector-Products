#ifndef MATRIX_VECTOR_MULTIPLY_H
#define MATRIX_VECTOR_MULTIPLY_H

#include <stdint.h>

void MatrixMulVector(uint32_t* result, uint32_t* matrix, uint32_t* vector, uint32_t rows, uint32_t cols);

void MatrixColXOR(uint32_t* result, uint32_t* matrix, uint32_t* vector, uint32_t rows, uint32_t cols);

void MatrixColXORByBlock(uint32_t* result_1, uint32_t* result_2, uint32_t* matrix, uint32_t* vector_1, uint32_t* vector_2, uint32_t rows, uint32_t cols, uint32_t block_size);

#endif
