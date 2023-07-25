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
package plato

import "github.com/bytedance/pid_limits/util"

type EntryCtx struct {
	startTime uint64
	pe        *PlatoEntry
}

func (c *EntryCtx) GetEntry() *PlatoEntry {
	return c.pe
}

func (c *EntryCtx) GetStartTime() uint64 {
	return c.startTime
}

func NewCtx(entry *PlatoEntry) *EntryCtx {
	return &EntryCtx{
		startTime: util.CurrentTimeMillis(),
		pe:        entry,
	}
}
