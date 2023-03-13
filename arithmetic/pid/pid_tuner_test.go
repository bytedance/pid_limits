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
	"github.com/bytedance/plato/util"
	"testing"
)

func TestTuner_Init(t *testing.T) {
	type fields struct {
		microseconds     uint64
		max              float64
		min              float64
		targetInputValue float64
		minOutput        float64
		maxOutput        float64
		t1               uint64
		t2               uint64
		i                int64
		cycles           int64
		output           bool
		outputValue      float64
		tHigh            uint64
		tLow             uint64
		pAverage         float64
		iAverage         float64
		dAverage         float64
		loopInterval     int64
		kp               float64
		ki               float64
		kd               float64
	}
	type args struct {
		threshold float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"test",
			fields{
				microseconds:     0,
				max:              1,
				min:              0,
				targetInputValue: 0,
				minOutput:        0,
				maxOutput:        0,
				t1:               util.CurrentTimeMillis(),
				t2:               util.CurrentTimeMillis(),
				i:                0,
				cycles:           400,
				output:           true,
				outputValue:      0,
				tHigh:            0,
				tLow:             0,
				pAverage:         0,
				iAverage:         0,
				dAverage:         0,
				loopInterval:     0,
				kp:               0,
				ki:               0,
				kd:               0,
			},
			args{threshold: 0.8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tuner{
				microseconds:     tt.fields.microseconds,
				max:              tt.fields.max,
				min:              tt.fields.min,
				targetInputValue: tt.fields.targetInputValue,
				minOutput:        tt.fields.minOutput,
				maxOutput:        tt.fields.maxOutput,
				t1:               tt.fields.t1,
				t2:               tt.fields.t2,
				i:                tt.fields.i,
				cycles:           tt.fields.cycles,
				output:           tt.fields.output,
				outputValue:      tt.fields.outputValue,
				tHigh:            tt.fields.tHigh,
				tLow:             tt.fields.tLow,
				pAverage:         tt.fields.pAverage,
				iAverage:         tt.fields.iAverage,
				dAverage:         tt.fields.dAverage,
				loopInterval:     tt.fields.loopInterval,
				kp:               tt.fields.kp,
				ki:               tt.fields.ki,
				kd:               tt.fields.kd,
			}
			tr.Init(tt.args.threshold)
		})
	}
}

