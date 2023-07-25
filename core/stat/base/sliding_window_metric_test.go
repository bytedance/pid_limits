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
	"github.com/bytedance/pid_limits/util"
	"reflect"
	"testing"
	"time"

	"github.com/bytedance/pid_limits/core/base"
)

func TestSlidingWindowMetric_getBucketStartRange(t *testing.T) {
	sw := &SlidingWindowMetric{
		Real:           NewBucketLeapArray(20, 2000),
		LastPassedTime: 0,
	}
	time.Sleep(40 * time.Millisecond)
	type fields struct {
		LastPassedTime uint64
		Real           *BucketLeapArray
	}
	type args struct {
		timeMs uint64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantStart uint64
		wantEnd   uint64
	}{
		{
			"test",
			fields{
				LastPassedTime: sw.LastPassedTime,
				Real:           sw.Real,
			},
			args{timeMs: util.CurrentTimeMillis()},
			calculateStartTime(util.CurrentTimeMillis(), sw.Real.IntervalInMs()/sw.Real.SampleCount()) - uint64(sw.Real.IntervalInMs()) + uint64(sw.Real.IntervalInMs()/sw.Real.SampleCount()),
			calculateStartTime(util.CurrentTimeMillis(), sw.Real.IntervalInMs()/sw.Real.SampleCount()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlidingWindowMetric{
				LastPassedTime: tt.fields.LastPassedTime,
				Real:           tt.fields.Real,
			}
			gotStart, gotEnd := m.getBucketStartRange(tt.args.timeMs)
			if gotStart != tt.wantStart {
				t.Errorf("SlidingWindowMetric.getBucketStartRange() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("SlidingWindowMetric.getBucketStartRange() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}

func TestSlidingWindowMetric_getBucketRange(t *testing.T) {
	sw := &SlidingWindowMetric{
		Real:           NewBucketLeapArray(200, 2000),
		LastPassedTime: 0,
	}
	type fields struct {
		LastPassedTime uint64
		Real           *BucketLeapArray
	}
	type args struct {
		startMs uint64
		endMs   uint64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantStart uint64
		wantEnd   uint64
	}{
		{
			"test",
			fields{
				LastPassedTime: sw.LastPassedTime,
				Real:           sw.Real,
			},
			args{
				startMs: 12,
				endMs:   43,
			},
			10,
			40,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlidingWindowMetric{
				LastPassedTime: tt.fields.LastPassedTime,
				Real:           tt.fields.Real,
			}
			gotStart, gotEnd := m.getBucketRange(tt.args.startMs, tt.args.endMs)
			if gotStart != tt.wantStart {
				t.Errorf("SlidingWindowMetric.getBucketRange() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("SlidingWindowMetric.getBucketRange() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}

func TestSlidingWindowMetric_count(t *testing.T) {
	sw := &SlidingWindowMetric{
		Real:           NewBucketLeapArray(200, 2000),
		LastPassedTime: 0,
	}
	sw.Real.AddCount(base.MetricEventComplete, 2)
	type fields struct {
		LastPassedTime uint64
		Real           *BucketLeapArray
	}
	type args struct {
		event  base.MetricEvent
		values []*bucketWrap
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			"test",
			fields{
				LastPassedTime: sw.LastPassedTime,
				Real:           sw.Real,
			},
			args{
				event:  base.MetricEventComplete,
				values: sw.Real.Values(util.CurrentTimeMillis()),
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlidingWindowMetric{
				LastPassedTime: tt.fields.LastPassedTime,
				Real:           tt.fields.Real,
			}
			if got := m.count(tt.args.event, tt.args.values); got != tt.want {
				t.Errorf("SlidingWindowMetric.count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindowMetric_pct(t *testing.T) {
	sw := &SlidingWindowMetric{
		LastPassedTime: 0,
		Real:           NewBucketLeapArray(200, 2000),
	}
	sw.Real.AddCount(base.MetricEventRt, 1)
	sw.Real.AddCount(base.MetricEventRt, 2)
	sw.Real.AddCount(base.MetricEventRt, 3)
	sw.Real.AddCount(base.MetricEventRt, 4)
	time.Sleep(10 * time.Millisecond)
	sw.Real.AddCount(base.MetricEventRt, 9)
	type fields struct {
		LastPassedTime uint64
		Real           *BucketLeapArray
	}
	type args struct {
		event  base.MetricEvent
		values []*bucketWrap
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int
	}{
		{
			"test",
			fields{
				LastPassedTime: sw.LastPassedTime,
				Real:           sw.Real,
			},
			args{
				event:  base.MetricEventComplete,
				values: sw.Real.Values(util.CurrentTimeMillis()),
			},
			[]int{1, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlidingWindowMetric{
				LastPassedTime: tt.fields.LastPassedTime,
				Real:           tt.fields.Real,
			}
			if got := m.pct(tt.args.event, tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SlidingWindowMetric.pct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindowMetric_GetSum(t *testing.T) {
	sw := &SlidingWindowMetric{
		LastPassedTime: 0,
		Real:           NewBucketLeapArray(200, 2000),
	}
	sw.Real.AddCount(base.MetricEventComplete, 1)
	sw.Real.AddCount(base.MetricEventComplete, 2)
	sw.Real.AddCount(base.MetricEventComplete, 3)
	type fields struct {
		LastPassedTime uint64
		Real           *BucketLeapArray
	}
	type args struct {
		event base.MetricEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			"test",
			fields{
				LastPassedTime: 0,
				Real:           sw.Real,
			},
			args{event: base.MetricEventComplete},
			6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlidingWindowMetric{
				LastPassedTime: tt.fields.LastPassedTime,
				Real:           tt.fields.Real,
			}
			if got := m.GetSum(tt.args.event); got != tt.want {
				t.Errorf("SlidingWindowMetric.GetSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindowMetric_GetValuesWithTime(t *testing.T) {
	sw := &SlidingWindowMetric{
		LastPassedTime: 0,
		Real:           NewBucketLeapArray(200, 2000),
	}
	sw.Real.AddCount(base.MetricEventComplete, 1)
	sw.Real.AddCount(base.MetricEventComplete, 2)
	time.Sleep(20 * time.Millisecond)
	sw.Real.AddCount(base.MetricEventComplete, 2)
	type fields struct {
		LastPassedTime uint64
		Real           *BucketLeapArray
	}
	type args struct {
		now   uint64
		event base.MetricEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int64
	}{
		{
			"test",
			fields{
				LastPassedTime: 0,
				Real:           sw.Real,
			},
			args{
				now:   util.CurrentTimeMillis(),
				event: base.MetricEventComplete,
			},
			[]int64{3, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlidingWindowMetric{
				LastPassedTime: tt.fields.LastPassedTime,
				Real:           tt.fields.Real,
			}
			if got := m.GetValuesWithTime(tt.args.now, tt.args.event); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SlidingWindowMetric.GetValuesWithTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindowMetric_GetAvgWithTime(t *testing.T) {
	sw := &SlidingWindowMetric{
		LastPassedTime: 0,
		Real:           NewBucketLeapArray(200, 2000),
	}
	sw.Real.AddCount(base.MetricEventComplete, 10)
	sw.Real.AddCount(base.MetricEventComplete, 20)
	time.Sleep(20 * time.Millisecond)
	sw.Real.AddCount(base.MetricEventComplete, 20)
	type fields struct {
		LastPassedTime uint64
		Real           *BucketLeapArray
	}
	type args struct {
		event base.MetricEvent
		start uint64
		end   uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			"test",
			fields{
				LastPassedTime: 0,
				Real:           sw.Real,
			},
			args{
				event: base.MetricEventComplete,
				start: util.CurrentTimeMillis() - 30,
				end:   util.CurrentTimeMillis() + 20,
			},
			1000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlidingWindowMetric{
				LastPassedTime: tt.fields.LastPassedTime,
				Real:           tt.fields.Real,
			}
			if got := m.GetAvgWithTime(tt.args.event, tt.args.start, tt.args.end); got != tt.want {
				t.Errorf("SlidingWindowMetric.GetAvgWithTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindowMetric_GetSumWithTime(t *testing.T) {
	sw := &SlidingWindowMetric{
		LastPassedTime: 0,
		Real:           NewBucketLeapArray(200, 2000),
	}
	sw.Real.AddCount(base.MetricEventComplete, 10)
	sw.Real.AddCount(base.MetricEventComplete, 20)
	time.Sleep(20 * time.Millisecond)
	sw.Real.AddCount(base.MetricEventComplete, 20)
	type fields struct {
		LastPassedTime uint64
		Real           *BucketLeapArray
	}
	type args struct {
		now   uint64
		event base.MetricEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			"test",
			fields{
				LastPassedTime: 0,
				Real:           sw.Real,
			},
			args{
				now:   util.CurrentTimeMillis(),
				event: base.MetricEventComplete,
			},
			50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlidingWindowMetric{
				LastPassedTime: tt.fields.LastPassedTime,
				Real:           tt.fields.Real,
			}
			if got := m.GetSumWithTime(tt.args.now, tt.args.event); got != tt.want {
				t.Errorf("SlidingWindowMetric.GetSumWithTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindowMetric_GetPctWithTime(t *testing.T) {
	sw := &SlidingWindowMetric{
		LastPassedTime: 0,
		Real:           NewBucketLeapArray(200, 2000),
	}
	sw.Real.AddCount(base.MetricEventRt, 10)
	sw.Real.AddCount(base.MetricEventRt, 20)
	time.Sleep(20 * time.Millisecond)
	sw.Real.AddCount(base.MetricEventRt, 20)
	type fields struct {
		LastPassedTime uint64
		Real           *BucketLeapArray
	}
	type args struct {
		now   uint64
		event base.MetricEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int
	}{
		{
			"test",
			fields{
				LastPassedTime: 0,
				Real:           sw.Real,
			},
			args{
				now:   util.CurrentTimeMillis(),
				event: base.MetricEventComplete,
			},
			[]int{10, 20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlidingWindowMetric{
				LastPassedTime: tt.fields.LastPassedTime,
				Real:           tt.fields.Real,
			}
			if got := m.GetPctWithTime(tt.args.now, tt.args.event); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SlidingWindowMetric.GetPctWithTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindowMetric_GetQPS(t *testing.T) {
	sw := &SlidingWindowMetric{
		LastPassedTime: 0,
		Real:           NewBucketLeapArray(200, 2000),
	}
	sw.Real.AddCount(base.MetricEventComplete, 10)
	sw.Real.AddCount(base.MetricEventComplete, 20)
	time.Sleep(20 * time.Millisecond)
	sw.Real.AddCount(base.MetricEventComplete, 20)
	type fields struct {
		LastPassedTime uint64
		Real           *BucketLeapArray
	}
	type args struct {
		event base.MetricEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			"test",
			fields{
				LastPassedTime: 0,
				Real:           sw.Real,
			},
			args{event: base.MetricEventComplete},
			25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlidingWindowMetric{
				LastPassedTime: tt.fields.LastPassedTime,
				Real:           tt.fields.Real,
			}
			if got := m.GetQPS(tt.args.event); got != tt.want {
				t.Errorf("SlidingWindowMetric.GetQPS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindowMetric_GetQPSWithTime(t *testing.T) {
	type fields struct {
		LastPassedTime uint64
		Real           *BucketLeapArray
	}
	type args struct {
		now   uint64
		event base.MetricEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlidingWindowMetric{
				LastPassedTime: tt.fields.LastPassedTime,
				Real:           tt.fields.Real,
			}
			if got := m.GetQPSWithTime(tt.args.now, tt.args.event); got != tt.want {
				t.Errorf("SlidingWindowMetric.GetQPSWithTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindowMetric_getIntervalInSecond(t *testing.T) {
	type fields struct {
		LastPassedTime uint64
		Real           *BucketLeapArray
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SlidingWindowMetric{
				LastPassedTime: tt.fields.LastPassedTime,
				Real:           tt.fields.Real,
			}
			if got := m.getIntervalInSecond(); got != tt.want {
				t.Errorf("SlidingWindowMetric.getIntervalInSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}
