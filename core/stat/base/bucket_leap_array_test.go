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
	"github.com/agiledragon/gomonkey"
	"github.com/bytedance/plato/util"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bytedance/plato/core/base"
)

func TestNewBucketLeapArray(t *testing.T) {
	defer func() {
		//恢复程序的控制权
		_ = recover()
	}()
	type args struct {
		sampleCount  uint32
		intervalInMs uint32
	}
	tests := []struct {
		name string
		args args
		want *BucketLeapArray
	}{
		{
			"test1",
			args{
				sampleCount:  100,
				intervalInMs: 200,
			},
			NewBucketLeapArray(100, 200),
		},
		{
			"test2",
			args{
				sampleCount:  101,
				intervalInMs: 200,
			},
			NewBucketLeapArray(101, 200),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBucketLeapArray(tt.args.sampleCount, tt.args.intervalInMs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBucketLeapArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketLeapArray_SampleCount(t *testing.T) {
	tests := []struct {
		name   string
		fields *BucketLeapArray
		want   uint32
	}{
		// TODO: Add test cases.
		{
			"test",
			NewBucketLeapArray(200, 2000),
			200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := &BucketLeapArray{
				data:     tt.fields.data,
				dataType: tt.fields.dataType,
			}
			if got := bla.SampleCount(); got != tt.want {
				t.Errorf("BucketLeapArray.SampleCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketLeapArray_IntervalInMs(t *testing.T) {
	tests := []struct {
		name   string
		fields *BucketLeapArray
		want   uint32
	}{
		{
			"test",
			NewBucketLeapArray(200, 2000),
			2000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := &BucketLeapArray{
				data:     tt.fields.data,
				dataType: tt.fields.dataType,
			}
			if got := bla.IntervalInMs(); got != tt.want {
				t.Errorf("BucketLeapArray.IntervalInMs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketLeapArray_BucketLengthInMs(t *testing.T) {
	tests := []struct {
		name   string
		fields *BucketLeapArray
		want   uint32
	}{
		{
			"test",
			NewBucketLeapArray(200, 2000),
			10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := &BucketLeapArray{
				data:     tt.fields.data,
				dataType: tt.fields.dataType,
			}
			if got := bla.BucketLengthInMs(); got != tt.want {
				t.Errorf("BucketLeapArray.BucketLengthInMs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketLeapArray_DataType(t *testing.T) {
	tests := []struct {
		name   string
		fields *BucketLeapArray
		want   string
	}{
		{
			"test",
			NewBucketLeapArray(200, 2000),
			"MetricBucket",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := &BucketLeapArray{
				data:     tt.fields.data,
				dataType: tt.fields.dataType,
			}
			if got := bla.DataType(); got != tt.want {
				t.Errorf("BucketLeapArray.DataType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketLeapArray_GetIntervalInSecond(t *testing.T) {
	tests := []struct {
		name   string
		fields *BucketLeapArray
		want   float64
	}{
		{
			"test",
			NewBucketLeapArray(2, 2000),
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := &BucketLeapArray{
				data:     tt.fields.data,
				dataType: tt.fields.dataType,
			}
			if got := bla.GetIntervalInSecond(); got != tt.want {
				t.Errorf("BucketLeapArray.GetIntervalInSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketLeapArray_newEmptyBucket(t *testing.T) {
	tests := []struct {
		name   string
		fields *BucketLeapArray
		want   interface{}
	}{
		{
			"test",
			NewBucketLeapArray(200, 2000),
			NewMetricBucket(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := tt.fields
			if got := bla.newEmptyBucket(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BucketLeapArray.newEmptyBucket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketLeapArray_resetBucketTo(t *testing.T) {
	wbl := NewBucketLeapArray(200, 2000)
	wbl.AddCount(base.MetricEventComplete, 1)
	bw := &bucketWrap{
		bucketStart: uint64(time.Now().Nanosecond()),
		value:       atomic.Value{},
	}
	type fields struct {
		data     leapArray
		dataType string
	}
	type args struct {
		ww        *bucketWrap
		startTime uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *bucketWrap
	}{
		{
			"test",
			fields{
				data:     wbl.data,
				dataType: "MetricBucket",
			},
			args{
				ww:        bw,
				startTime: uint64(time.Now().Nanosecond()),
			},
			bw,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := &BucketLeapArray{
				data:     tt.fields.data,
				dataType: tt.fields.dataType,
			}
			if got := bla.resetBucketTo(tt.args.ww, tt.args.startTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BucketLeapArray.resetBucketTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketLeapArray_AddCount(t *testing.T) {
	blw := NewBucketLeapArray(200, 2000)
	type fields struct {
		data     leapArray
		dataType string
	}
	type args struct {
		event base.MetricEvent
		count int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"test",
			fields{
				data:     blw.data,
				dataType: blw.dataType,
			},
			args{
				event: base.MetricEventComplete,
				count: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := &BucketLeapArray{
				data:     tt.fields.data,
				dataType: tt.fields.dataType,
			}
			bla.AddCount(tt.args.event, tt.args.count)
		})
	}
}

func TestBucketLeapArray_addCountWithTime(t *testing.T) {
	blw := NewBucketLeapArray(200, 2000)
	blw.AddCount(base.MetricEventComplete, 1)
	type fields struct {
		data     leapArray
		dataType string
	}
	type args struct {
		now   uint64
		event base.MetricEvent
		count int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"test",
			fields{
				data:     blw.data,
				dataType: blw.dataType,
			},
			args{
				now:   util.CurrentTimeMillis(),
				event: base.MetricEventComplete,
				count: 2,
			},
		},
		{
			"test2",
			fields{
				data:     blw.data,
				dataType: blw.dataType,
			},
			args{
				now:   util.CurrentTimeMillis() / 2,
				event: base.MetricEventComplete,
				count: 2,
			},
		},
		{
			"test3",
			fields{
				data:     blw.data,
				dataType: blw.dataType,
			},
			args{
				now:   util.CurrentTimeMillis(),
				event: base.MetricEventComplete,
				count: 2,
			},
		},
		{
			"test4",
			fields{
				data:     blw.data,
				dataType: blw.dataType,
			},
			args{
				now:   util.CurrentTimeMillis(),
				event: base.MetricEventComplete,
				count: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := &BucketLeapArray{
				data:     tt.fields.data,
				dataType: tt.fields.dataType,
			}
			if tt.name == "test3" {
				ad := &atomic.Value{}
				patch := gomonkey.ApplyMethod(reflect.TypeOf(ad), "Load", func(*atomic.Value) (x interface{}) {
					return nil
				})
				defer patch.Reset()
			}
			bla.addCountWithTime(tt.args.now, tt.args.event, tt.args.count)
		})
	}
}

func TestBucketLeapArray_Count(t *testing.T) {
	blw := NewBucketLeapArray(200, 2000)
	blw.AddCount(base.MetricEventComplete, 2)
	type fields struct {
		data     leapArray
		dataType string
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
				data:     blw.data,
				dataType: blw.dataType,
			},
			args{event: base.MetricEventComplete},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := &BucketLeapArray{
				data:     tt.fields.data,
				dataType: tt.fields.dataType,
			}
			if got := bla.Count(tt.args.event); got != tt.want {
				t.Errorf("BucketLeapArray.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketLeapArray_CountWithTime(t *testing.T) {
	blw := NewBucketLeapArray(10, 2000)
	blw.AddCount(base.MetricEventComplete, 4)
	type fields struct {
		data     leapArray
		dataType string
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
				data:     blw.data,
				dataType: blw.dataType,
			},
			args{
				now:   util.CurrentTimeMillis(),
				event: base.MetricEventComplete,
			},
			4,
		},
		{
			"test2",
			fields{
				data:     blw.data,
				dataType: blw.dataType,
			},
			args{
				now:   util.CurrentTimeMillis() / 2,
				event: base.MetricEventComplete,
			},
			0,
		},
		{
			"test3",
			fields{
				data:     blw.data,
				dataType: blw.dataType,
			},
			args{
				now:   util.CurrentTimeMillis(),
				event: base.MetricEventComplete,
			},
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := &BucketLeapArray{
				data:     tt.fields.data,
				dataType: tt.fields.dataType,
			}
			if tt.name == "test3" {
				time.Sleep(2000 * time.Millisecond)
			}
			if got := bla.CountWithTime(tt.args.now, tt.args.event); got != tt.want {
				t.Errorf("BucketLeapArray.CountWithTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketLeapArray_Values(t *testing.T) {
	blw := NewBucketLeapArray(10, 3000)
	type fields struct {
		data     leapArray
		dataType string
	}
	type args struct {
		now uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			"test",
			fields{
				data:     blw.data,
				dataType: blw.dataType,
			},
			args{now: util.CurrentTimeMillis()},
			10,
		},
		{
			"test2",
			fields{
				data:     blw.data,
				dataType: blw.dataType,
			},
			args{now: util.CurrentTimeMillis() / 2},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := &BucketLeapArray{
				data:     tt.fields.data,
				dataType: tt.fields.dataType,
			}
			if got := bla.Values(tt.args.now); len(got) != tt.want {
				t.Errorf("BucketLeapArray.Values() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketLeapArray_ValuesConditional(t *testing.T) {
	blw := NewBucketLeapArray(3, 1200)
	blw.AddCount(base.MetricEventComplete, 2)
	time.Sleep(400 * time.Millisecond)
	blw.AddCount(base.MetricEventComplete, 3)
	time.Sleep(400 * time.Millisecond)
	blw.AddCount(base.MetricEventComplete, 1)
	type fields struct {
		data     leapArray
		dataType string
	}
	type args struct {
		now       uint64
		predicate base.TimePredicate
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*bucketWrap
	}{
		{
			"test",
			fields{
				data:     blw.data,
				dataType: blw.dataType,
			},
			args{
				now: util.CurrentTimeMillis(),
				predicate: func(u uint64) bool {
					return false
				},
			},
			[]*bucketWrap{},
		},
		{
			"test2",
			fields{
				data:     blw.data,
				dataType: blw.dataType,
			},
			args{
				now: util.CurrentTimeMillis() / 2,
				predicate: func(u uint64) bool {
					return false
				},
			},
			[]*bucketWrap{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := &BucketLeapArray{
				data:     tt.fields.data,
				dataType: tt.fields.dataType,
			}
			if got := bla.ValuesConditional(tt.args.now, tt.args.predicate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BucketLeapArray.ValuesConditional() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketLeapArray_MinRt(t *testing.T) {
	blw := NewBucketLeapArray(10, 2000)
	blw.AddCount(base.MetricEventRt, 11)
	time.Sleep(30 * time.Millisecond)
	blw.AddCount(base.MetricEventRt, 3)
	time.Sleep(30 * time.Millisecond)
	blw.AddCount(base.MetricEventRt, 4)
	time.Sleep(30 * time.Millisecond)
	blw.AddCount(base.MetricEventRt, 20)
	time.Sleep(30 * time.Millisecond)
	blw.AddCount(base.MetricEventRt, 10)
	time.Sleep(30 * time.Millisecond)
	blw.AddCount(base.MetricEventRt, 40)
	time.Sleep(30 * time.Millisecond)
	type fields struct {
		data     leapArray
		dataType string
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			"test",
			fields{
				data:     blw.data,
				dataType: blw.dataType,
			},
			3,
		},
		{
			"test2",
			fields{
				data:     blw.data,
				dataType: blw.dataType,
			},
			60000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bla := &BucketLeapArray{
				data:     tt.fields.data,
				dataType: tt.fields.dataType,
			}

			if tt.name == "test2" {
				patch := gomonkey.ApplyFunc(util.CurrentTimeMillis, func() uint64 {
					return uint64(0)
				})
				defer patch.Reset()
			}

			if got := bla.MinRt(); got != tt.want {
				t.Errorf("BucketLeapArray.MinRt() = %v, want %v", got, tt.want)
			}
		})
	}
}
