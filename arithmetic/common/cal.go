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
package  common

import "math"

/**
define some common calculation, eg. sum, avg, standardDeviation, so on
 */
func SumFloat(x []float64) float64 {
	var c float64 = 0
	for _, i := range x {
		c += i
	}
	return c
}

func AverageFloat(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}
	return SumFloat(x)/float64(len(x))
}

// avg delegates the average value of x
func StandardDeviation(x []float64, avg float64) float64 {
	var variance float64
	for _,v := range x {
		variance += math.Pow(v - avg, 2)
	}
	sd := math.Sqrt(variance / float64(len(x)))
	return sd
}
