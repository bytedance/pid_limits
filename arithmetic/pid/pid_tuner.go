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
	"log"
	"math"

	"github.com/bytedance/plato/util"
)

const (
	kpConstant = 0.2
	tiConstant = 2
	tdConstant = 0.333
	M_PI       = 3.14159265358979323846
)

type Tuner struct {
	microseconds                 uint64
	max                          float64
	min                          float64
	targetInputValue             float64
	minOutput                    float64
	maxOutput                    float64
	t1, t2                       uint64
	i                            int64
	cycles                       int64
	output                       bool
	outputValue                  float64
	tHigh, tLow                  uint64
	pAverage, iAverage, dAverage float64
	loopInterval                 int64
	kp, ki, kd                   float64
}

func (t *Tuner) Init(threshold float64) {
	t.cycles = 400  // Cycle counter
	t.output = true // Current output state
	t.outputValue = t.maxOutput
	t.t1, t.t2 = util.CurrentTimeMillis(), util.CurrentTimeMillis() // Times used for calculating period
	t.microseconds, t.tHigh, t.tLow = 0, 0, 0                       // More time variables
	t.max = 1                                                       // Max input
	t.min = 0                                                       // Min input
	t.pAverage, t.iAverage, t.dAverage = 0, 0, 0
	t.maxOutput = 0
	t.minOutput = -10000
	t.targetInputValue = threshold
	t.loopInterval = 100
}

func (t *Tuner) TunePID(input float64) float64 {

	if t.i > t.cycles {
		return 0
	}

	// Calculate time delta
	//prevMicroseconds := t.microseconds
	t.microseconds = util.CurrentTimeMillis()
	//deltaT := t.microseconds - prevMicroseconds;

	// Calculate max and min
	t.max = math.Max(t.max, input)
	t.min = math.Min(t.min, input)

	// Output is on and input signal has risen to target
	if t.output && input > t.targetInputValue {
		// Turn output off, record current time as t1, calculate tHigh, and reset maximum
		t.output = false
		t.outputValue = t.minOutput
		t.t1 = util.CurrentTimeMillis()
		t.tHigh = t.t1 - t.t2
		t.max = t.targetInputValue
	}

	// Output is off and input signal has dropped to target
	if !t.output && input < t.targetInputValue {
		// Turn output on, record current time as t2, calculate tLow
		t.output = true
		t.outputValue = t.maxOutput
		t.t2 = util.CurrentTimeMillis()
		t.tLow = t.t2 - t.t1

		// Calculate Ku (ultimate gain)
		// Formula given is Ku = 4d / πa
		// d is the amplitude of the output signal
		// a is the amplitude of the input signal
		ku := (4.0 * ((t.maxOutput - t.minOutput) / 2.0)) / (M_PI * (t.max - t.min) / 2.0)

		// Calculate Tu (period of output oscillations)
		tu := t.tLow + t.tHigh

		// Calculate gains
		t.kp = kpConstant * ku
		t.ki = (t.kp / (tiConstant * float64(tu))) * float64(t.loopInterval)
		t.kd = (tdConstant * t.kp * float64(tu)) / float64(t.loopInterval)

		// Average all gains after the first two cycles
		if t.i > 1 {
			t.pAverage += t.kp
			t.iAverage += t.ki
			t.dAverage += t.kd
		}

		// Reset minimum
		t.min = t.targetInputValue

		// Increment cycle count
		t.i += 1
		log.Printf("pidtuner: %v", t.i)
	}

	// If loop is done, disable output and calculate averages
	if t.i >= t.cycles {
		t.output = false
		t.outputValue = t.minOutput
		t.kp = t.pAverage / (float64(t.i) - 1)
		t.ki = t.iAverage / (float64(t.i) - 1)
		t.kd = t.dAverage / (float64(t.i) - 1)
		log.Printf("end = %v %v %v", t.kp, t.ki, t.kd)
	}

	return t.outputValue
}

func (t *Tuner) GetP() float64 {
	return t.kp
}

func (t *Tuner) GetI() float64 {
	return t.ki
}

func (t *Tuner) GetD() float64 {
	return t.kd
}
