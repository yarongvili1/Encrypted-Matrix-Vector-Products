// matrix_vector_multiply.c
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>

#include "answer.h"

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

// Function to perform XOR of Columns
void MatrixColXOR(uint32_t* result, uint32_t* matrix, uint32_t* vector, uint32_t rows, uint32_t cols) {
    for (uint32_t i = 0; i < rows; i++) {
        if (vector[i] == 1) {
            uint32_t start = i * cols;
            for (uint32_t j = 0; j < cols; j++) {
                result[j] ^= matrix[start + j];
            }
        }
    }
}

// Function to perform XOR of Columns in Blocks
void MatrixColXORByBlock(uint32_t* result_1, uint32_t* result_2, uint32_t* matrix, uint32_t* vector_1, uint32_t* vector_2, uint32_t rows, uint32_t cols, uint32_t block_size) {
    for (uint32_t k = 0; k < rows/block_size; k++) {
        for (uint32_t i = k * block_size; i < k * block_size + block_size; i++) {
            if (vector_1[i] == 1) {
                uint32_t start = i * cols;
                for (uint32_t j = 0; j < cols; j++) {
                    result_1[k * cols + j] ^= matrix[start + j];
                }
            }
            if (vector_2[i] == 1) {
                uint32_t start = i * cols;
                for (uint32_t j = 0; j < cols; j++) {
                    result_2[k * cols + j] ^= matrix[start + j];
                }
            }
        }
    }

}
