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
package  plato

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/bytedance/plato/util"
)

func TestEntryCtx_GetEntry(t *testing.T) {
	ec := DefaultEntryCalculated("test")
	type fields struct {
		startTime uint64
		pe        *PlatoEntry
	}
	tests := []struct {
		name   string
		fields fields
		want   *PlatoEntry
	}{
		{
			"test",
			fields{
				startTime: util.CurrentTimeMillis(),
				pe:        ec,
			},
			ec,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EntryCtx{
				startTime: tt.fields.startTime,
				pe:        tt.fields.pe,
			}
			if got := c.GetEntry(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EntryCtx.GetEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryCtx_GetStartTime(t *testing.T) {
	ec := DefaultEntryCalculated("test")
	type fields struct {
		startTime uint64
		pe        *PlatoEntry
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			"test",
			fields{
				startTime: util.CurrentTimeMillis(),
				pe:        ec,
			},
			util.CurrentTimeMillis(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EntryCtx{
				startTime: tt.fields.startTime,
				pe:        tt.fields.pe,
			}
			if got := c.GetStartTime(); got != tt.want {
				t.Errorf("EntryCtx.GetStartTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCtx(t *testing.T) {
	ec := DefaultEntryCalculated("test")
	ctx, _ := ec.Entry()
	type args struct {
		entry *PlatoEntry
	}
	tests := []struct {
		name string
		args args
		want *EntryCtx
	}{
		{
			"test",
			args{entry: ec},
			ctx,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCtx(tt.args.entry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}


