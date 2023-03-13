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
package  executors

import "testing"

func Test_dropRejectPolicy_RejectWhenComplete(t *testing.T) {
	type args struct {
		in0 Complete
	}
	tests := []struct {
		name string
		d    dropRejectPolicy
		args args
		want bool
	}{
		{
			"test",
			dropRejectPolicy{},
			args{in0: newDefaultFuture()},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := dropRejectPolicy{}
			if got := d.RejectWhenComplete(tt.args.in0); got != tt.want {
				t.Errorf("dropRejectPolicy.RejectWhenComplete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dropRejectPolicy_OnReject(t *testing.T) {
	type args struct {
		callable Callable
		f        Complete
	}
	tests := []struct {
		name string
		d    dropRejectPolicy
		args args
	}{
		{
			"test",
			dropRejectPolicy{},
			args{
				callable: func() (interface{}, error) {
					return "ok", nil
				},
				f:        newDefaultFuture(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := dropRejectPolicy{}
			d.OnReject(tt.args.callable, tt.args.f)
		})
	}
}
