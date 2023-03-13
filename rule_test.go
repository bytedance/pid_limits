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
package  plato

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bytedance/plato/util"

	"github.com/stretchr/testify/assert"
)

func TestChainAndCalculate(t *testing.T) {
	ce := NewEntryCalculated("test", Qps)
	for i := 0; i < 100; i++ {
		ce.Run(func() error {
			return errors.New("")
		})
		//sleep some time so background goroutine can calculate the metrics
	}
	time.Sleep(500 * time.Millisecond)
	v := Chain(ce, Qps)
	t.Log(v)
	assert.True(t, math.Abs(v-100) < 0.0001)

	v = Chain(ce, AvgRT)
	assert.True(t, v < 0.0001)

	ce = NewPlatoEntry("test_no_cache", Qps)
	for i := 0; i < 100; i++ {
		ce.Run(func() error {
			return errors.New("")
		})
	}
	time.Sleep(10 * time.Millisecond)
	v = Chain(ce, Qps)
	assert.Equal(t, v, 0.0)
	v = Calculate(ce, Qps)
	assert.True(t, math.Abs(v-100) < 0.0001)
	v = Calculate(ce, AvgRT)
	assert.Equal(t, v, 0.0)

	ce = DefaultEntryCalculated("test_default_entry")
	for i := 0; i < 100; i++ {
		ce.Run(func() error {
			time.Sleep(1 * time.Millisecond)
			return errors.New("")
		})
	}
	time.Sleep(400 * time.Millisecond)
	v = Chain(ce, Qps)
	t.Log(v)
	assert.True(t, math.Abs(v-100) < 0.0001)
	v = Calculate(ce, Qps)
	t.Log(v)
	assert.True(t, math.Abs(v-100) < 0.0001)
	v = Calculate(ce, AvgRT)
	t.Log(v)
	assert.True(t, v > 0)
}

const qps = 1000

func BenchmarkRand2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < qps; i++ {
			go rand.Intn(1000)
		}
	}
}

func BenchmarkRand3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < qps; i++ {
			rand.New(rand.NewSource(int64(time.Nanosecond))).Int()
			rand.Intn(100)
		}
	}
}

func BenchmarkFastRand(b *testing.B) {
	b.SetParallelism(4)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 所有 goroutine 一起，循环一共执行 b.N 次
			_ = util.Uint32n(10000) < 5000
		}
	})
}

//_ = rand.Intn(10000) < 50
func BenchmarkRand(b *testing.B) {
	b.SetParallelism(4)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 所有 goroutine 一起，循环一共执行 b.N 次
			_ = rand.Intn(10000) < 5000
		}
	})
}

func BenchmarkNano(b *testing.B) {
	b.SetParallelism(4)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 所有 goroutine 一起，循环一共执行 b.N 次
			_ = time.Now().Nanosecond()%10000 < 5000
		}
	})
}

func BenchmarkAtomicValue(b *testing.B) {
	b.SetParallelism(4)
	a := atomic.Value{}
	a.Store(1.2)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 所有 goroutine 一起，循环一共执行 b.N 次
			_, _ = a.Load().(float64)
		}
	})
}

func BenchmarkCAS(b *testing.B) {
	b.SetParallelism(4)
	rate := 1.2
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 所有 goroutine 一起，循环一共执行 b.N 次
			util.GetFloat64(&rate)
		}
	})
}

func BenchmarkFMT(b *testing.B) {
	b.SetParallelism(4)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fmt.Errorf("block by pid")
		}
	})
}

func TestRuntime(t *testing.T) {
	fmt.Println(runtime.GOMAXPROCS(0))
}

func TestBB(t *testing.T) {
	fmt.Printf("%v\n", len([]byte("ababfdfdasb")))
}

func TestMarshal(t *testing.T) {
	a := math.MaxInt64
	fmt.Println(len(strconv.Itoa(a)))
	a2 := strconv.FormatInt(int64(a), 36)
	fmt.Println(a2)
	fmt.Println(len(a2))
}

func BenchmarkFMT2(b *testing.B) {
	b.SetParallelism(4)
	a := math.MaxInt64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = strconv.FormatInt(int64(a), 36)
		}
	})
}

type USER struct {
	name string
}

func TestPointer(t *testing.T) {
	u := &USER{name: "aaa"}
	fix(u)
	fmt.Println(u)
}

func fix(u *USER) {
	u.name = "bbbb"
}

func BenchmarkAtomic(b *testing.B) {
	b.SetParallelism(4)
	a := atomic.Value{}
	a.Store(1)
	util.LoopWithInterval(func() {
		a.Store(0.3)
	}, time.Millisecond*100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			a.Load()
		}
	})
}

func BenchmarkUint32(b *testing.B) {
	b.SetParallelism(4)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000)
		}
	})
}

func TestABC(t *testing.T) {
	a := -8000
	fmt.Println(uint32(a))
	fmt.Println(-uint32(a))
}
