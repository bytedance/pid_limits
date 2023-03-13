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
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

type passRule struct{}

func (p *passRule) Decide(ctx *EntryCtx) bool {
	return true
}

type rejectRule struct{}

func (r *rejectRule) Decide(ctx *EntryCtx) bool {
	return false
}

func TestPCT(t *testing.T)  {
	entry := DefaultEntry("")
	Init([]*PlatoEntry{entry})
	go func() {
		for{
			entry.Run(func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(5)))
				return nil
			})
		}
	}()
	go func() {
		for{
			fmt.Println("AvgRT: ", Chain(entry, AvgRT))
			fmt.Println("PctRt:", Chain(entry, PctRT))
			fmt.Println("QPS: ", Chain(entry, Qps))
			time.Sleep(time.Second)
		}
	}()
}

func TestEntry(t *testing.T) {
	entry := DefaultEntry("")
	entry.Rule = &passRule{}
	//may pass
	var i int
	e, b := entry.Run(func() error {
		i = 1
		return nil
	})
	assert.True(t, e == nil)
	assert.True(t, b)
	assert.True(t, i == 1)

	//return error
	e, b = entry.Run(func() error {
		i = 2
		return errors.New("test")
	})
	assert.True(t, b)
	assert.True(t, i == 2)
	assert.True(t, e.Error() == "test")

	//reject
	entry.Rule = &rejectRule{}
	var i1 int
	e, b = entry.Run(func() error {
		i1 = 1
		return errors.New("test")
	})
	assert.True(t, e == ErrRejectByRule)
	assert.True(t, i1 == 0)
	assert.False(t, b)
}
