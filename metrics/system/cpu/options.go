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
package  cpu

const (
	score = 2.4
	threshold = 0.9
	lowerThreshold = threshold * 0.9
	alg = Raw
)

type MonitorAlg int

const (
	ZScore MonitorAlg = iota
	Raw
)

// Option .
type Option struct {
	f func(*Options)
}

type Options struct {
	upperThreshold      func()float64
	lowerThreshold		func()float64
	alg                 MonitorAlg
	score               float64
}

func newOptions() *Options {
	return &Options{
		upperThreshold: func() float64 {
			return threshold
		},
		lowerThreshold: func() float64 {
			return lowerThreshold
		},
		alg:                 alg,
		score:               score,
	}
}

// WithThresholdScore is used to set z-score critical value
func WithThresholdScore(score float64) Option {
	return Option{f: func(options *Options) {
		options.score = score
	}}
}

// WithUpperBound is used to set threshold
func WithUpperBound(value func()float64) Option {
	return Option{f: func(options *Options) {
		options.upperThreshold = value
	}}
}

// WithLowerBound is used to set low bound of threshold
func WithLowerBound(value func()float64) Option {
	return Option{f: func(options *Options) {
		options.lowerThreshold = value
	}}
}

// WithAlg is used to set algorithm for the overload decision
func WithAlg(alg MonitorAlg) Option {
	return Option{f: func(options *Options) {
		options.alg = alg
	}}
}