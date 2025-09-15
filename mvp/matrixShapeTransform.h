// transform.h
#ifndef TRANSFORM_H
#define TRANSFORM_H

#ifdef __cplusplus
extern "C" {
#endif

void TransformRowMajorToBlockRowMajor(
    const uint32_t* mat,
    uint32_t* matBlocked,
    uint32_t n, uint32_t m, uint32_t s
);

#ifdef __cplusplus
}
#endif

#endif  // TRANSFORM_H
