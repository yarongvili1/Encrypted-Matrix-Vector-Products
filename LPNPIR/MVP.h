#ifndef WRAPPER_H
#define WRAPPER_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

void MatVecProduct(const uint32_t* mat, const uint32_t* vec, uint32_t* result, uint32_t n, uint32_t m, uint32_t p);

#ifdef __cplusplus
}
#endif

#endif