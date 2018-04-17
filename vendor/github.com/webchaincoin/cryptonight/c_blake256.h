#ifndef _BLAKE256_H_
#define _BLAKE256_H_

#include <stdint.h>

typedef struct {
  uint32_t h[8], s[4], t[2];
  int buflen, nullt;
  uint8_t buf[64];
} state;


void blake256_init(state *);

void blake256_update(state *, const uint8_t *, uint64_t);

void blake256_final(state *, uint8_t *);

void blake256_hash(uint8_t *, const uint8_t *, uint64_t);

#endif /* _BLAKE256_H_ */
