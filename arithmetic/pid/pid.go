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

import "github.com/bytedance/plato/util"

const (
	OUTMAX = 0
	OUTMIN = -10000
)

type PID struct {
	kp, ki, kd  float64
	setPoint    float64
	getSetPoint func() float64
	errSum      float64
	lastErr     float64
	lastTime    uint64
	outMax      float64
	outMin      float64
}

type Option struct {
	GetSetPoint func() float64
}

type OptionFunc func(*Option)

func defaultOption() *Option {
	return &Option{GetSetPoint: nil}
}

func WithDynamicPoint(f func() float64) OptionFunc {
	return func(option *Option) {
		option.GetSetPoint = f
	}
}

func SetTunings(kp, ki, kd, setPoint float64, opts ...OptionFunc) *PID {
	option := defaultOption()
	for _, opt := range opts {
		opt(option)
	}
	pid := &PID{
		kp:       kp,
		ki:       ki,
		kd:       kd,
		setPoint: setPoint,
		getSetPoint: func() float64 {
			return setPoint
		},
		errSum:   0,
		lastTime: util.CurrentTimeMillis(),
		lastErr:  0,
		outMax:   OUTMAX,
		outMin:   OUTMIN,
	}
	if option.GetSetPoint != nil {
		pid.getSetPoint = option.GetSetPoint
	}
	return pid
}

func (pid *PID) GetThreshold() float64 {
	return pid.getSetPoint()
}

func (pid *PID) setOutLimit(max, min float64) {
	if max < min {
		return
	}
	if max < OUTMAX {
		pid.outMax = max
	}
	if min > OUTMIN {
		pid.outMin = min
	}
}

func (pid *PID) Compute(input float64) float64 {

	now := util.CurrentTimeMillis()
	timeChange := now - pid.lastTime
	err := pid.getSetPoint() - input
	old := pid.errSum
	pid.errSum = pid.errSum + err*(float64(timeChange))
	dErr := (err - pid.lastErr) / float64(timeChange)

	pid.lastErr = err
	pid.lastTime = now
	out := pid.kp*err + pid.ki*pid.errSum + pid.kd*dErr
	if out > pid.outMax {
		pid.errSum = old
		return pid.outMax
	} else if out < pid.outMin {
		pid.errSum = old
		return pid.outMin
	} else {
		return out
	}
}
