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
	"reflect"
	"testing"

	"github.com/bytedance/plato/core/base"
)

func TestNewMetricBucket(t *testing.T) {
	tests := []struct {
		name string
		want *MetricBucket
	}{
		{
			"test",
			&MetricBucket{
				minRt: base.DefaultStatisticMaxRt,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMetricBucket(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricBucket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricBucket_Add(t *testing.T) {
	mb := NewMetricBucket()
	type fields struct {
		counter [base.MetricEventTotal]int64
		minRt   int64
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
				counter: mb.counter,
				minRt:   mb.minRt,
			},
			args{
				event: base.MetricEventComplete,
				count: 2,
			},
		},
		{
			"test2",
			fields{
				counter: mb.counter,
				minRt:   mb.minRt,
			},
			args{
				event: base.MetricEventRt,
				count: 2,
			},
		},
		{
			"test3",
			fields{
				counter: mb.counter,
				minRt:   mb.minRt,
			},
			args{
				event: -1,
				count: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := &MetricBucket{
				counter: tt.fields.counter,
				minRt:   tt.fields.minRt,
			}
			defer func() {
				_ = recover()
			}()
			mb.Add(tt.args.event, tt.args.count)
		})
	}
}

func TestMetricBucket_addCount(t *testing.T) {
	mb := NewMetricBucket()
	type fields struct {
		counter [base.MetricEventTotal]int64
		minRt   int64
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
				counter: mb.counter,
				minRt:   mb.minRt,
			},
			args{
				event: base.MetricEventComplete,
				count: 2,
			},
		},
		{
			"test2",
			fields{
				counter: mb.counter,
				minRt:   mb.minRt,
			},
			args{
				event: base.MetricEventRt,
				count: 2,
			},
		},
		{
			"test3",
			fields{
				counter: mb.counter,
				minRt:   mb.minRt,
			},
			args{
				event: -1,
				count: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := &MetricBucket{
				counter: tt.fields.counter,
				minRt:   tt.fields.minRt,
			}
			defer func() {
				_ = recover()
			}()
			mb.addCount(tt.args.event, tt.args.count)
		})
	}
}

func TestMetricBucket_Get(t *testing.T) {
	mb := NewMetricBucket()
	mb.Add(base.MetricEventComplete, 2)
	type fields struct {
		counter [base.MetricEventTotal]int64
		minRt   int64
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
				counter: mb.counter,
				minRt:   mb.minRt,
			},
			args{event: base.MetricEventComplete},
			2,
		},
		{
			"test",
			fields{
				counter: mb.counter,
				minRt:   mb.minRt,
			},
			args{event: -1},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := &MetricBucket{
				counter: tt.fields.counter,
				minRt:   tt.fields.minRt,
			}
			defer func() {
				_ = recover()
			}()
			if got := mb.Get(tt.args.event); got != tt.want {
				t.Errorf("MetricBucket.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricBucket_reset(t *testing.T) {
	mb := NewMetricBucket()
	type fields struct {
		counter [base.MetricEventTotal]int64
		minRt   int64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"test",
			fields{
				counter: mb.counter,
				minRt:   mb.minRt,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := &MetricBucket{
				counter: tt.fields.counter,
				minRt:   tt.fields.minRt,
			}
			mb.reset()
		})
	}
}

func TestMetricBucket_AddRt(t *testing.T) {
	mb := NewMetricBucket()
	type fields struct {
		counter [base.MetricEventTotal]int64
		minRt   int64
	}
	type args struct {
		rt int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"test",
			fields{
				counter: mb.counter,
				minRt:   mb.minRt,
			},
			args{rt: 20},
		},
		{
			"test2",
			fields{
				counter: mb.counter,
				minRt:   mb.minRt,
			},
			args{rt: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := &MetricBucket{
				counter: tt.fields.counter,
				minRt:   tt.fields.minRt,
			}
			mb.AddRt(tt.args.rt)
		})
	}
}

func TestMetricBucket_MinRt(t *testing.T) {
	mb := NewMetricBucket()
	mb.AddRt(9)
	mb.AddRt(7)
	mb.AddRt(10)
	type fields struct {
		counter [base.MetricEventTotal]int64
		minRt   int64
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			"test",
			fields{
				counter: mb.counter,
				minRt:   mb.minRt,
			},
			7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := &MetricBucket{
				counter: tt.fields.counter,
				minRt:   tt.fields.minRt,
			}
			if got := mb.MinRt(); got != tt.want {
				t.Errorf("MetricBucket.MinRt() = %v, want %v", got, tt.want)
			}
		})
	}
}
