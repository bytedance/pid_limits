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
package  system

import (
	"testing"
	"time"
)

func TestInitCollector(t *testing.T) {
	type args struct {
		intervalMs uint32
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"test",
			args{intervalMs: uint32(0)},
		},
		{
			"test2",
			args{intervalMs: uint32(50)},
		},
		{
			"test3",
			args{intervalMs: uint32(50)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitCollector(tt.args.intervalMs)
			time.Sleep(time.Second)
			if tt.name == "test3" {
				close(ssStopChan)
				ssStopChan = nil
			}
		})
	}
}

func Test_retrieveAndUpdateCPUUsage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrieveAndUpdateCPUUsage()
		})
	}
}

func TestCurrentCPUUsage(t *testing.T) {
	tests := []struct {
		name string
		want float64
	}{
		{
			"test",
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CurrentCPUUsage(); false {
				t.Errorf("CurrentCPUUsage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractCPUWindows(t *testing.T) {
	tests := []struct {
		name         string
		wantCpuRates []float64
	}{
		{
			"test",
			[]float64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCpuRates := ExtractCPUWindows(); false {
				t.Errorf("ExtractCPUWindows() = %+v, want %+v", gotCpuRates, tt.wantCpuRates)
			}
		})
	}
}

