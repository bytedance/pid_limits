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

//Chain will return cached value of target metric. It return 0 when input PlatoEntry does not have target metric, or input PlatoEntry has not been pass to plato.Init method
func Chain(entry *PlatoEntry, metricType MetricFactory) float64 {
	m := entry.Metrics[metricType]
	if m == nil {
		return 0
	}
	return m.getValue()
}

//Calculate will calculate and return the value of target metric. It return 0 when input PlatoEntry does not have target metric
func Calculate(entry *PlatoEntry, metricType MetricFactory) float64 {
	m := entry.Metrics[metricType]
	if m == nil {
		return 0
	}
	return m.cul()
}

type RuleInterface interface {
	Decide(ctx *EntryCtx) bool
}
