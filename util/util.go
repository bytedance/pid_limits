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
package  util

import (
	"log"
	"math"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	TimeFormat         = "2006-01-02 15:04:05"
	DateFormat         = "2006-01-02"
	UnixTimeUnitOffset = uint64(time.Millisecond / time.Nanosecond)
)

func CurrentTimeMillis() uint64 {
	return uint64(time.Now().UnixNano()) / UnixTimeUnitOffset
}

func SetFloat64(addr *float64, new float64) float64 {
	unsafeAddr := (*uint64)(unsafe.Pointer(addr))

	for {
		oldValue := math.Float64bits(*addr)
		newValue := math.Float64bits(new)

		if atomic.CompareAndSwapUint64(unsafeAddr, oldValue, newValue) {
			return new
		}
	}
}

//Atomic load a float64 value
func GetFloat64(p *float64) float64 {
	if p == nil {
		return 0
	}
	return math.Float64frombits(atomic.LoadUint64((*uint64)(unsafe.Pointer(p))))
}

// Returns the current Unix timestamp in nanoseconds.
func CurrentTimeNano() uint64 {
	return uint64(time.Now().UnixNano())
}

func LoopWithInterval(runnable func(), interval time.Duration) {
	funcName := runtime.FuncForPC(reflect.ValueOf(runnable).Pointer()).Name()
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("error: [LoopWithInterval]%v(adaptive) panic, error : %v, stack : %s", funcName, err, debug.Stack())
				}
			}()
			runnable()
		}()
		time.Sleep(interval)
	}
}

func IsDocker() bool {
	return os.Getenv("IS_DOCKER_ENV") == "true"
}
