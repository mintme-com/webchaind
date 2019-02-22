// Copyright 2015 The go-ethereum Authors
// Copyright 2015 Lefteris Karapetsas <lefteris@refu.co>
// Copyright 2015 Matthew Wampler-Doty <matthew.wampler.doty@gmail.com>
// Copyright 2018-2019 Webchain project
// This file is part of Webchain.
//
// Webchain is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Webchain is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Webchain. If not, see <http://www.gnu.org/licenses/>.

package cryptonight

// TODO: Refactor this package. "cryptonight" isn't right place for
// LYRA2 implementation.

/*
#cgo amd64 CFLAGS: -maes
#include "cryptonight.h"
#include "Lyra2.h"
#include <stdlib.h>
*/
import "C"

import (
	"encoding/binary"
	"math/big"
	"math/rand"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/webchain-network/webchaind/common"
	"github.com/webchain-network/webchaind/core/types"
	"github.com/webchain-network/webchaind/logger"
	"github.com/webchain-network/webchaind/logger/glog"
	"github.com/webchain-network/webchaind/pow"

	"github.com/klauspost/cpuid"
)

var (
	maxUint256  = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))
)

type Cryptonight struct {
	lyra2_height uint64
	hashRate     int32
	turbo        bool
}

func (pow *Cryptonight) GetHashrate() int64 {
	return int64(atomic.LoadInt32(&pow.hashRate))
}

func (pow *Cryptonight) Search(block pow.Block, stop <-chan struct{}, index int) (nonce uint64) {
	is_lyra2 := block.NumberU64() >= pow.lyra2_height
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	diff := block.Difficulty()

	i := int64(0)
	start := time.Now().UnixNano()
	previousHashrate := int32(0)

	nonce = uint64(r.Int63())
	target := new(big.Int).Div(maxUint256, diff)
	var ctx unsafe.Pointer
	if is_lyra2 {
		ctx = C.LYRA2_create()
	} else {
		ctx = C.cryptonight_create()
	}
	headerBytes := types.HeaderToBytes(block.Header())
	for {
		select {
		case <-stop:
			atomic.AddInt32(&pow.hashRate, -previousHashrate)
			if is_lyra2 {
				C.LYRA2_destroy(ctx)
			} else {
				C.cryptonight_destroy(ctx)
			}
			return 0
		default:
			i++

			// we don't have to update hash rate on every nonce, so update after
			// first nonce check and then after 2^X nonces
			if i == 2 || ((i % (1 << 7)) == 0) {
				elapsed := time.Now().UnixNano() - start
				hashes := (float64(1e9) / float64(elapsed)) * float64(i)
				hashrateDiff := int32(hashes) - previousHashrate
				previousHashrate = int32(hashes)
				atomic.AddInt32(&pow.hashRate, hashrateDiff)
			}

			var result *big.Int
			if is_lyra2 {
				result = pow.computeLYRA2(ctx, headerBytes, nonce).Big()
			} else {
				result = pow.compute(ctx, headerBytes, nonce).Big()
			}

			// TODO: disagrees with the spec https://github.com/ethereum/wiki/wiki/Ethash#mining
			if result.Cmp(target) <= 0 {
				atomic.AddInt32(&pow.hashRate, -previousHashrate)
				if is_lyra2 {
					C.LYRA2_destroy(ctx)
				} else {
					C.cryptonight_destroy(ctx)
				}
				return nonce
			}
			nonce += 1
		}

		if !pow.turbo {
			time.Sleep(20 * time.Microsecond)
		}
	}
}

func (pow *Cryptonight) Turbo(on bool) {
	pow.turbo = on
}

func bytesToHash(in unsafe.Pointer) common.Hash {
	return *(*common.Hash)(in)
}

func (pow *Cryptonight) compute(ctx unsafe.Pointer, blockBytes []byte, nonce uint64) common.Hash {
	binary.BigEndian.PutUint64(blockBytes[len(blockBytes)-8:], nonce)

	var in unsafe.Pointer = C.CBytes(blockBytes)
	var out unsafe.Pointer = C.malloc(common.HashLength)
	if cpuid.CPU.AesNi() {
		C.cryptonight_hash_aesni(ctx, (*C.char)(in), (*C.char)(out), C.uint32_t(len(blockBytes)))
	} else {
		C.cryptonight_hash(ctx, (*C.char)(in), (*C.char)(out), C.uint32_t(len(blockBytes)))
	}

	var hash common.Hash = bytesToHash(unsafe.Pointer(out))

	C.free(in)
	C.free(out)

	return hash
}

func (pow *Cryptonight) computeLYRA2(lyra2_ctx unsafe.Pointer, blockBytes []byte, nonce uint64) common.Hash {
	binary.BigEndian.PutUint64(blockBytes[len(blockBytes)-8:], nonce)

	var in unsafe.Pointer = C.CBytes(blockBytes)
	var out unsafe.Pointer = C.malloc(common.HashLength)

	C.LYRA2(lyra2_ctx, unsafe.Pointer(out), common.HashLength, unsafe.Pointer(in), C.int32_t(len(blockBytes)))

	var hash common.Hash = bytesToHash(unsafe.Pointer(out))

	C.free(in)
	C.free(out)

	return hash
}

func (pow *Cryptonight) CalcHash(headerBytes []byte, nonce uint64) *big.Int {
	var cn_ctx unsafe.Pointer = C.cryptonight_create()
	result := pow.compute(cn_ctx, headerBytes, nonce)
	C.cryptonight_destroy(cn_ctx)

	return result.Big()
}

func (pow *Cryptonight) CalcHashLYRA2(headerBytes []byte, nonce uint64) *big.Int {
	var lyra2_ctx unsafe.Pointer = C.LYRA2_create()
	result := pow.computeLYRA2(lyra2_ctx, headerBytes, nonce)
	C.LYRA2_destroy(lyra2_ctx)

	return result.Big()
}

func (pow *Cryptonight) Verify(block pow.Block) bool {
	difficulty := block.Difficulty()
	headerBytes := types.HeaderToBytes(block.Header())

	/* Cannot happen if block header diff is validated prior to PoW, but can
		 happen if PoW is checked first due to parallel PoW checking.
	*/
	if difficulty.Cmp(common.Big0) == 0 {
		glog.V(logger.Debug).Infof("invalid block difficulty")
		return false
	}

	var result *big.Int
	if block.NumberU64() >= pow.lyra2_height {
		result = pow.CalcHashLYRA2(headerBytes, block.Nonce())
	} else {
		result = pow.CalcHash(headerBytes, block.Nonce())
	}

	// The actual check.
	target := new(big.Int).Div(maxUint256, difficulty)
	return result.Cmp(target) <= 0
}

func New(lyra2_height uint64) *Cryptonight {
	return &Cryptonight{lyra2_height: lyra2_height}
}

func NewShared(lyra2_height uint64) *Cryptonight {
	return &Cryptonight{lyra2_height: lyra2_height}
}

func NewForTesting(lyra2_height uint64) (*Cryptonight, error) {
	return &Cryptonight{lyra2_height: lyra2_height}, nil
}
