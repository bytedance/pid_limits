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
package system

import (
	"errors"
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytedance/pid_limits/core/stat"
	"golang.org/x/sys/unix"
)

const (
	notRetrievedValue   float64 = 0
	retrieveValueFailed float64 = -1
	scale               float64 = 10000000
)

var (
	currentCPUUsage    atomic.Value
	initOnce           sync.Once
	ssStopChan         = make(chan struct{})
	slidingWindow      = stat.NewSlidingWindow(100, 6*time.Second)
	getCPURate         = getCPURateByStat
	retrieveValueError = errors.New("can not retrieveValue from ")
	errPrevStatsNil    = errors.New("PREV STAT IS NIL")
	disablePIDLimit    = false
)

func init() {
	// judge the file related with CGroup exits of not
	currentCPUUsage.Store(notRetrievedValue)
	if isCgroup2UnifiedMode() && initCGroupV2() {
		log.Println("current version is cgroupv2")
		getCPURate = getCPURateByCGroupV2
		return
	}
	if initCGroup() {
		log.Println("current version is cgroupv1")
		getCPURate = getCPURateByCGroup
		return
	}
	disablePIDLimit = true
}

func isCgroup2UnifiedMode() bool {
	var st unix.Statfs_t
	err := unix.Statfs(cGroupV2prefixPath, &st)
	if err != nil {
		return false
	}
	// 	原本为 st.Type == unix.CGROUP2_SUPER_MAGIC， 但是 mac 环境没有 CGROUP2_SUPER_MAGIC 变量，影响本地编译测试，故使用原本的魔术
	return st.Type == 0x63677270
}

func InitCollector(intervalMs uint32) {
	if intervalMs == 0 {
		return
	}
	initOnce.Do(func() {
		// Initial retrieval.
		retrieveAndUpdateCPUUsage()

		ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic, error : %v, stack : %s", err, debug.Stack())
				}
				time.Sleep(time.Second)
			}()
			for {
				select {
				case <-ticker.C:
					retrieveAndUpdateCPUUsage()
				case <-ssStopChan:
					ticker.Stop()
					return
				}
			}
		}()
	})
}

func retrieveAndUpdateCPUUsage() {
	cpuRate, err := getCPURate()
	if err != nil {
		return
	}
	currentCPUUsage.Store(cpuRate)
	slidingWindow.Add(int(cpuRate * scale))
}

func CurrentCPUUsage() float64 {
	if disablePIDLimit {
		return 0
	}
	r, ok := currentCPUUsage.Load().(float64)
	if !ok {
		return notRetrievedValue
	}
	return r
}

// checkout the cpu slice in order of time
func ExtractCPUWindows() (cpuRates []float64) {
	record := slidingWindow.GetData()
	for _, rate := range record {
		if rate == 0 {
			continue
		}
		cpuRates = append(cpuRates, float64(rate)/scale)
	}
	return cpuRates
}
