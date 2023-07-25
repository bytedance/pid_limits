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
package plato

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytedance/pid_limits/util"
	"github.com/bytedance/pid_limits/util/executors"
)

type metricsConainer interface {
	AddMetrics(...*Metric) int
	NextMetric() *Metric
}

func NewCalculateManager(workersCount int, queueSize int) (*CalculateManager, error) {
	p := executors.NewFixedSizeExecutorsPool(workersCount, queueSize)
	if p == nil {
		return nil, errors.New("fail to init workers pool")
	}

	return &CalculateManager{
		wp: p,
		mc: NewCopyOnWriteMetricsContainer(),
	}, nil
}

type CalculateManager struct {
	wp       executors.ExecutorPool
	mc       metricsConainer
	initOnce sync.Once
}

func (c *CalculateManager) AddMetrics(metrics ...*Metric) int {
	return c.mc.AddMetrics(metrics...)
}

func (c *CalculateManager) Start() {
	c.initOnce.Do(func() {
		go func() {
			util.LoopWithInterval(func() {
				m := c.mc.NextMetric()
				if m == nil {
					return
				}
				_ = c.wp.Run(func() error {
					util.SetFloat64(&m.value, m.cul())
					return nil
				})
			}, 10*time.Millisecond)
		}()
	})
}

func NewCopyOnWriteMetricsContainer() metricsConainer {
	c := &CopyOnWriteMetricsContainer{}
	c.l.Store(make([]*Metric, 0))
	c.em = map[*Metric]struct{}{}
	return c
}

type CopyOnWriteMetricsContainer struct {
	l  atomic.Value         // atomic value to hold metrics slice
	m  sync.Mutex           // mutex to protect all metrics will be write to l
	em map[*Metric]struct{} //avoid duplicate insert
	n  uint64               // next metric to calculate
}

func (c *CopyOnWriteMetricsContainer) NextMetric() *Metric {
	metrics := c.l.Load().([]*Metric)
	base := uint64(len(metrics))
	if base == 0 {
		return nil
	}
	idx := atomic.AddUint64(&c.n, 1) - 1
	return metrics[int(idx%base)]
}

func (c *CopyOnWriteMetricsContainer) AddMetrics(ms ...*Metric) int {
	if len(ms) == 0 {
		return 0
	}
	c.m.Lock()
	defer c.m.Unlock()

	//copy on write
	oldMetrics := c.l.Load().([]*Metric)
	newMetrics := make([]*Metric, len(oldMetrics), len(oldMetrics)+len(ms))
	if len(oldMetrics) != 0 {
		copy(newMetrics, oldMetrics)
	}

	var addNum int
	for _, m := range ms {
		if _, ok := c.em[m]; ok {
			continue
		}

		newMetrics = append(newMetrics, m)
		c.em[m] = struct{}{}
		addNum++
	}
	c.l.Store(newMetrics)
	return addNum
}
