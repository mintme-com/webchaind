// Copyright (c) 2012-2013 The Cryptonote developers
// Copyright (c) 2018 Webchain project
// Distributed under the MIT/X11 software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

#include "oaes_lib.h"
#include "c_keccak.h"
#include "c_groestl.h"
#include "c_blake256.h"
#include "c_jh.h"
#include "c_skein.h"
#include "int-util.h"
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

static void do_blake_hash(const void* input, size_t len, char* output) {
    blake256_hash((uint8_t*)output, input, len);
}

void do_groestl_hash(const void* input, size_t len, char* output) {
    groestl(input, len * 8, (uint8_t*)output);
}

static void do_jh_hash(const void* input, size_t len, char* output) {
    int r = jh_hash(HASH_SIZE * 8, input, 8 * len, (uint8_t*)output);
    assert(SUCCESS == r);
}

static void do_skein_hash(const void* input, size_t len, char* output) {
    int r = c_skein_hash(8 * HASH_SIZE, input, 8 * len, (uint8_t*)output);
    assert(SKEIN_SUCCESS == r);
}

static void (* const extra_hashes[4])(const void *, size_t, char *) = {
    do_blake_hash, do_groestl_hash, do_jh_hash, do_skein_hash
};

extern int aesb_single_round(const uint8_t *in, uint8_t*out, const uint8_t *expandedKey);
extern int aesb_pseudo_round(const uint8_t *in, uint8_t *out, const uint8_t *expandedKey);

static inline size_t e2i(const uint8_t* a) {
    return (*((uint64_t*) a) / AES_BLOCK_SIZE) & (MEMORY / AES_BLOCK_SIZE - 1);
}

static void mul_sum_xor_dst(const uint8_t* a, uint8_t* c, uint8_t* dst) {
    uint64_t hi, lo = mul128(((uint64_t*) a)[0], ((uint64_t*) dst)[0], &hi) + ((uint64_t*) c)[1];
    hi += ((uint64_t*) c)[0];

    ((uint64_t*) c)[0] = ((uint64_t*) dst)[0] ^ hi;
    ((uint64_t*) c)[1] = ((uint64_t*) dst)[1] ^ lo;
    ((uint64_t*) dst)[0] = hi;
    ((uint64_t*) dst)[1] = lo;
}

static inline void copy_block(uint8_t* dst, const uint8_t* src) {
    ((uint64_t*) dst)[0] = ((uint64_t*) src)[0];
    ((uint64_t*) dst)[1] = ((uint64_t*) src)[1];
}

static inline void xor_blocks(uint8_t* a, const uint8_t* b) {
    ((uint64_t*) a)[0] ^= ((uint64_t*) b)[0];
    ((uint64_t*) a)[1] ^= ((uint64_t*) b)[1];
}

static inline void xor_blocks_dst(const uint8_t* a, const uint8_t* b, uint8_t* dst) {
    ((uint64_t*) dst)[0] = ((uint64_t*) a)[0] ^ ((uint64_t*) b)[0];
    ((uint64_t*) dst)[1] = ((uint64_t*) a)[1] ^ ((uint64_t*) b)[1];
}

struct cryptonight_ctx {
    uint8_t long_state[MEMORY];
    union cn_slow_hash_state state;
    uint8_t text[INIT_SIZE_BYTE];
    uint8_t a[AES_BLOCK_SIZE];
    uint8_t b[AES_BLOCK_SIZE];
    uint8_t c[AES_BLOCK_SIZE];
    uint8_t aes_key[AES_KEY_SIZE];
    oaes_ctx* aes_ctx;
};

void cryptonight_hash(void *ctx2, const char* input, char* output, uint32_t len) {
    //struct cryptonight_ctx *ctx = alloca(sizeof(struct cryptonight_ctx));
	struct cryptonight_ctx *ctx = ctx2;
    hash_process(&ctx->state.hs, (const uint8_t*) input, len);
    memcpy(ctx->text, ctx->state.init, INIT_SIZE_BYTE);
    memcpy(ctx->aes_key, ctx->state.hs.b, AES_KEY_SIZE);
    ctx->aes_ctx = (oaes_ctx*) oaes_alloc();
    size_t i, j;

    oaes_key_import_data(ctx->aes_ctx, ctx->aes_key, AES_KEY_SIZE);
    for (i = 0; i < MEMORY / INIT_SIZE_BYTE; i++) {
        for (j = 0; j < INIT_SIZE_BLK; j++) {
            aesb_pseudo_round(&ctx->text[AES_BLOCK_SIZE * j],
                    &ctx->text[AES_BLOCK_SIZE * j],
                    ctx->aes_ctx->key->exp_data);
        }
        memcpy(&ctx->long_state[i * INIT_SIZE_BYTE], ctx->text, INIT_SIZE_BYTE);
    }

    for (i = 0; i < 16; i++) {
        ctx->a[i] = ctx->state.k[i] ^ ctx->state.k[32 + i];
        ctx->b[i] = ctx->state.k[16 + i] ^ ctx->state.k[48 + i];
    }

    for (i = 0; i < ITER / 2; i++) {
        /* Dependency chain: address -> read value ------+
         * written value <-+ hard function (AES or MUL) <+
         * next address  <-+
         */
        /* Iteration 1 */
        j = e2i(ctx->a);
        aesb_single_round(&ctx->long_state[j * AES_BLOCK_SIZE], ctx->c, ctx->a);
        xor_blocks_dst(ctx->c, ctx->b, &ctx->long_state[j * AES_BLOCK_SIZE]);
        VARIANT_WEB_1_1((uint8_t*)&ctx->long_state[j * AES_BLOCK_SIZE]);
        /* Iteration 2 */
        mul_sum_xor_dst(ctx->c, ctx->a,
                &ctx->long_state[e2i(ctx->c) * AES_BLOCK_SIZE]);
        copy_block(ctx->b, ctx->c);
        VARIANT_WEB_1_2((uint8_t*)
                &ctx->long_state[e2i(ctx->c) * AES_BLOCK_SIZE]);
    }

    memcpy(ctx->text, ctx->state.init, INIT_SIZE_BYTE);
    oaes_key_import_data(ctx->aes_ctx, &ctx->state.hs.b[32], AES_KEY_SIZE);
    for (i = 0; i < MEMORY / INIT_SIZE_BYTE; i++) {
        for (j = 0; j < INIT_SIZE_BLK; j++) {
            xor_blocks(&ctx->text[j * AES_BLOCK_SIZE],
                    &ctx->long_state[i * INIT_SIZE_BYTE + j * AES_BLOCK_SIZE]);
            aesb_pseudo_round(&ctx->text[j * AES_BLOCK_SIZE],
                    &ctx->text[j * AES_BLOCK_SIZE],
                    ctx->aes_ctx->key->exp_data);
        }
    }
    memcpy(ctx->state.init, ctx->text, INIT_SIZE_BYTE);
    hash_permutation(&ctx->state.hs);
    /*memcpy(hash, &state, 32);*/
    extra_hashes[ctx->state.hs.b[0] & 3](&ctx->state, 200, output);
    oaes_free((OAES_CTX **) &ctx->aes_ctx);
}

void *cryptonight_create(void)
{
	return malloc(sizeof(struct cryptonight_ctx));
}

void cryptonight_destroy(void *ctx)
{
	if (ctx) free(ctx);
}
