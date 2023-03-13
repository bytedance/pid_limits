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
package  zscore

import (
	"fmt"
	"github.com/bytedance/plato/arithmetic/common"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestZScore(t *testing.T) {
	type args struct {
		record []float64
		score  float64
	}
	tests := []struct {
		name        string
		args        args
		wantWindows []float64
	}{
		{
			"test",
			args{
				record: []float64{0, 3.1, 8, 3, 2.3, 3.5, 3.6, 8.9, 9.3},
				score:  1.0,
			},
			[]float64{3.1, 3, 2.3, 3.5, 3.6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWindows := ZScore(tt.args.record, tt.args.score); !reflect.DeepEqual(gotWindows, tt.wantWindows) {
				t.Errorf("ZScore() = %v, want %v", gotWindows, tt.wantWindows)
			}
		})
	}
}

func TestZScore2(t *testing.T) {
	s := "1.5281353 0.6911482 0.1731542 0.8581307 0.7311312 1.5246627 0.8110686 0.1806212 0.7451248 0.7601824 0.702663 1.4230084 0.7260351 0.2555277 0.8630126 2.0968866 8.4408668 1.4632263 0.1967254 0.2879235 0.9739697 0.7640313 0.7421751 1.5682512 0.752595 0.2449329 0.8759364 0.8317237 1.4537346 0.6274189 0.2395038 0.7559726 0.776308 0.7124939 1.1583864 0.7883947 0.3943375 0.7674496 0.7440499 0.872271 1.3544351 0.6739148 0.2980677 0.8086892 0.7360282 0.7264415 0.8100931 1.5179045 0.6396336 0.236121 0.7461853 0.75243 0.6925321 1.5219135 0.7780892 0.1941268 0.755833 0.7449936 0.746055"
	ss := strings.Split(s, " ")
	f := make([]float64, 0)
	for _, s2 := range ss {
		k, _ := strconv.ParseFloat(s2, 64)
		f = append(f, k)
	}
	z := ZScore(f, 0.75)
	fmt.Println(z)
	fmt.Println(common.AverageFloat(z))
}

