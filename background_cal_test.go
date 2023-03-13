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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetricsContainer(t *testing.T) {
	mc := NewCopyOnWriteMetricsContainer()
	entry := DefaultEntry("test")
	assert.True(t, mc.AddMetrics(entry.Metrics[Qps]) == 1)
	mc.AddMetrics(entry.Metrics[PctRT])
	mc.AddMetrics(entry.Metrics[AvgRT])
	assert.True(t, mc.AddMetrics(entry.Metrics[Qps]) == 0)
	assert.Equal(t, mc.NextMetric(), entry.Metrics[Qps])
	assert.Equal(t, mc.NextMetric(), entry.Metrics[PctRT])
	assert.Equal(t, mc.NextMetric(), entry.Metrics[AvgRT])
	assert.Equal(t, mc.NextMetric(), entry.Metrics[Qps])
}

func TestCalculateManager(t *testing.T) {
	manager, e := NewCalculateManager(10, 10)
	assert.Nil(t, e)
	entry := DefaultEntry("test")
	for _, m := range entry.Metrics {
		manager.AddMetrics(m)
	}

	manager.Start()

	for i := 0; i < 100; i++ {
		_, _ = entry.Run(func() error {
			return errors.New("err")
		})
		time.Sleep(3 * time.Millisecond)
	}

	assert.True(t, Chain(entry, Qps) > 0)
	assert.True(t, Chain(entry, ErrRate) > 0)
}
