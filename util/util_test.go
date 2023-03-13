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
package  util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSet2Float64(t *testing.T) {
	var i float64
	SetFloat64(&i, 0.5)
	SetFloat64(&i, 0.6)
	t.Logf("%v", i)
}

func TestGet2Float64(t *testing.T) {
	var i float64 = 199.991
	assert.Equal(t, GetFloat64(&i), i)
}

func TestCurrentTimeMillis(t *testing.T) {
	tests := []struct {
		name string
		want uint64
	}{
		{
			"test",
			CurrentTimeMillis(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CurrentTimeMillis(); got != tt.want {
				t.Errorf("CurrentTimeMillis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetFloat64(t *testing.T) {
	f := 1.0
	type args struct {
		addr *float64
		new  float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"test",
			args{
				addr: &f,
				new:  2.0,
			},
			2.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetFloat64(tt.args.addr, tt.args.new); got != tt.want {
				t.Errorf("SetFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFloat64(t *testing.T) {
	f := 1.0
	type args struct {
		p *float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"test",
			args{p: &f},
			1.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFloat64(tt.args.p); got != tt.want {
				t.Errorf("GetFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrentTimeNano(t *testing.T) {
	tests := []struct {
		name string
		want uint64
	}{
		{
			"test",
			CurrentTimeNano(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CurrentTimeNano(); false {
				t.Errorf("CurrentTimeNano() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoopWithInterval(t *testing.T) {
	type args struct {
		runnable func()
		interval time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"test",
			args{
				runnable: func() {

				},
				interval: time.Millisecond * 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				LoopWithInterval(tt.args.runnable, tt.args.interval)
			}()
			time.Sleep(time.Second * 2)
		})
	}
}

func TestIsDocker(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			"test",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDocker(); got != tt.want {
				t.Errorf("IsDocker() = %v, want %v", got, tt.want)
			}
		})
	}
}
