#include <iostream>

inline uint32_t mod_add(uint32_t a, uint32_t b, uint32_t p) {
    uint64_t sum = (uint64_t)a + b;
    return (sum >= p) ? sum - p : (uint32_t)sum;
}

inline uint32_t mod_mul(uint32_t a, uint32_t b, uint32_t p) {
    return (uint32_t)(((uint64_t)a * b) % p);
}

extern "C" {

void BlockMatrixVectorProduct(const uint32_t* mat, const uint32_t* vec, uint32_t* result, uint32_t n, uint32_t m, uint32_t s, uint32_t p)
{
    int b = n / s;
    for (int i = 0; i < s * m; ++i) {
        result[i] = 0;
    }

    for (int blk = 0; blk < s; ++blk) {
        int row_start = blk * b;
        for (int row = row_start; row < row_start + b; ++row) {
            uint32_t v = vec[row];
            const uint32_t* row_ptr = mat + row * m;
            uint32_t* res_ptr = result + blk * m;

            for (int col = 0; col < m; ++col) {
                uint32_t prod = mod_mul(row_ptr[col], v, p);
                res_ptr[col] = mod_add(res_ptr[col], prod, p);
            }
        }
    }
}

void BlockMatrixVectorProductByColumn(const uint32_t* mat, const uint32_t* vec, uint32_t* result, uint32_t n, uint32_t m, uint32_t s, uint32_t p)
{

    int b = m / s;
    for (int i = 0; i < s * n; ++i) {
        result[i] = 0;
    }

    for (int blk = 0; blk < s; ++blk) {
        int col_start = blk * b;
        uint32_t* res_ptr = result + blk * n; // result has shape (s x n), stored row-major

        for (int row = 0; row < n; ++row) {
            const uint32_t* row_ptr = mat + row * m;

            uint32_t acc = 0;
            for (int j = 0; j < b; ++j) {
                int col = col_start + j;
                acc = mod_add(acc, mod_mul(row_ptr[col], vec[col], p), p);
            }
            res_ptr[row] = acc;
        }
    }
}

}
