#include "NTT.h"
#include <iostream>
#include <cstdint>
#include <cstdlib>  // for rand()
#include <cassert>



u32 mod_pow(u32 base, u32 exp, u32 mod) {
    u64 result = 1;
    u64 b = base;
    while (exp > 0) {
        if (exp & 1) result = (result * b) % mod;
        b = (b * b) % mod;
        exp >>= 1;
    }
    return (u32)result;
}

u32 mod_inv(u32 a, u32 mod) {
    return mod_pow(a, mod - 2, mod); // Fermat's little theorem, assuming mod is prime
}

// Modular exponentiation: (base^exp) % mod
uint32_t modExponent(uint32_t base, uint32_t exp, uint32_t mod) {
    uint64_t result = 1;
    uint64_t b = base % mod;
    while (exp > 0) {
        if (exp & 1)
            result = (result * b) % mod;
        b = (b * b) % mod;
        exp >>= 1;
    }
    return (uint32_t)result;
}

// Check if beta^k == 1 mod M for some 2 ≤ k < N
bool existSmallN(uint32_t beta, uint32_t M, uint32_t N) {
    for (uint32_t k = 2; k < N; ++k) {
        if (modExponent(beta, k, M) == 1)
            return true;
    }
    return false;
}

// Return a primitive N-th root of unity modulo M
uint32_t NthRootOfUnity(uint32_t M, uint32_t N) {
    assert(M > 1);
    assert((M - 1) % N == 0);  // Ensure N divides M-1
    uint32_t phi = M - 1;

    while (true) {
        uint32_t alpha = 1 + (rand() % (M - 1)); // pick alpha ∈ [1, M-1]
        uint32_t beta = modExponent(alpha, phi / N, M);
        if (!existSmallN(beta, M, N)) {
            return beta;
        }
    }
}

void print_array(const char* name, const u32* arr, size_t n) {
    std::cout << name << ": [ ";
    for (size_t i = 0; i < n; ++i) {
        std::cout << arr[i];
        if (i + 1 < n) std::cout << ", ";
    }
    std::cout << " ]" << std::endl;
}


void ntt(u32* a, size_t n, u32 root, u32 mod) {
    // Bit-reversal permutation
    size_t j = 0;
    for (size_t i = 1; i < n; ++i) {
        size_t bit = n >> 1;
        while (j & bit) {
            j ^= bit;
            bit >>= 1;
        }
        j ^= bit;
        if (i < j) {
            std::swap(a[i], a[j]);
        }
    }

    for (size_t len = 2; len <= n; len <<= 1) {
        u32 wlen = mod_pow(root, n / len, mod);

        for (size_t i = 0; i < n; i += len) {
            u32 w = 1;
            for (size_t j = 0; j < len / 2; ++j) {
                u32 u = a[i + j];
                u32 v = (u64)a[i + j + len / 2] * w % mod;
                a[i + j] = (u + v) % mod;
                a[i + j + len / 2] = (mod + u - v) % mod;
                w = (u64)w * wlen % mod;
            }
        }
    }
}

void intt(u32* a, size_t n, u32 root, u32 mod) {
    u32 inv_root = mod_inv(root, mod);

    ntt(a, n, inv_root, mod);

    u32 inv_n = mod_inv(n, mod);
    for (size_t i = 0; i < n; ++i) {
        a[i] = (u64)a[i] * inv_n % mod;
    }
}

// Convolution using NTT
void ntt_convolution(const u32* a, const u32* b, u32* result, size_t n, u32 root, u32 mod) {
    u32* fa = new u32[n];
    u32* fb = new u32[n];
    
    for (size_t i = 0; i < n; ++i) {
        fa[i] = a[i];
        fb[i] = b[i];
    }

    ntt(fa, n, root, mod);
    ntt(fb, n, root, mod);

    for (size_t i = 0; i < n; ++i) {
        fa[i] = (u64)fa[i] * fb[i] % mod;
    }

    intt(fa, n, root, mod);

    for (size_t i = 0; i < n; ++i) {
        result[i] = fa[i];
    }

    delete[] fa;
    delete[] fb;
}