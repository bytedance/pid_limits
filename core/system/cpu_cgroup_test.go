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
	"log"
	"math"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shirou/gopsutil/cpu"
	"github.com/stretchr/testify/assert"
)

func Test_extractSystemPath(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test",
			args{line: "4:cpu,cpuacct:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6"},
			"/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractSystemPath(tt.args.line); got != tt.want {
				t.Errorf("extractSystemPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPodSystemPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		want string
	}{
		{
			"test",
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPodSystemPath(); false {
				t.Errorf("getPodSystemPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDockerSystemMetric(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			"test",
			args{filepath: "/test"},
			0,
		},
		{
			"test2",
			args{filepath: "./cpu_cgroup_test.out"},
			20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDockerSystemMetric(tt.args.filepath); false {
				t.Errorf("getDockerSystemMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initCGroup(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			"test",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initCGroup(); false {
				t.Errorf("initCGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileExist(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"test",
			args{filename: "./cpu_cgroup_test.go"},
			true,
		},
		{
			"test2",
			args{filename: "./cpu_cgroup_test_v2.go"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileExist(tt.args.filename); got != tt.want {
				t.Errorf("fileExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCPURateByCGroup(t *testing.T) {
	tests := []struct {
		name     string
		wantRate float64
	}{
		{
			"test",
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRate, _ := getCPURateByCGroup(); gotRate != tt.wantRate {
				t.Errorf("getCPURateByCGroup() = %v, want %v", gotRate, tt.wantRate)
			}
		})
	}
}

func TestExtractFilePath(t *testing.T) {
	var dataBox = []struct {
		input    string
		expected string
	}{{
		`10:memory:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
9:devices:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
8:blkio:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
7:pids:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
6:cpuset:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
5:freezer:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
4:cpu,cpuacct:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
3:perf_event:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
2:net_cls,net_prio:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
1:name=systemd:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6`,
		"/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6",
	}, {
		`4:cpu,cpuacct:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
10:memory:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
5:freezer:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
2:net_cls,net_prio:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
1:name=systemd:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6`,
		"/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6",
	}, {
		`10:memory:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
9:devices:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
8:blkio:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
7:pids:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
6:cpuset:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
5:freezer:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
3:perf_event:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
2:net_cls,net_prio:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
1:name=systemd:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6`,
		"",
	}, {
		`10:memory:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
9:devices:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
8:blkio:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
7:pids:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
6:cpuset:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
5:freezer:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
4:2:cpu,cpuacct:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
3:perf_event:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
2:net_cls,net_prio:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6
1:name=systemd:/kubepods/burstable/podf86cf783-f26a-11ea-bbcc-b8599fdf6096/0c1e3e53f1f37149bb95dd47c6b579ce5cf714af2d860f48d14c442f64e114f6`,
		"",
	}, {
		"",
		"",
	}}
	for _, tt := range dataBox {
		assert.Equal(t, tt.expected, extractSystemPath(tt.input))
	}
}

func TestGetCPURateByCGroup(t *testing.T) {
	if !initCGroup() {
		log.Println("can not use cgroup file in this machine")
		return
	}
	go func() {
		values, _ := cpu.Percent(0, false)
		cpuShirou := values[0]
		cpuCGroup, _ := getCPURateByCGroup()
		assert.Equal(t, true, math.Abs(cpuShirou-cpuCGroup) < 0.15)
		time.Sleep(100 * time.Millisecond)
	}()
	time.Sleep(10 * time.Second)
}

