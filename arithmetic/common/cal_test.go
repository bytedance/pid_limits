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
package  common

import "testing"

func TestSumFloat(t *testing.T) {
	type args struct {
		x []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"test",
			args{x: []float64{1,2,3,4}},
			10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumFloat(tt.args.x); got != tt.want {
				t.Errorf("SumFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAverageFloat(t *testing.T) {
	type args struct {
		x []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"test",
			args{x: []float64{1,2,4,10,1, 102}},
			20,
		},
		{
			"test2",
			args{x: []float64{}},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AverageFloat(tt.args.x); got != tt.want {
				t.Errorf("AverageFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStandardDeviation(t *testing.T) {
	type args struct {
		x   []float64
		avg float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"test",
			args{
				x:   []float64{1,2,4,10,1, 102},
				avg: 20,
			},
			36.80126809409335,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StandardDeviation(tt.args.x, tt.args.avg); got != tt.want {
				t.Errorf("StandardDeviation() = %v, want %v", got, tt.want)
			}
		})
	}
}
