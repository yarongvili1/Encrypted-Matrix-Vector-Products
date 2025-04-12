#ifndef PRIME_FIELD_UTILS_H
#define PRIME_FIELD_UTILS_H

// #include <vector>
// #include <random>
#include <stdint.h>
// #include <chrono>

#ifdef __cplusplus
extern "C" {
#endif

// C-callable wrappers
void* GenerateRandomColsOfC(uint32_t K, uint32_t L, uint32_t p, int64_t seed);
void* GenerateRandomColsOfD(uint32_t K, uint32_t L, uint32_t p, int64_t seed);
void* SampleVectorFromNullSpace(uint32_t* out, uint32_t K, uint32_t L, uint32_t p, int64_t seed);

#ifdef __cplusplus
}
#endif

#endif // PRIME_FIELD_UTILS_H