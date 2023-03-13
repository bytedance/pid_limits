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
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetricMap(t *testing.T) {
	pe := NewPlatoEntry("")
	pe.AddMetric(PctRT)
	assert.True(t, pe.Metrics[PctRT] != nil)
}

func TestQps(t *testing.T) {
	pe := DefaultEntry("")
	Init([]*PlatoEntry{pe})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			t := time.Tick(20 * time.Millisecond)
			pe.Run(func() error {
				return nil
			})
			for {
				select {
				case <-t:
					pe.Run(func() error {
						return nil
					})
				}
			}
		}()
		wg.Done()
	}
	wg.Wait()
}

func TestErrRate(t *testing.T) {
	pe := DefaultEntry("")
	Init([]*PlatoEntry{DefaultEntry(""), DefaultEntry(""), pe})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			t := time.Tick(20 * time.Millisecond)
			for {
				select {
				case <-t:
					pe.Run(func() error {
						if rand.Intn(100) < 50 {
							return errors.New("111")
						}
						return nil
					})
				}
			}
		}()
		wg.Done()
	}
	wg.Wait()
}
