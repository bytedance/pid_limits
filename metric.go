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
	"sort"

	"github.com/bytedance/pid_limits/core/base"
	"github.com/bytedance/pid_limits/util"
)

var (
	ErrRate = CreateMetricFactory(func(entry *PlatoEntry) *Metric {
		return &Metric{
			cul: func() float64 {
				curMillis := util.CurrentTimeMillis()
				err := float64(entry.rts.GetSumWithTime(curMillis, base.MetricEventError))
				total := float64((entry.rts.GetSumWithTime(curMillis, base.MetricEventComplete)) + (entry.rts.GetSumWithTime(curMillis, base.MetricEventBlock)))
				if total == 0 {
					return 0
				} else {
					return err / total
				}
			},
		}
	})

	Qps = CreateMetricFactory(func(entry *PlatoEntry) *Metric {
		return &Metric{
			cul: func() float64 {
				intervalMs := entry.rts.Real.IntervalInMs()
				sampleCount := entry.rts.Real.SampleCount()
				now := util.CurrentTimeMillis()
				end := now - uint64(intervalMs/sampleCount)
				return entry.rts.GetAvgWithTime(base.MetricEventComplete, end-1000, end)
			},
		}
	})

	PctRT = CreateMetricFactory(func(entry *PlatoEntry) *Metric {
		return &Metric{
			cul: func() float64 {
				a := entry.rts.GetPctWithTime(util.CurrentTimeMillis(), base.MetricEventRt)
				sort.Ints(a)
				if len(a) == 0 {
					return 0
				}
				return float64(a[int(len(a)/10*9)])
			},
		}
	})

	AvgRT = CreateMetricFactory(func(entry *PlatoEntry) *Metric {
		return &Metric{
			cul: func() float64 {
				complete := float64(entry.rts.GetSum(base.MetricEventComplete))
				if complete == 0 {
					return 0
				} else {
					return float64(entry.rts.GetSum(base.MetricEventRt)) / complete
				}
			},
		}
	})
)

type Metric struct {
	value float64
	cul   func() float64
}

func (m *Metric) getValue() float64 {
	return util.GetFloat64(&m.value)
}

type MetricFactoryFunc func(*PlatoEntry) *Metric

type MetricFactory interface {
	CreateMetric(*PlatoEntry) *Metric
}

type defaultMetricFactory struct {
	fn MetricFactoryFunc
}

func (d *defaultMetricFactory) CreateMetric(pe *PlatoEntry) *Metric {
	return d.fn(pe)
}

//Helper function to generate MetricFactory
func CreateMetricFactory(mf MetricFactoryFunc) MetricFactory {
	return &defaultMetricFactory{fn: mf}
}
