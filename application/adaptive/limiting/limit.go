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
package  limiting

import (
	"math"

	"github.com/bytedance/plato/application/adaptive/config"
	"github.com/bytedance/plato/arithmetic/pid"
	"github.com/bytedance/plato/metrics/system/cpu"
)

// develop to orient interface, use limit() function to determine weather limit cpu rate
type RateLimit interface {
	Limit() bool
	LimitRatio() float64
}

func NewPidLimitingHttpDefault(cpuThreshold float64, opts ...config.OptionFunc) RateLimit {
	option := config.NewOptions()
	for _, opt := range opts {
		opt(option)
	}
	if option.EnableOverloadScene {
		return NewPidLimiting(5130.083602420542, 44.491571338644654, 123658.09189447836, cpuThreshold, opts...)
	}
	return NewPidLimiting(5351.821461335851, 12.030101184005932, 0.03, cpuThreshold, opts...)
}


func NewPidLimiting(kp float64, ki float64, kd float64, setPoint float64, opts ...config.OptionFunc) *PIDLimiting {
	option := config.NewOptions()
	for _, opt := range opts {
		opt(option)
	}
	var monitor cpu.Monitor
	var upperBound, lowerBound func() float64
	if option.DynamicPoint != nil {
		upperBound = func() float64 {
			return math.Min(0.99, option.DynamicPoint()+option.Drift)
		}
		lowerBound = func() float64 {
			return math.Max(0.01, option.DynamicPoint()-option.Drift)
		}
	} else {
		upper := math.Min(0.99, setPoint+option.Drift)
		upperBound = func() float64 {
			return upper
		}
		lower := math.Max(0.01, setPoint-option.Drift)
		lowerBound = func() float64 {
			return lower
		}
	}
	monitor = cpu.NewCPUMonitor(cpu.WithUpperBound(upperBound), cpu.WithLowerBound(lowerBound), cpu.WithAlg(option.MonitorAlg))
	limit := &PIDLimiting{
		rate:                0,
		pid:                 pid.SetTunings(kp, ki, kd, setPoint, pid.WithDynamicPoint(option.DynamicPoint)),
		monitor:             monitor,
		enableMetric:        option.EnableMetric,
		enableOverloadScene: option.EnableOverloadScene,
	}
	limit.start()
	limit.enablePid.Store(true)
	return limit
}
