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

	"github.com/shirou/gopsutil/cpu"
)

/*
Get cpu usage rate by cgroup
*/
var (
	prevCPUStat *cpu.TimesStat
)

func getCPURateByStat() (float64, error) {
	var cpuRate float64
	cpuStats, err := cpu.Times(false)
	if err != nil {
		log.Printf("warning: Failed to retrieve current CPU usage: %+v", err)
	}
	if len(cpuStats) > 0 {
		curCPUStat := &cpuStats[0]
		cpuRate, err = recordCPUUsage(prevCPUStat, curCPUStat)
		if err == retrieveValueError {
			return retrieveValueFailed, err
		}
		// Cache the latest CPU stat info.
		prevCPUStat = curCPUStat
	} else {
		log.Println("error: pid_error: len(cpuStats) == 0")
	}
	return cpuRate, nil
}

func recordCPUUsage(prev, curCPUStat *cpu.TimesStat) (float64, error) {
	var cpuUsage = float64(0)
	if prev != nil && curCPUStat != nil {
		prevTotal := calculateTotalCPUTick(prev)
		curTotal := calculateTotalCPUTick(curCPUStat)
		tDiff := curTotal - prevTotal
		// if tDiff == 0, then do not change cpu usage, keep the old value should be best
		if tDiff != 0 {
			prevUsed := calculateUserCPUTick(prev) + calculateKernelCPUTick(prev)
			curUsed := calculateUserCPUTick(curCPUStat) + calculateKernelCPUTick(curCPUStat)
			cpuUsage = (curUsed - prevUsed) / tDiff
			cpuUsage = math.Max(0.0, cpuUsage)
		} else {
			return retrieveValueFailed, retrieveValueError
		}
	} else {
		return retrieveValueFailed, errPrevStatsNil
	}
	return cpuUsage, nil
}

func calculateTotalCPUTick(stat *cpu.TimesStat) float64 {
	return stat.User + stat.Nice + stat.System + stat.Idle +
		stat.Iowait + stat.Irq + stat.Softirq + stat.Steal
}

func calculateUserCPUTick(stat *cpu.TimesStat) float64 {
	return stat.User + stat.Nice
}

func calculateKernelCPUTick(stat *cpu.TimesStat) float64 {
	return stat.System + stat.Irq + stat.Softirq
}
