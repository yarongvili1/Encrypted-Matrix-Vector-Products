// matrix_vector_multiply.c
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include "BitMVP.h"

// Function to perform matrix-vector multiplication over F2
void MatrixMulVector(uint32_t* result, uint32_t* matrix, uint32_t* vector, uint32_t rows, uint32_t cols) {
    uint64_t index = 0;
    for (int i = 0; i < rows; i++) {
        uint32_t sum = 0;
        for (int j = 0; j < cols; j++) {
            sum ^= matrix[index++] * vector[j];
        }
        result[i] = sum;
    }
}

// Function to perform Vector^T * Matrix while matrix is flattened by rows
void VecMatrixMulF2(uint32_t* result, uint32_t* matrix, uint32_t* vector, uint32_t rows, uint32_t cols) {
    for (uint32_t i = 0; i < rows; i++) {
        if (vector[i] == 1) {
            uint32_t start = i * cols;
            for (uint32_t j = 0; j < cols; j++) {
                result[j] ^= matrix[start + j];
            }
        }
    }
}

// Function to perform Vector * Matrix(Flattened) in Blocks over F2, and vectors are packed by 32 entries.
void MatrixColXORByBlock2D(uint32_t* result_1, uint32_t* result_2, uint32_t* matrix, uint32_t* vector_1, uint32_t* vector_2, uint32_t rows, uint32_t cols, uint32_t block_size) {
    uint32_t nblocks = (rows + block_size - 1) / block_size;

    for (uint32_t k = 0; k < nblocks; k++) {
        uint32_t block_start = k * block_size;
        uint32_t block_end   = block_start + block_size;
        if (block_end > rows) block_end = rows;

        // base pointers for results
        uint32_t* R1 = result_1 + (size_t)k * cols;
        uint32_t* R2 = result_2 + (size_t)k * cols;

        // find first/last word indices
        uint32_t w_start = block_start >> 5;
        uint32_t w_end   = (block_end - 1) >> 5;

        for (uint32_t w = w_start; w <= w_end; w++) {
            uint32_t base_row = w << 5;
            uint32_t b0 = (block_start > base_row) ? (block_start - base_row) : 0;
            uint32_t b1 = (block_end   - base_row > 32) ? 32 : (block_end - base_row);

            uint32_t word1 = vector_1[w];
            uint32_t word2 = vector_2[w];

            for (uint32_t b = b0; b < b1; b++) {
                uint32_t row = base_row + b;

                if ((word1 >> b) & 1u) {
                    const uint32_t* rowp = matrix + (size_t)row * cols;
                    for (uint32_t j = 0; j < cols; j++)
                        R1[j] ^= rowp[j];
                }
                if ((word2 >> b) & 1u) {
                    const uint32_t* rowp = matrix + (size_t)row * cols;
                    for (uint32_t j = 0; j < cols; j++)
                        R2[j] ^= rowp[j];
                }
            }
        }
    }
}

void MatrixColXORByBlock(uint32_t* result, uint32_t* matrix, uint32_t* vector, uint32_t rows, uint32_t cols, uint32_t block_size) {
    for (uint32_t k = 0; k < rows/block_size; k++) {
        for (uint32_t i = k * block_size; i < k * block_size + block_size; i++) {
            if (vector[i] == 1) {
                uint32_t start = i * cols;
                for (uint32_t j = 0; j < cols; j++) {
                    result[k * cols + j] ^= matrix[start + j];
                }
            }
        }
    }
}