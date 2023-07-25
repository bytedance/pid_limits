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
package adaptive

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bytedance/pid_limits/application/adaptive/config"
	"github.com/bytedance/pid_limits/application/adaptive/limiting"
	"github.com/bytedance/pid_limits/arithmetic/pid"
	"github.com/bytedance/pid_limits/core/system"
	"github.com/gin-gonic/gin"
)

type tunePID struct {
	rate float64
}

func (t *tunePID) initTuner(threshold float64) {
	tuner := &pid.Tuner{}
	tuner.Init(threshold)
	go func() {
		for {
			t.rate = tuner.TunePID(system.CurrentCPUUsage())
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func PlatoMiddlewareGinDefault(threshold float64, opts ...config.OptionFunc) gin.HandlerFunc {
	limit := limiting.NewPidLimitingHttpDefault(threshold, opts...)
	return func(c *gin.Context) {
		if limit.Limit() {
			_ = c.AbortWithError(510, fmt.Errorf("block by pid"))
			return
		}
		c.Next()
	}
}

func PlatoMiddlewareGin(kp float64, ki float64, kd float64, threshold float64, opts ...config.OptionFunc) gin.HandlerFunc {
	limit := limiting.NewPidLimiting(kp, ki, kd, threshold, opts...)
	return func(c *gin.Context) {
		if limit.Limit() {
			_ = c.AbortWithError(510, fmt.Errorf("block by pid"))
			return
		}
		c.Next()
	}
}

func TunePIDMiddlewareGin(threshold float64) gin.HandlerFunc {
	tPID := &tunePID{rate: 0}
	tPID.initTuner(threshold)
	return func(c *gin.Context) {
		if rand.Intn(10000) < int(-tPID.rate) {
			_ = c.AbortWithError(510, fmt.Errorf("block by pid"))
			return
		}
		c.Next()
	}
}
