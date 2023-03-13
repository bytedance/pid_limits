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
package  pid

import (
	"reflect"
	"testing"
)

func TestSetTunings(t *testing.T) {
	type args struct {
		kp       float64
		ki       float64
		kd       float64
		setPoint float64
	}
	tests := []struct {
		name string
		args args
		want *PID
	}{
		{
			"test",
			args{
				kp:       1,
				ki:       2,
				kd:       3,
				setPoint: 0.8,
			},
			&PID{
				kp:       1,
				ki:       2,
				kd:       3,
				setPoint: 0.8,
				errSum:   0,
				lastErr:  0,
				outMax:   OUTMAX,
				outMin:   OUTMIN,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetTunings(tt.args.kp, tt.args.ki, tt.args.kd, tt.args.setPoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetTunings() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestPID_GetThreshold(t *testing.T) {
	type fields struct {
		kp       float64
		ki       float64
		kd       float64
		setPoint float64
		errSum   float64
		lastErr  float64
		lastTime uint64
		outMax   float64
		outMin   float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			"test",
			fields{
				kp:       1,
				ki:       2,
				kd:       3,
				setPoint: 4,
				errSum:   5,
				lastErr:  6,
				lastTime: 7,
				outMax:   8,
				outMin:   9,
			},
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid := &PID{
				kp:       tt.fields.kp,
				ki:       tt.fields.ki,
				kd:       tt.fields.kd,
				setPoint: tt.fields.setPoint,
				errSum:   tt.fields.errSum,
				lastErr:  tt.fields.lastErr,
				lastTime: tt.fields.lastTime,
				outMax:   tt.fields.outMax,
				outMin:   tt.fields.outMin,
			}
			if got := pid.GetThreshold(); got != tt.want {
				t.Errorf("PID.GetThreshold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPID_setOutLimit(t *testing.T) {
	type fields struct {
		kp       float64
		ki       float64
		kd       float64
		setPoint float64
		errSum   float64
		lastErr  float64
		lastTime uint64
		outMax   float64
		outMin   float64
	}
	type args struct {
		max float64
		min float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"test",
			fields{
				kp:       1,
				ki:       2,
				kd:       3,
				setPoint: 4,
				errSum:   5,
				lastErr:  6,
				lastTime: 7,
				outMax:   8,
				outMin:   9,
			},
			args{
				max: 1,
				min: 2,
			},
		},
		{
			"test2",
			fields{
				kp:       1,
				ki:       2,
				kd:       3,
				setPoint: 4,
				errSum:   5,
				lastErr:  6,
				lastTime: 7,
				outMax:   8,
				outMin:   9,
			},
			args{
				max: -20,
				min: -40,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid := &PID{
				kp:       tt.fields.kp,
				ki:       tt.fields.ki,
				kd:       tt.fields.kd,
				setPoint: tt.fields.setPoint,
				errSum:   tt.fields.errSum,
				lastErr:  tt.fields.lastErr,
				lastTime: tt.fields.lastTime,
				outMax:   tt.fields.outMax,
				outMin:   tt.fields.outMin,
			}
			pid.setOutLimit(tt.args.max, tt.args.min)
		})
	}
}

func TestPID_Compute(t *testing.T) {
	type fields struct {
		kp       float64
		ki       float64
		kd       float64
		setPoint float64
		errSum   float64
		lastErr  float64
		lastTime uint64
		outMax   float64
		outMin   float64
	}
	type args struct {
		input float64
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
				kp:       1,
				ki:       2,
				kd:       3,
				setPoint: 4,
				errSum:   5,
				lastErr:  6,
				lastTime: 7,
				outMax:   8,
				outMin:   9,
			},
			args{input: 10},
			9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid := &PID{
				kp:       tt.fields.kp,
				ki:       tt.fields.ki,
				kd:       tt.fields.kd,
				setPoint: tt.fields.setPoint,
				errSum:   tt.fields.errSum,
				lastErr:  tt.fields.lastErr,
				lastTime: tt.fields.lastTime,
				outMax:   tt.fields.outMax,
				outMin:   tt.fields.outMin,
			}
			if got := pid.Compute(tt.args.input); got != tt.want {
				t.Errorf("PID.Compute() = %v, want %v", got, tt.want)
			}
		})
	}
}
