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
package zscore

import (
	"github.com/bytedance/pid_limits/arithmetic/common"
	"math"
)

// use z-score to eliminate the abnormal value
func ZScore(record []float64, score float64) (windows []float64) {
	u := common.AverageFloat(record)
	sd := common.StandardDeviation(record, u)
	windows = make([]float64, 0)
	count := 0
	var total float64 = 0
	for _, v := range record {
		if zScore := (v - u) / sd; math.Abs(zScore) < score {
			// the value of z-score less than score should be normal
			// recommend the score value is 2.2
			count++
			total += v
			windows = append(windows, v)
		}
	}
	if count == 0 {
		return windows
	}
	return windows
}
