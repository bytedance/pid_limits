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
package  cpu

import (
	"reflect"
	"testing"
)

func Test_newOptions(t *testing.T) {
	tests := []struct {
		name string
		want *Options
	}{
		{
			"test",
			&Options{
				thresholdValue: 0.9,
				score:          2.2,
				thresholdLess:  12,
				thresholdOver:  4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithThresholdValue(t *testing.T) {
	type args struct {
		threshold float64
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"test",
			args{threshold: 0.8},
			Option{f: func(options *Options) {
				options.thresholdValue = 0.8
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithThresholdValue(tt.args.threshold); NewCPUMonitor(got).threshold.value != NewCPUMonitor(tt.want).threshold.value {
				t.Errorf("WithThresholdValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithThresholdScore(t *testing.T) {
	type args struct {
		score float64
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"test",
			args{score: 2.1},
			Option{f: func(options *Options) {
				options.score = 2.1
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithThresholdScore(tt.args.score); NewCPUMonitor(got).threshold.score != NewCPUMonitor(tt.want).threshold.score {
				t.Errorf("WithThresholdScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithThresholdLess(t *testing.T) {
	type args struct {
		less int64
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"test",
			args{less: 2},
			Option{f: func(options *Options) {
				options.thresholdLess = 2
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithThresholdLess(tt.args.less); NewCPUMonitor(got).threshold.less != NewCPUMonitor(tt.want).threshold.less {
				t.Errorf("WithThresholdLess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithThresholdOver(t *testing.T) {
	type args struct {
		over int64
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"test",
			args{over: 2},
			Option{f: func(options *Options) {
				options.thresholdOver = 2
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithThresholdOver(tt.args.over); NewCPUMonitor(got).threshold.over != NewCPUMonitor(tt.want).threshold.over {
				t.Errorf("WithThresholdOver() = %v, want %v", got, tt.want)
			}
		})
	}
}
