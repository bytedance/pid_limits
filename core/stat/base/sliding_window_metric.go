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
package base

import (
	"log"

	"github.com/bytedance/pid_limits/core/base"
	"github.com/bytedance/pid_limits/util"
)

// SlidingWindowMetric represents the sliding window metric wrapper.
// It does not store any data and is the wrapper of BucketLeapArray to adapt to different internal bucket
// SlidingWindowMetric is used for SentinelRules and BucketLeapArray is used for monitor
// BucketLeapArray is per resource, and SlidingWindowMetric support only read operation.
type SlidingWindowMetric struct {
	LastPassedTime uint64
	Real           *BucketLeapArray
}

// Get the start time range of the bucket for the provided time.
// The actual time span is: [start, end + in.bucketTimeLength)
func (m *SlidingWindowMetric) getBucketStartRange(timeMs uint64) (start, end uint64) {
	curBucketStartTime := calculateStartTime(timeMs, m.Real.IntervalInMs()/m.Real.SampleCount())
	end = curBucketStartTime
	start = end - uint64(m.Real.IntervalInMs()) + uint64(m.Real.IntervalInMs()/m.Real.SampleCount())
	return
}

func (m *SlidingWindowMetric) getBucketRange(startMs, endMs uint64) (start, end uint64) {
	defaultBucketInMs := m.Real.IntervalInMs() / m.Real.SampleCount()
	start = calculateStartTime(startMs, defaultBucketInMs)
	end = calculateStartTime(endMs, defaultBucketInMs)
	return
}

func (m *SlidingWindowMetric) count(event base.MetricEvent, values []*bucketWrap) int64 {
	ret := int64(0)
	for _, ww := range values {
		mb := ww.value.Load()
		if mb == nil {
			log.Println("error: Illegal state: current bucket value is nil when summing count")
			continue
		}
		counter, ok := mb.(*MetricBucket)
		if !ok {
			log.Printf("error: Fail to cast data value(%+v) to MetricBucket type", mb)
			continue
		}
		ret += counter.Get(event)
	}
	return ret
}

func (m *SlidingWindowMetric) pct(event base.MetricEvent, values []*bucketWrap) []int {
	ret := make([]int, 0)
	for _, ww := range values {
		mb := ww.value.Load()
		if mb == nil {
			log.Println("error: Illegal state: current bucket value is nil when summing count")
			continue
		}
		counter, ok := mb.(*MetricBucket)
		if !ok {
			log.Printf("error: Fail to cast data value(%+v) to MetricBucket type", mb)
			continue
		}
		if c := int(counter.MinRt()); c != 0 && c != int(base.DefaultStatisticMaxRt) {
			ret = append(ret, c)
		}
	}
	return ret
}

func (m *SlidingWindowMetric) GetSum(event base.MetricEvent) int64 {
	return m.GetSumWithTime(util.CurrentTimeMillis(), event)
}

// Enter window sequence in chronological order
func (m *SlidingWindowMetric) GetValuesWithTime(now uint64, event base.MetricEvent) []int64 {
	record := make([]int64, 0)
	for _, ww := range m.Real.Values(now) {
		mb := ww.value.Load()
		if mb == nil {
			log.Println("error: Current bucket's value is nil.")
			continue
		}
		counter, ok := mb.(*MetricBucket)
		if !ok {
			log.Printf("error: Fail to cast data value(%+v) to MetricBucket type", mb)
			continue
		}
		v := counter.Get(event)
		if v != 0 {
			record = append(record, v)
		}
	}
	return record
}

func (m *SlidingWindowMetric) GetAvgWithTime(event base.MetricEvent, start, end uint64) float64 {
	//todo : border check
	if end < start || (end-start) < uint64(m.Real.IntervalInMs()/m.Real.SampleCount()) {
		return 0
	}

	start, end = m.getBucketRange(start, end)
	satisfiedBuckets := m.Real.ValuesConditional(end, func(ws uint64) bool {
		return ws >= start && ws <= end
	})
	total := m.count(event, satisfiedBuckets)
	return float64(total) / (float64(end-start) / 1000)
}

func (m *SlidingWindowMetric) GetSumWithTime(now uint64, event base.MetricEvent) int64 {
	start, end := m.getBucketStartRange(now)
	satisfiedBuckets := m.Real.ValuesConditional(now, func(ws uint64) bool {
		return ws >= start && ws <= end
	})
	return m.count(event, satisfiedBuckets)
}

func (m *SlidingWindowMetric) GetPctWithTime(now uint64, event base.MetricEvent) []int {
	start, end := m.getBucketStartRange(now)
	satisfiedBuckets := m.Real.ValuesConditional(now, func(ws uint64) bool {
		return ws >= start && ws <= end
	})
	return m.pct(event, satisfiedBuckets)
}

func (m *SlidingWindowMetric) GetQPS(event base.MetricEvent) float64 {
	return m.GetQPSWithTime(util.CurrentTimeMillis(), event)
}

func (m *SlidingWindowMetric) GetQPSWithTime(now uint64, event base.MetricEvent) float64 {
	return float64(m.GetSumWithTime(now, event)) / m.getIntervalInSecond()
}

func (m *SlidingWindowMetric) getIntervalInSecond() float64 {
	return m.Real.GetIntervalInSecond()
}
