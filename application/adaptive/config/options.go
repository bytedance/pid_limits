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
package   config

import (
	"github.com/bytedance/plato/metrics/system/cpu"
)

type Options struct {
	EnableMetric        bool
	EnableOverloadScene bool
	MonitorAlg          cpu.MonitorAlg
	DynamicPoint        func() float64
	Drift               float64
}

type OptionFunc func(*Options)

func NewOptions() *Options {
	return &Options{
		EnableMetric:        true,
		EnableOverloadScene: false,
		MonitorAlg:          cpu.ZScore,
		DynamicPoint:        nil,
		Drift:               0.1,
	}
}

func WithDisableMetric() OptionFunc {
	return func(options *Options) {
		options.EnableMetric = false
	}
}

func WithEnableOverloadScene() OptionFunc {
	return func(options *Options) {
		options.EnableOverloadScene = true
	}
}

func WithMonitorAlg(alg cpu.MonitorAlg) OptionFunc {
	return func(options *Options) {
		options.MonitorAlg = alg
	}
}

func WithDynamicPoint(f func() float64) OptionFunc {
	return func(options *Options) {
		options.DynamicPoint = f
	}
}

func WithDrift(f float64) OptionFunc {
	return func(options *Options) {
		options.Drift = f
	}
}
