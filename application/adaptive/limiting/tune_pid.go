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
	"github.com/bytedance/pid_limits/arithmetic/pid"
	"github.com/bytedance/pid_limits/metrics/system/cpu"
	"github.com/bytedance/pid_limits/util"
	"sync/atomic"
	"time"
)

type PIDTune struct {
	rate  uint32
	tuner *pid.Tuner
}

func (t *PIDTune) start() {
	go func() {
		util.LoopWithInterval(func() {
			// async to calculate cup rate
			cpuUsage := cpu.GetUsage()
			rate := t.tuner.TunePID(cpuUsage)
			atomic.StoreUint32(&t.rate, uint32(rate))
		}, 100*time.Millisecond)
	}()
}

func (t *PIDTune) Limit() bool {
	return util.Uint32n(10000) < -atomic.LoadUint32(&t.rate)
}

// Rate the probability is form 0 ~ 10000
func (t *PIDTune) LimitRatio() uint32 {
	return -atomic.LoadUint32(&t.rate)
}
