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
package limiting

import (
	"github.com/bytedance/pid_limits/application/adaptive/config"
	"github.com/go-playground/assert/v2"
	"testing"
	"time"
)

func TestNewPidLimitingHttpDefault(t *testing.T) {
	limiting := NewPidLimitingHttpDefault(0.8, config.WithDisableMetric())
	time.Sleep(2 * time.Second)
	assert.Equal(t, false, limiting.Limit())
}