func TestTuner_TunePID(t *testing.T) {
	type fields struct {
		microseconds     uint64
		max              float64
		min              float64
		targetInputValue float64
		minOutput        float64
		maxOutput        float64
		t1               uint64
		t2               uint64
		i                int64
		cycles           int64
		output           bool
		outputValue      float64
		tHigh            uint64
		tLow             uint64
		pAverage         float64
		iAverage         float64
		dAverage         float64
		loopInterval     int64
		kp               float64
		ki               float64
		kd               float64
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
				microseconds:     0,
				max:              1,
				min:              0,
				targetInputValue: 0,
				minOutput:        0,
				maxOutput:        0,
				t1:               util.CurrentTimeMillis(),
				t2:               util.CurrentTimeMillis(),
				i:                0,
				cycles:           400,
				output:           true,
				outputValue:      0,
				tHigh:            0,
				tLow:             0,
				pAverage:         0,
				iAverage:         0,
				dAverage:         0,
				loopInterval:     0,
				kp:               0,
				ki:               0,
				kd:               0,
			},
			args{input: 10},
			0,
		},
		{
			"test2",
			fields{
				microseconds:     0,
				max:              1,
				min:              0,
				targetInputValue: 0,
				minOutput:        0,
				maxOutput:        0,
				t1:               util.CurrentTimeMillis(),
				t2:               util.CurrentTimeMillis(),
				i:                100,
				cycles:           20,
				output:           true,
				outputValue:      0,
				tHigh:            0,
				tLow:             0,
				pAverage:         0,
				iAverage:         0,
				dAverage:         0,
				loopInterval:     0,
				kp:               0,
				ki:               0,
				kd:               0,
			},
			args{input: 10},
			0,
		},
		{
			"test",
			fields{
				microseconds:     0,
				max:              1,
				min:              0,
				targetInputValue: 20,
				minOutput:        0,
				maxOutput:        0,
				t1:               util.CurrentTimeMillis(),
				t2:               util.CurrentTimeMillis(),
				i:                0,
				cycles:           400,
				output:           false,
				outputValue:      0,
				tHigh:            0,
				tLow:             0,
				pAverage:         0,
				iAverage:         0,
				dAverage:         0,
				loopInterval:     0,
				kp:               0,
				ki:               0,
				kd:               0,
			},
			args{input: 1},
			0,
		},
		{
			"test",
			fields{
				microseconds:     0,
				max:              1,
				min:              0,
				targetInputValue: 20,
				minOutput:        0,
				maxOutput:        0,
				t1:               util.CurrentTimeMillis(),
				t2:               util.CurrentTimeMillis(),
				i:                2,
				cycles:           400,
				output:           false,
				outputValue:      0,
				tHigh:            0,
				tLow:             0,
				pAverage:         0,
				iAverage:         0,
				dAverage:         0,
				loopInterval:     0,
				kp:               0,
				ki:               0,
				kd:               0,
			},
			args{input: 1},
			0,
		},
		{
			"test",
			fields{
				microseconds:     0,
				max:              1,
				min:              0,
				targetInputValue: 20,
				minOutput:        0,
				maxOutput:        0,
				t1:               util.CurrentTimeMillis(),
				t2:               util.CurrentTimeMillis(),
				i:                400,
				cycles:           400,
				output:           false,
				outputValue:      0,
				tHigh:            0,
				tLow:             0,
				pAverage:         0,
				iAverage:         0,
				dAverage:         0,
				loopInterval:     0,
				kp:               0,
				ki:               0,
				kd:               0,
			},
			args{input: 1},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tuner{
				microseconds:     tt.fields.microseconds,
				max:              tt.fields.max,
				min:              tt.fields.min,
				targetInputValue: tt.fields.targetInputValue,
				minOutput:        tt.fields.minOutput,
				maxOutput:        tt.fields.maxOutput,
				t1:               tt.fields.t1,
				t2:               tt.fields.t2,
				i:                tt.fields.i,
				cycles:           tt.fields.cycles,
				output:           tt.fields.output,
				outputValue:      tt.fields.outputValue,
				tHigh:            tt.fields.tHigh,
				tLow:             tt.fields.tLow,
				pAverage:         tt.fields.pAverage,
				iAverage:         tt.fields.iAverage,
				dAverage:         tt.fields.dAverage,
				loopInterval:     tt.fields.loopInterval,
				kp:               tt.fields.kp,
				ki:               tt.fields.ki,
				kd:               tt.fields.kd,
			}
			if got := tr.TunePID(tt.args.input); got != tt.want {
				t.Errorf("Tuner.TunePID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuner_GetP(t *testing.T) {
	type fields struct {
		microseconds     uint64
		max              float64
		min              float64
		targetInputValue float64
		minOutput        float64
		maxOutput        float64
		t1               uint64
		t2               uint64
		i                int64
		cycles           int64
		output           bool
		outputValue      float64
		tHigh            uint64
		tLow             uint64
		pAverage         float64
		iAverage         float64
		dAverage         float64
		loopInterval     int64
		kp               float64
		ki               float64
		kd               float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			"test",
			fields{
				microseconds:     0,
				max:              1,
				min:              0,
				targetInputValue: 0,
				minOutput:        0,
				maxOutput:        0,
				t1:               util.CurrentTimeMillis(),
				t2:               util.CurrentTimeMillis(),
				i:                0,
				cycles:           400,
				output:           true,
				outputValue:      0,
				tHigh:            0,
				tLow:             0,
				pAverage:         0,
				iAverage:         0,
				dAverage:         0,
				loopInterval:     0,
				kp:               0,
				ki:               0,
				kd:               0,
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tuner{
				microseconds:     tt.fields.microseconds,
				max:              tt.fields.max,
				min:              tt.fields.min,
				targetInputValue: tt.fields.targetInputValue,
				minOutput:        tt.fields.minOutput,
				maxOutput:        tt.fields.maxOutput,
				t1:               tt.fields.t1,
				t2:               tt.fields.t2,
				i:                tt.fields.i,
				cycles:           tt.fields.cycles,
				output:           tt.fields.output,
				outputValue:      tt.fields.outputValue,
				tHigh:            tt.fields.tHigh,
				tLow:             tt.fields.tLow,
				pAverage:         tt.fields.pAverage,
				iAverage:         tt.fields.iAverage,
				dAverage:         tt.fields.dAverage,
				loopInterval:     tt.fields.loopInterval,
				kp:               tt.fields.kp,
				ki:               tt.fields.ki,
				kd:               tt.fields.kd,
			}
			if got := tr.GetP(); got != tt.want {
				t.Errorf("Tuner.GetP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuner_GetI(t *testing.T) {
	type fields struct {
		microseconds     uint64
		max              float64
		min              float64
		targetInputValue float64
		minOutput        float64
		maxOutput        float64
		t1               uint64
		t2               uint64
		i                int64
		cycles           int64
		output           bool
		outputValue      float64
		tHigh            uint64
		tLow             uint64
		pAverage         float64
		iAverage         float64
		dAverage         float64
		loopInterval     int64
		kp               float64
		ki               float64
		kd               float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			"test",
			fields{
				microseconds:     0,
				max:              1,
				min:              0,
				targetInputValue: 0,
				minOutput:        0,
				maxOutput:        0,
				t1:               util.CurrentTimeMillis(),
				t2:               util.CurrentTimeMillis(),
				i:                0,
				cycles:           400,
				output:           true,
				outputValue:      0,
				tHigh:            0,
				tLow:             0,
				pAverage:         0,
				iAverage:         0,
				dAverage:         0,
				loopInterval:     0,
				kp:               0,
				ki:               0,
				kd:               0,
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tuner{
				microseconds:     tt.fields.microseconds,
				max:              tt.fields.max,
				min:              tt.fields.min,
				targetInputValue: tt.fields.targetInputValue,
				minOutput:        tt.fields.minOutput,
				maxOutput:        tt.fields.maxOutput,
				t1:               tt.fields.t1,
				t2:               tt.fields.t2,
				i:                tt.fields.i,
				cycles:           tt.fields.cycles,
				output:           tt.fields.output,
				outputValue:      tt.fields.outputValue,
				tHigh:            tt.fields.tHigh,
				tLow:             tt.fields.tLow,
				pAverage:         tt.fields.pAverage,
				iAverage:         tt.fields.iAverage,
				dAverage:         tt.fields.dAverage,
				loopInterval:     tt.fields.loopInterval,
				kp:               tt.fields.kp,
				ki:               tt.fields.ki,
				kd:               tt.fields.kd,
			}
			if got := tr.GetI(); got != tt.want {
				t.Errorf("Tuner.GetI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuner_GetD(t *testing.T) {
	type fields struct {
		microseconds     uint64
		max              float64
		min              float64
		targetInputValue float64
		minOutput        float64
		maxOutput        float64
		t1               uint64
		t2               uint64
		i                int64
		cycles           int64
		output           bool
		outputValue      float64
		tHigh            uint64
		tLow             uint64
		pAverage         float64
		iAverage         float64
		dAverage         float64
		loopInterval     int64
		kp               float64
		ki               float64
		kd               float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			"test",
			fields{
				microseconds:     0,
				max:              1,
				min:              0,
				targetInputValue: 3,
				minOutput:        0,
				maxOutput:        0,
				t1:               util.CurrentTimeMillis(),
				t2:               util.CurrentTimeMillis(),
				i:                0,
				cycles:           400,
				output:           false,
				outputValue:      0,
				tHigh:            0,
				tLow:             0,
				pAverage:         0,
				iAverage:         0,
				dAverage:         0,
				loopInterval:     0,
				kp:               0,
				ki:               0,
				kd:               0,
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tuner{
				microseconds:     tt.fields.microseconds,
				max:              tt.fields.max,
				min:              tt.fields.min,
				targetInputValue: tt.fields.targetInputValue,
				minOutput:        tt.fields.minOutput,
				maxOutput:        tt.fields.maxOutput,
				t1:               tt.fields.t1,
				t2:               tt.fields.t2,
				i:                tt.fields.i,
				cycles:           tt.fields.cycles,
				output:           tt.fields.output,
				outputValue:      tt.fields.outputValue,
				tHigh:            tt.fields.tHigh,
				tLow:             tt.fields.tLow,
				pAverage:         tt.fields.pAverage,
				iAverage:         tt.fields.iAverage,
				dAverage:         tt.fields.dAverage,
				loopInterval:     tt.fields.loopInterval,
				kp:               tt.fields.kp,
				ki:               tt.fields.ki,
				kd:               tt.fields.kd,
			}
			if got := tr.GetD(); got != tt.want {
				t.Errorf("Tuner.GetD() = %v, want %v", got, tt.want)
			}
		})
	}
}
