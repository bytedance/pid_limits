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
package flow

import (
	"github.com/agiledragon/gomonkey"
	"github.com/bytedance/pid_limits/core/base"
	stat "github.com/bytedance/pid_limits/core/stat/base"
	"github.com/bytedance/pid_limits/util"
	"math"
	"sync/atomic"
	"testing"
	"time"
)

func TestDoCheck(t *testing.T) {
	bucket := stat.NewBucketLeapArray(200, 2000)
	bucket.AddCount(base.MetricEventComplete, 1)
	type args struct {
		acquireCount uint32
		threshold    float64
		node         *stat.SlidingWindowMetric
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"test1", args{
				uint32(10),
				9.,
				&stat.SlidingWindowMetric{
					LastPassedTime: 1,
					Real:           bucket,
				},
			}, false,
		},
		{
			"test2", args{
				uint32(10),
				12.,
				&stat.SlidingWindowMetric{
					LastPassedTime: 1,
					Real:           bucket,
				},
			}, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DoCheck(tt.args.acquireCount, tt.args.threshold, tt.args.node); got != tt.want {
				t.Errorf("DoCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoCheck2(t *testing.T) {
	bucket := stat.NewBucketLeapArray(200, 2000)
	bucket.AddCount(base.MetricEventComplete, 1)
	type args struct {
		acquireCount uint32
		threshold    float64
		node         *stat.SlidingWindowMetric
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"test1", args{
				uint32(10),
				9.,
				&stat.SlidingWindowMetric{
					LastPassedTime: 1,
					Real:           bucket,
				},
			}, true,
		},
		{
			"test2", args{
				uint32(10),
				12.,
				&stat.SlidingWindowMetric{
					LastPassedTime: 1,
					Real:           bucket,
				},
			}, true,
		},
		{
			"test3", args{
				uint32(0),
				12.,
				&stat.SlidingWindowMetric{
					LastPassedTime: 1,
					Real:           bucket,
				},
			}, true,
		},
		{
			"test4", args{
				uint32(10),
				0.,
				&stat.SlidingWindowMetric{
					LastPassedTime: 1,
					Real:           bucket,
				},
			}, false,
		},
		{
			"test5", args{
				uint32(100000),
				3.,
				&stat.SlidingWindowMetric{
					LastPassedTime: util.CurrentTimeNano(),
					Real:           bucket,
				},
			}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DoCheck2(tt.args.acquireCount, tt.args.threshold, tt.args.node); got != tt.want {
				t.Errorf("DoCheck2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoCheck2V2(t *testing.T) {
	bucket := stat.NewBucketLeapArray(200, 2000)
	bucket.AddCount(base.MetricEventComplete, 1)
	type args struct {
		acquireCount uint32
		threshold    float64
		node         *stat.SlidingWindowMetric
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"test1", args{
				uint32(20),
				200000.,
				&stat.SlidingWindowMetric{
					LastPassedTime: util.CurrentTimeNano(),
					Real:           bucket,
				},
			}, true,
		},
	}
	for _, tt := range tests {
		exeTime := 1
		t.Run(tt.name, func(t *testing.T) {
			patch := gomonkey.ApplyFunc(util.CurrentTimeNano, func() uint64 {
				defer func() { exeTime++ }()
				if exeTime == 1 {
					return 0
				}
				if exeTime == 2 {
					return atomic.LoadUint64(&tt.args.node.LastPassedTime) + uint64(math.Ceil(float64(tt.args.acquireCount)/tt.args.threshold*float64(nanoUnitOffset)))
				}
				if exeTime == 3 {
					return tt.args.node.LastPassedTime
				}
				return 0
			})
			defer patch.Reset()
			if got := DoCheck2(tt.args.acquireCount, tt.args.threshold, tt.args.node); got != tt.want {
				t.Errorf("DoCheck2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoCheck2V3(t *testing.T) {
	bucket := stat.NewBucketLeapArray(200, 2000)
	bucket.AddCount(base.MetricEventComplete, 1)
	type args struct {
		acquireCount uint32
		threshold    float64
		node         *stat.SlidingWindowMetric
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"test1", args{
				uint32(20),
				200000.,
				&stat.SlidingWindowMetric{
					LastPassedTime: util.CurrentTimeNano(),
					Real:           bucket,
				},
			}, false,
		},
	}
	for _, tt := range tests {
		exeTime := 1
		t.Run(tt.name, func(t *testing.T) {
			patch := gomonkey.ApplyFunc(util.CurrentTimeNano, func() uint64 {
				defer func() { exeTime++ }()
				if exeTime == 1 {
					return 0
				}
				if exeTime == 2 {
					return atomic.LoadUint64(&tt.args.node.LastPassedTime) + uint64(math.Ceil(float64(tt.args.acquireCount)/tt.args.threshold*float64(nanoUnitOffset)))
				}
				if exeTime == 3 {
					return uint64(time.Now().UnixNano())
				}
				return 0
			})
			defer patch.Reset()
			if got := DoCheck2(tt.args.acquireCount, tt.args.threshold, tt.args.node); got != tt.want {
				t.Errorf("DoCheck2() = %v, want %v", got, tt.want)
			}
		})
	}
}
