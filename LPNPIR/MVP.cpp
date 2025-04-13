#include <iostream>

#include <NTL/ZZ_p.h>
#include <NTL/ZZ_pX.h>
#include <cstdint>

inline uint32_t mod_add(uint32_t a, uint32_t b, uint32_t p) {
    uint64_t sum = (uint64_t)a + b;
    return (sum >= p) ? sum - p : (uint32_t)sum;
}

inline uint32_t mod_mul(uint32_t a, uint32_t b, uint32_t p) {
    return (uint32_t)(((uint64_t)a * b) % p);
}

extern "C" {

// M x v
void MatVecProduct(const uint32_t* mat, const uint32_t* vec, uint32_t* result, uint32_t n, uint32_t m, uint32_t p)
{
    for (int row = 0; row < n; ++row) {
        const uint32_t* row_ptr = mat + row * m;

        uint32_t acc = 0;
        for (int col = 0; col < m; ++col) {
            acc = mod_add(acc, mod_mul(row_ptr[col], vec[col], p), p);
        }

        result[row] = acc;
    }
}

}