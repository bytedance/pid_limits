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
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	cGroupPath    = "/proc/self/cgroup"
	prefixPath    = "/sys/fs/cgroup/cpu,cpuacct"
	cpuAccounting = "cpu,cpuacct"

	notRetrievedPath = ""

	usageLog  = "/cpuacct.usage"
	periodLog = "/cpu.cfs_period_us"
	quotaLog  = "/cpu.cfs_quota_us"
)

type cGroupStat struct {
	cpuUsage  int64
	timeStamp int64
}

var (
	prevCGroupStat cGroupStat

	dockerCPUUsagePath  string
	dockerCPUPeriodPath string
	dockerCPUQuotaPath  string
	cfsPeriod           float64
	cfsQuota            float64
)

func extractSystemPath(line string) string {
	for _, line := range strings.Split(line, "\n") {
		parts := strings.Split(line, ":")
		if len(parts) == 3 && parts[1] == cpuAccounting {
			return parts[2]
		}
	}
	return notRetrievedPath
}

func getPodSystemPath() string {
	data, err := ioutil.ReadFile(cGroupPath)
	if err != nil {
		log.Printf("error: read /proc/self/cgroup error: %v\n", err)
		return notRetrievedPath
	}
	return extractSystemPath(string(data))
}

// use bufio to read file should be better than ioutil
func getDockerSystemMetric(filepath string) int64 {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("error: fetch file from path %v error: %v\n", filepath, err)
		return 0
	}
	var value int64
	line := strings.ReplaceAll(string(data), "\n", "")
	if value, err = strconv.ParseInt(line, 10, 64); err != nil {
		log.Printf("error: parse cpu usage file error %v\n", err)
		return 0
	}
	return value
}

func initCGroup() bool {
	// todo use file join instead
	podPath := getPodSystemPath()

	dockerCPUUsagePath = filepath.Join(prefixPath, podPath, usageLog)

	dockerCPUPeriodPath = filepath.Join(prefixPath, podPath, periodLog)

	dockerCPUQuotaPath = filepath.Join(prefixPath, podPath, quotaLog)

	cfsPeriod = float64(getDockerSystemMetric(dockerCPUPeriodPath))

	cfsQuota = float64(getDockerSystemMetric(dockerCPUQuotaPath))

	prevCGroupStat = cGroupStat{
		cpuUsage:  0,
		timeStamp: 0,
	}

	return fileExist(dockerCPUUsagePath) && fileExist(dockerCPUPeriodPath) && fileExist(dockerCPUQuotaPath) &&
		cfsQuota != -1 && cfsQuota != 0 && cfsPeriod != 0
}

// check file exist or not
func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func getCPURateByCGroup() (rate float64, err error) {
	usage := getDockerSystemMetric(dockerCPUUsagePath)
	if usage == prevCGroupStat.cpuUsage {
		// 打日志
		log.Println("error: retrieveValueFailed")
		return retrieveValueFailed, retrieveValueError
	}
	now := time.Now().UnixNano()
	if prevCGroupStat.timeStamp != 0 && now != prevCGroupStat.timeStamp {
		prevUsage := prevCGroupStat.cpuUsage
		preTimeStamp := prevCGroupStat.timeStamp
		rate = float64(usage-prevUsage) / float64(now-preTimeStamp) * cfsPeriod / cfsQuota
	}
	prevCGroupStat.timeStamp = now
	prevCGroupStat.cpuUsage = usage
	rate = math.Max(0.0, rate)
	//rate = math.Min(1.0, rate)
	return rate, nil
}
