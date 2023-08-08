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
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytedance/pid_limits/arithmetic/common"
	"github.com/bytedance/pid_limits/arithmetic/zscore"
	"github.com/bytedance/pid_limits/core/system"
	"github.com/bytedance/pid_limits/util"
)

/**
使用 zscore 判断CPU是否过载
*/

type MonitorZScore struct {
	upperThreshold func() float64
	lowerThreshold func() float64
	score          float64
	overload       atomic.Value
	initOnce       sync.Once
	continuousTime uint32 // 记录连续低于阈值的次数
}

const (
	// 连续多少次超过【低于】阈值，则开启【关闭】限流
	continuousTimes = uint32(30)
)

func NewMonitorZScore(opts *Options) Monitor {
	monitor := &MonitorZScore{
		score:          opts.score,
		upperThreshold: opts.upperThreshold,
		lowerThreshold: opts.lowerThreshold,
		overload:       atomic.Value{},
	}
	monitor.overload.Store(false)
	monitor.start()
	return monitor
}

// IsOverload method used to output whether the cpu is overload with the statistical method - zscore
func (monitor *MonitorZScore) IsOverload() bool {
	if monitor == nil {
		log.Println("error: adaptive cpu Monitor is nil")
		return false
	}
	if overload, ok := monitor.overload.Load().(bool); ok {
		return overload
	}
	log.Println("error: adaptive failed to get overload information from Monitor")
	return false
}

func (monitor *MonitorZScore) start() {
	monitor.initOnce.Do(func() {
		go util.LoopWithInterval(func() {
			monitor.decide()
		}, 100*time.Millisecond)
	})
}

func (monitor *MonitorZScore) decide() {
	if monitor == nil {
		log.Println("error: cpu Monitor is nil")
		return
	}
	rateWindows := system.ExtractCPUWindows()
	windows := zscore.ZScore(rateWindows, monitor.score)
	log.Printf("pid_debug: rateWindows = %+v\n", rateWindows)
	if len(windows) == 0 {
		return
	}
	avgCPU := common.AverageFloat(windows)
	log.Printf("pid_debug: avgCPU = %+v, monitor.IsOverlod() = %+v\n", avgCPU, monitor.IsOverload())
	log.Printf("pid_debug: upperThreshold = %+v, lowerThreshold = %+v\n", monitor.upperThreshold(), monitor.lowerThreshold())
	if monitor.IsOverload() {
		monitor.decideLessLoad(avgCPU, rateWindows, windows)
	} else {
		monitor.decideOverLoad(avgCPU, rateWindows, windows)
	}
}

func (monitor *MonitorZScore) decideOverLoad(avgCPU float64, rateWindows, windows []float64) {
	// 在没有过载的时候，需要连续30个计算周期【3秒】中，每次CPU平均值高于阈值上限
	if avgCPU >= monitor.upperThreshold() && atomic.AddUint32(&monitor.continuousTime, 1) > continuousTimes {
		log.Printf("warning: [adaptive limiting] start, rateWindows=%v windows=%v", rateWindows, windows)
		monitor.overload.Store(true)
		atomic.StoreUint32(&monitor.continuousTime, 0)
	}
	if avgCPU < monitor.lowerThreshold() {
		atomic.StoreUint32(&monitor.continuousTime, 0)
	}
}

func (monitor *MonitorZScore) decideLessLoad(avgCPU float64, rateWindows, windows []float64) {
	if avgCPU < monitor.lowerThreshold() && atomic.AddUint32(&monitor.continuousTime, 1) > continuousTimes {
		log.Printf("warning: [adaptive limiting] end, rateWindows=%v windows=%v", rateWindows, windows)
		monitor.overload.Store(false)
		atomic.StoreUint32(&monitor.continuousTime, 0)
	}
	if avgCPU >= monitor.upperThreshold() {
		// 如果有一次阈值高于上限，都需要将其清零
		atomic.StoreUint32(&monitor.continuousTime, 0)
	}
}
