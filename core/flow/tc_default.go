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
package  flow

import (
	"github.com/bytedance/plato/core/base"
	stat "github.com/bytedance/plato/core/stat/base"
	"github.com/bytedance/plato/util"
	"math"
	"sync/atomic"
	"time"
)

const (
	nanoUnitOffset    = time.Second / time.Nanosecond
	maxQueueingTimeNs = 0 * uint64(time.Millisecond/time.Nanosecond)
)

func DoCheck(acquireCount uint32, threshold float64, node *stat.SlidingWindowMetric) bool {
	curCount := node.GetQPS(base.MetricEventComplete)

	if curCount+float64(acquireCount) > threshold {
		return false
	}
	return true
}

func DoCheck2(acquireCount uint32, threshold float64, node *stat.SlidingWindowMetric) bool {
	// Pass when acquire count is less or equal than 0.
	if acquireCount <= 0 {
		return true
	}
	if threshold <= 0 {
		return false
	}
	// Here we use nanosecond so that we could control the queueing time more accurately.
	curNano := util.CurrentTimeNano()
	// The interval between two requests (in nanoseconds).
	interval := uint64(math.Ceil(float64(acquireCount) / threshold * float64(nanoUnitOffset)))

	// Expected pass time of this request.
	expectedTime := atomic.LoadUint64(&node.LastPassedTime) + interval

	if expectedTime <= curNano {
		// Contention may exist here, but it's okay.
		atomic.StoreUint64(&node.LastPassedTime, curNano)
		return true
	}
	estimatedQueueingDuration := atomic.LoadUint64(&node.LastPassedTime) + interval - util.CurrentTimeNano()
	if estimatedQueueingDuration > maxQueueingTimeNs {
		return false
	}

	oldTime := atomic.AddUint64(&node.LastPassedTime, interval)
	estimatedQueueingDuration = oldTime - util.CurrentTimeNano()
	if estimatedQueueingDuration > maxQueueingTimeNs {
		// Subtract the interval.
		atomic.AddUint64(&node.LastPassedTime, ^(interval - 1))
		return false
	}
	return true
}
