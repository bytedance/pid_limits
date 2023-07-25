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
package limiting

import (
	"log"
	"math"
	"sync/atomic"
	"time"

	"github.com/bytedance/pid_limits/arithmetic/pid"
	"github.com/bytedance/pid_limits/metrics/system/cpu"
	"github.com/bytedance/pid_limits/util"
)

type PIDLimiting struct {
	rate                uint32
	pid                 *pid.PID
	monitor             cpu.Monitor
	enableMetric        bool
	enableOverloadScene bool
	enablePid           atomic.Value
}

func (l *PIDLimiting) Limit() bool {
	if l.monitor.IsOverload() && util.Uint32n(10000) < atomic.LoadUint32(&l.rate) {
		return true
	}
	return false
}

// Rate the probability is form 0 ~ 10000
func (l *PIDLimiting) LimitRatio() float64 {
	if l.monitor.IsOverload() {
		return math.Min(10000, math.Max(0, float64(atomic.LoadUint32(&l.rate))))
	}
	return 0
}

func (l *PIDLimiting) start() {
	go util.LoopWithInterval(func() {
		cpuUsage := cpu.GetUsage()
		if !l.enableOverloadScene {
			cpuUsage = math.Min(cpuUsage, 1)
		}
		rate := l.pid.Compute(cpuUsage)
		atomic.StoreUint32(&l.rate, uint32(-rate))
		if l.enableMetric {
			log.Printf(
				"cpu usage is: %f, threshould is: %f, reject rate is %f, overloaded: %v",
				cpuUsage, l.pid.GetThreshold(), -rate, l.monitor.IsOverload(),
			)
		}
	}, 100*time.Millisecond)
}
