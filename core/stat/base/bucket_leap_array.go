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
package  base

import (
	"fmt"
	"log"
	"sync/atomic"

	"github.com/bytedance/plato/core/base"
	"github.com/bytedance/plato/util"
	"github.com/pkg/errors"
)

// The implementation of sliding window based on LeapArray (as the sliding window infrastructure)
// and MetricBucket (as the data type). The MetricBucket is used to record statistic
// metrics per minimum time unit (i.e. the bucket time span).
type BucketLeapArray struct {
	data     leapArray
	dataType string
}

// sampleCount is the number of slots
// intervalInMs is the time length of sliding window
func NewBucketLeapArray(sampleCount uint32, intervalInMs uint32) *BucketLeapArray {
	if intervalInMs%sampleCount != 0 {
		panic(fmt.Sprintf("Invalid parameters, intervalInMs is %d, sampleCount is %d.", intervalInMs, sampleCount))
	}
	bucketLengthInMs := intervalInMs / sampleCount
	ret := &BucketLeapArray{
		data: leapArray{
			bucketLengthInMs: bucketLengthInMs,
			sampleCount:      sampleCount,
			intervalInMs:     intervalInMs,
			array:            nil,
		},
		dataType: "MetricBucket",
	}
	arr := newAtomicBucketWrapArray(int(sampleCount), bucketLengthInMs, ret)
	ret.data.array = arr
	return ret
}

func (bla *BucketLeapArray) SampleCount() uint32 {
	return bla.data.sampleCount
}

func (bla *BucketLeapArray) IntervalInMs() uint32 {
	return bla.data.intervalInMs
}

func (bla *BucketLeapArray) BucketLengthInMs() uint32 {
	return bla.data.bucketLengthInMs
}

func (bla *BucketLeapArray) DataType() string {
	return bla.dataType
}

func (bla *BucketLeapArray) GetIntervalInSecond() float64 {
	return float64(bla.IntervalInMs()) / 1000.0
}

func (bla *BucketLeapArray) newEmptyBucket() interface{} {
	return NewMetricBucket()
}

func (bla *BucketLeapArray) resetBucketTo(ww *bucketWrap, startTime uint64) *bucketWrap {
	atomic.StoreUint64(&ww.bucketStart, startTime)
	ww.value.Store(NewMetricBucket())
	return ww
}

// Write method
// It might panic
func (bla *BucketLeapArray) AddCount(event base.MetricEvent, count int64) {
	bla.addCountWithTime(util.CurrentTimeMillis(), event, count)
}

func (bla *BucketLeapArray) addCountWithTime(now uint64, event base.MetricEvent, count int64) {
	curBucket, err := bla.data.currentBucketOfTime(now, bla)
	if err != nil {
		log.Printf("error: Failed to get current bucket, current ts=%d, err: %+v.", now, errors.WithStack(err))
		return
	}
	if curBucket == nil {
		log.Println("error: Failed to add count: current bucket is nil")
		return
	}
	mb := curBucket.value.Load()
	if mb == nil {
		log.Println("error: Failed to add count: current bucket atomic value is nil")
		return
	}
	b, ok := mb.(*MetricBucket)
	if !ok {
		log.Println("Failed to add count: bucket data type error")
		return
	}
	b.Add(event, count)
}

// Read method, need to adapt upper application
// it might panic
func (bla *BucketLeapArray) Count(event base.MetricEvent) int64 {
	return bla.CountWithTime(util.CurrentTimeMillis(), event)
}

func (bla *BucketLeapArray) CountWithTime(now uint64, event base.MetricEvent) int64 {
	_, err := bla.data.currentBucketOfTime(now, bla)
	if err != nil {
		log.Printf("error: Fail to get current bucket, err: %+v.", errors.WithStack(err))
	}
	count := int64(0)
	for _, ww := range bla.data.valuesWithTime(now) {
		mb := ww.value.Load()
		if mb == nil {
			log.Println("error: Current bucket's value is nil.")
			continue
		}
		b, ok := mb.(*MetricBucket)
		if !ok {
			log.Println("error: Fail to assert MetricBucket type.")
			continue
		}
		count += b.Get(event)
	}
	return count
}

// Read method, get all bucketWrap.
func (bla *BucketLeapArray) Values(now uint64) []*bucketWrap {
	_, err := bla.data.currentBucketOfTime(now, bla)
	if err != nil {
		log.Printf("error: Fail to get current(%d) bucket, err: %+v.", now, errors.WithStack(err))
	}
	return bla.data.valuesWithTime(now)
}

func (bla *BucketLeapArray) ValuesConditional(now uint64, predicate base.TimePredicate) []*bucketWrap {
	_, err := bla.data.currentBucketOfTime(now, bla)
	if err != nil {
		log.Printf("error: Fail to get current(%d) bucket, err: %+v.", now, errors.WithStack(err))
	}
	return bla.data.ValuesConditional(now, predicate)
}

func (bla *BucketLeapArray) MinRt() int64 {
	_, err := bla.data.currentBucket(bla)
	if err != nil {
		log.Printf("error: Fail to get current bucket, err: %+v.", errors.WithStack(err))
	}

	ret := base.DefaultStatisticMaxRt

	for _, v := range bla.data.values() {
		mb := v.value.Load()
		if mb == nil {
			log.Println("error: Current bucket's value is nil.")
			continue
		}
		b, ok := mb.(*MetricBucket)
		if !ok {
			log.Println("error: Fail to cast data as MetricBucket type")
			continue
		}
		mr := b.MinRt()
		if ret > mr {
			ret = mr
		}
	}
	return ret
}
