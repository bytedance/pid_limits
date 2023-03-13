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
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytedance/plato/util"
)

const (
	waitTime = 6 * time.Second
)

type MonitorRaw struct {
	upperThreshold float64
	lowerThreshold float64
	overload       atomic.Value
	initOnce       sync.Once
}

func NewMonitorRaw(opts *Options) Monitor {
	monitor := &MonitorRaw{
		upperThreshold: opts.upperThreshold(),
		lowerThreshold: opts.lowerThreshold(),
		overload:       atomic.Value{},
		initOnce:       sync.Once{},
	}
	monitor.overload.Store(false)
	monitor.start()
	return monitor
}

func (monitor *MonitorRaw) IsOverload() bool {
	if monitor == nil {
		log.Println("error: adaptive cpu Monitor is nil")
		return false
	}
	if overload, ok := monitor.overload.Load().(bool); ok {
		return overload
	}
	log.Println("adaptive failed to get overload information from Monitor")
	return false
}

func (monitor *MonitorRaw) start() {
	monitor.initOnce.Do(func() {
		go util.LoopWithInterval(func() {
			usage := GetUsage()
			if usage >= monitor.upperThreshold && !monitor.IsOverload() {
				time.Sleep(waitTime)
				if GetUsage() >= monitor.upperThreshold {
					monitor.overload.Store(true)
				}
				return
			}
			if usage < monitor.lowerThreshold && monitor.IsOverload() {
				time.Sleep(waitTime)
				if GetUsage() < monitor.lowerThreshold {
					monitor.overload.Store(false)
				}
				return
			}
		}, 100*time.Millisecond)
	})
}
