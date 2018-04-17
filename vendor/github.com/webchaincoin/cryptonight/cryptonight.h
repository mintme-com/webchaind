#ifndef CRYPTONIGHT_H
#define CRYPTONIGHT_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

void cryptonight_hash(void *ctx, const char* input, char* output, uint32_t len);

void *cryptonight_create(void);
void cryptonight_destroy(void *ctx);

#ifdef __cplusplus
}
#endif

#endif
