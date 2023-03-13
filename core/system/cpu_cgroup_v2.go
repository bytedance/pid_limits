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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	cGroupV2Path       = "/proc/1/cgroup"
	cGroupV2prefixPath = "/sys/fs/cgroup"
	nullString         = ""

	scopePath = "/init.scope"

	cpuStatFile = "cpu.stat"
	cpuMaxFile  = "cpu.max"
)

var (
	dockerCPUStatPath string
	dockerCPUMaxPath  string
)

func extractSystemPathByCGroupV2(line string) string {
	for _, line := range strings.Split(line, "\n") {
		parts := strings.Split(line, ":")
		if len(parts) == 3 && parts[1] == nullString {
			return strings.TrimSuffix(parts[2], scopePath)
		}
	}
	return notRetrievedPath
}

func getPodSystemPathByCGroupV2() string {
	data, err := ioutil.ReadFile(cGroupV2Path)
	if err != nil {
		log.Printf("read /proc/1/cgroup error: %v\n", err)
		return notRetrievedPath
	}
	return extractSystemPathByCGroupV2(string(data))
}

func getCPUQuota() (cfsQuota, cfsPeriod float64, err error) {
	if !fileExist(dockerCPUMaxPath) {
		return 0, 0, fmt.Errorf("[getCPUQuota] docker cpu.max path does not exist, path:%s", dockerCPUMaxPath)
	}
	data, err := ioutil.ReadFile(dockerCPUMaxPath)
	if err != nil {
		return 0, 0, fmt.Errorf("[getCPUQuota] fetch file from path %v error: %v", dockerCPUMaxPath, err)
	}
	line := strings.ReplaceAll(string(data), "\n", "")
	values := strings.Split(line, " ")
	if len(values) != 2 {
		return 0, 0, fmt.Errorf("[getCPUQuota] cpu.max returns more than 2 values, cpu.max_path=%v, value=%s, error: %v", dockerCPUMaxPath, line, err)
	}
	quota, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("[getCPUQuota] parse quota value failed, line:%s, err:%v", line, err)
	}
	period, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("[getCPUQuota] parse period value failed, line:%s, err:%v", line, err)
	}
	return quota, period, nil
}

func initCGroupV2() bool {
	// todo use file join instead
	podPath := getPodSystemPathByCGroupV2()

	dockerCPUMaxPath = filepath.Join(cGroupV2prefixPath, podPath, cpuMaxFile)
	dockerCPUStatPath = filepath.Join(cGroupV2prefixPath, podPath, cpuStatFile)

	var err error
	cfsQuota, cfsPeriod, err = getCPUQuota()
	if err != nil {
		log.Printf("error: initCGroup failed, error:%v", err)
		return false
	}
	prevCGroupStat = cGroupStat{
		cpuUsage:  0,
		timeStamp: 0,
	}
	return fileExist(dockerCPUMaxPath) && fileExist(dockerCPUStatPath) && cfsQuota != -1 && cfsQuota != 0 && cfsPeriod != 0
}

func readCPUUsageByCPUStat() (usage int64, err error) {
	if !fileExist(dockerCPUStatPath) {
		return 0, fmt.Errorf("[getCPUQuota] docker cpu stat does not exist, path:%s", dockerCPUStatPath)
	}
	data, err := ioutil.ReadFile(dockerCPUStatPath)
	if err != nil {
		return 0, fmt.Errorf("[getCPUQuota] fetch file from path %v error: %v", dockerCPUMaxPath, err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		data := strings.Split(line, " ")
		if len(data) != 2 {
			continue
		}
		if data[0] == "usage_usec" {
			if value, err := strconv.ParseInt(data[1], 10, 64); err != nil {
				return 0, err
			} else {
				return value, nil
			}
		}
	}
	return 0, errors.New("not found usage_usec")
}

func getCPURateByCGroupV2() (rate float64, err error) {
	usage, err := readCPUUsageByCPUStat()
	if err != nil {
		log.Printf("error: read cpu usage by stat failed, err: %v", err)
		return retrieveValueFailed, err
	}
	if usage == prevCGroupStat.cpuUsage {
		log.Println("error: retrieveValueFailed")
		return retrieveValueFailed, retrieveValueError
	}
	// UnixMicro() 方法是在 go 1.17 之后引入的，为了兼容之前的版本，这里不直接使用该方法
	now := time.Now().UnixNano() / 1000
	if prevCGroupStat.timeStamp != 0 && now != prevCGroupStat.timeStamp {
		prevUsage := prevCGroupStat.cpuUsage
		preTimeStamp := prevCGroupStat.timeStamp
		rate = float64(usage-prevUsage) / float64(now-preTimeStamp) * cfsPeriod / cfsQuota
	}
	prevCGroupStat.timeStamp = now
	prevCGroupStat.cpuUsage = usage
	rate = math.Max(0.0, rate)
	return rate, nil
}
