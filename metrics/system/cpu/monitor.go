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

import (
	"log"

	"github.com/bytedance/plato/core/system"
)

const (
	cpuCollectorIntervalMs = 100
)

func init() {
	//init system metrics collector
	// as we consideration many times, we should keep to collect cpu usage every 100ms, for getting value timely
	system.InitCollector(cpuCollectorIntervalMs)
}

// GetUsage is used to get cpu current usage
func GetUsage() float64 {
	return system.CurrentCPUUsage()
}

type Monitor interface {
	IsOverload() bool
}

func NewCPUMonitor(ops ...Option) Monitor {
	opts := newOptions()
	for _, do := range ops {
		do.f(opts)
	}
	threshold := opts.upperThreshold()
	if threshold < 0 || threshold > 1 {
		log.Fatal("cpu usage threshold should be in 0 ~ 1")
	}

	switch opts.alg {
	case ZScore:
		return NewMonitorZScore(opts)
	case Raw:
		return NewMonitorRaw(opts)
	}
	return NewMonitorRaw(opts)
}
