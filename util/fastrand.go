/*
 * Copyright 2021 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package util

import (
	"sync"
	"time"
)

var pool sync.Pool

type RandPool struct {
	x uint32
}

func Uint32() uint32 {
	v := pool.Get()
	if v == nil {
		v = &RandPool{}
	}
	r := v.(*RandPool)
	x := r.Uint32()
	pool.Put(r)
	return x
}

func Uint32n(maxN uint32) uint32 {
	x := Uint32()
	return uint32((uint64(x) * uint64(maxN)) >> 32)
}


func (r *RandPool) Uint32() uint32 {
	for r.x == 0 {
		r.x = getRandomUint32()
	}
	x := r.x
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	r.x = x
	return x
}

func (r *RandPool) Uint32n(maxN uint32) uint32 {
	x := r.Uint32()
	return uint32((uint64(x) * uint64(maxN)) >> 32)
}

func getRandomUint32() uint32 {
	x := time.Now().UnixNano()
	return uint32((x >> 32) ^ x)
}
