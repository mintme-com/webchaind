#ifndef CRYPTONIGHT_H
#define CRYPTONIGHT_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include "hash-ops.h"

#define MEMORY         (1 << 21) /* 2 MiB */
#define ITER           (1 << 20)
#define AES_BLOCK_SIZE  16
#define AES_KEY_SIZE    32 /*16*/
#define INIT_SIZE_BLK   8
#define INIT_SIZE_BYTE (INIT_SIZE_BLK * AES_BLOCK_SIZE)

#define VARIANT_WEB_1_1(p) \
    const uint8_t tmp1 = ((const uint8_t*)(p))[11]; \
    static const uint32_t table1 = 0x86420; \
    const uint8_t index1 = (((tmp1 >> 3) & 6) | (tmp1 & 1)) << 1; \
    ((uint8_t*)(p))[11] = tmp1 ^ ((table1 >> index1) & 0x30);

#define VARIANT_WEB_1_2(p) \
    const uint8_t tmp2 = ((const uint8_t*)(p))[1]; \
    static const uint32_t table2 = 0x75310; \
    const uint8_t index2 = (((tmp2 >> 3) & 6) | (tmp2 & 1)) << 1; \
    ((uint8_t*)(p))[1] = tmp2 ^ ((table2 >> index2) & 0x33);

#pragma pack(push, 1)
union cn_slow_hash_state {
    union hash_state hs;
    struct {
        uint8_t k[64];
        uint8_t init[INIT_SIZE_BYTE];
    };
};
#pragma pack(pop)

extern int aesb_single_round(const uint8_t *in, uint8_t*out, const uint8_t *expandedKey);
extern int aesb_pseudo_round(const uint8_t *in, uint8_t *out, const uint8_t *expandedKey);

static inline size_t e2i(const uint8_t* a) {
    return (*((uint64_t*) a) / AES_BLOCK_SIZE) & (MEMORY / AES_BLOCK_SIZE - 1);
}

static inline void copy_block(uint8_t* dst, const uint8_t* src) {
    ((uint64_t*) dst)[0] = ((uint64_t*) src)[0];
    ((uint64_t*) dst)[1] = ((uint64_t*) src)[1];
}

static inline void xor_blocks(uint8_t *a, const uint8_t *b) {
    ((uint64_t *) a)[0] ^= ((uint64_t *) b)[0];
    ((uint64_t *) a)[1] ^= ((uint64_t *) b)[1];
}

static inline void xor_blocks_dst(const uint8_t *a, const uint8_t *b, uint8_t *dst) {
    ((uint64_t *) dst)[0] = ((uint64_t *) a)[0] ^ ((uint64_t *) b)[0];
    ((uint64_t *) dst)[1] = ((uint64_t *) a)[1] ^ ((uint64_t *) b)[1];
}

void cryptonight_hash(void *ctx, const char* input, char* output, uint32_t len);
void cryptonight_hash_aesni(void *ctx, const char* input, char* output, uint32_t len);

void *cryptonight_create(void);
void cryptonight_destroy(void *ctx);

extern void (* const extra_hashes[4])(const void *, size_t, char *);

#ifdef __cplusplus
}
#endif

#endif
