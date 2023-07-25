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
package cpu

import (
	"log"
	"time"

	"github.com/bytedance/pid_limits/metrics/system/cpu"
)

// monitor should be global variable
var monitor cpu.Monitor

func cpuMonitorDemo() {
	// default cpu overload threshold is 0.9
	monitor = cpu.NewCPUMonitor(cpu.WithUpperBound(func() float64 {
		return 0.9
	}))
	for {
		log.Printf("CPU usage: %v, CPU overload: %v\n", cpu.GetUsage(), monitor.IsOverload())
		time.Sleep(time.Second)
	}
}
