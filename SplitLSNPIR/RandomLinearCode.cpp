#include <vector>
#include <random>
#include <cstdint>
#include <chrono>

extern "C" {

void GenerateRandomColsOfC(uint32_t* out, uint32_t K, uint32_t L, uint32_t p, int64_t seed) {
    std::mt19937 rng(seed);
    std::uniform_int_distribution<uint32_t> dist(0, p - 1);

    for (uint32_t i = 0; i < K; ++i) {
        for (uint32_t j = 0; j < L; ++j) {
            out[i * L + j] = dist(rng);
        }
    }
}

void GenerateRandomColsOfD(uint32_t* out, uint32_t K, uint32_t L, uint32_t p, int64_t seed) {
    std::vector<uint32_t> P(K * L);
    GenerateRandomColsOfC(P.data(), K, L, p, seed);

    // Fill top L rows with -P^T mod p
    for (uint32_t i = 0; i < L; ++i) {
        for (uint32_t j = 0; j < K; ++j) {
            out[i * (K + L) + j] = (p - P[j * L + i]) % p;
        }
    }

    // Fill identity block (bottom L x L) into rows
    for (uint32_t i = 0; i < L; ++i) {
        for (uint32_t j = 0; j < L; ++j) {
            out[i * (K + L) + K + j] = (i == j) ? 1 : 0;
        }
    }
}

void SampleVectorFromNullSpace(uint32_t* out, uint32_t K, uint32_t L, uint32_t p, int64_t seed) {
    std::mt19937 rng(seed + 999);
    std::uniform_int_distribution<uint32_t> dist(1, p - 1); // non-zero coeffs

    std::vector<uint32_t> coeff(L);
    for (uint32_t i = 0; i < L; ++i) {
        coeff[i] = dist(rng);
    }

    std::vector<uint32_t> D(L * (K + L));
    GenerateRandomColsOfD(D.data(), K, L, p, seed);

    for (uint32_t i = 0; i < L; ++i) {
        uint32_t c = coeff[i];
        for (uint32_t j = 0; j < K + L; ++j) {
            uint64_t prod = static_cast<uint64_t>(D[i * (K + L) + j]) * c;
            out[j] = static_cast<uint32_t>((static_cast<uint64_t>(out[j]) + prod) % p);
        }
    }
}

} // extern "C"
