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
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func Test_getCPURateByStat(t *testing.T) {
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
			if got, _ := getCPURateByStat(); got != tt.want {
				t.Errorf("getCPURateByStat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_recordCPUUsage(t *testing.T) {
	getCPURateByStat()
	time.Sleep(time.Second)
	cpuStats, err := cpu.Times(false)
	if err != nil {
		t.Errorf("cpu.Times error: %v\n", err)
		return
	}
	value, err := cpu.Percent(time.Second, false)
	if err != nil {
		t.Errorf("cpu.Percent error: %v\n", err)
		return
	}
	type args struct {
		prev       *cpu.TimesStat
		curCPUStat *cpu.TimesStat
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"test",
			args{
				prev:       prevCPUStat,
				curCPUStat: &cpuStats[0],
			},
			value[0]/100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := recordCPUUsage(tt.args.prev, tt.args.curCPUStat); got / tt.want < 0.1 || tt.want / got < 0.1 {
				t.Errorf("recordCPUUsage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateTotalCPUTick(t *testing.T) {
	type args struct {
		stat *cpu.TimesStat
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"test",
			args{stat: &cpu.TimesStat{
				CPU:       "cpu0",
				User:      1,
				System:    1,
				Idle:      1,
				Nice:      1,
				Iowait:    1,
				Irq:       1,
				Softirq:   1,
				Steal:     1,
				Guest:     1,
				GuestNice: 1,
			}},
			8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateTotalCPUTick(tt.args.stat); got != tt.want {
				t.Errorf("calculateTotalCPUTick() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateUserCPUTick(t *testing.T) {
	type args struct {
		stat *cpu.TimesStat
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"test",
			args{stat: &cpu.TimesStat{
				CPU:       "cpu0",
				User:      1,
				System:    1,
				Idle:      1,
				Nice:      1,
				Iowait:    1,
				Irq:       1,
				Softirq:   1,
				Steal:     1,
				Guest:     1,
				GuestNice: 1,
			}},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateUserCPUTick(tt.args.stat); got != tt.want {
				t.Errorf("calculateUserCPUTick() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateKernelCPUTick(t *testing.T) {
	type args struct {
		stat *cpu.TimesStat
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"test",
			args{stat: &cpu.TimesStat{
				CPU:       "cpu0",
				User:      1,
				System:    1,
				Idle:      1,
				Nice:      1,
				Iowait:    1,
				Irq:       1,
				Softirq:   1,
				Steal:     1,
				Guest:     1,
				GuestNice: 1,
			}},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateKernelCPUTick(tt.args.stat); got != tt.want {
				t.Errorf("calculateKernelCPUTick() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetByStat(t *testing.T) {
	for {
		f, e := getCPURateByStat()
		assert.Nil(t, e)
		fmt.Println(f)
		time.Sleep(time.Second)
	}
}