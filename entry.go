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

import (
	"errors"
	"time"

	"github.com/bytedance/pid_limits/core/base"
	"github.com/bytedance/pid_limits/core/stat"
	"github.com/bytedance/pid_limits/util"
)

var (
	ErrRejectByRule = errors.New("reject by rule")
)

//A function create a *PlatoEntry that has metrics represented by calculatedMetrics.
//Metrics of this *PlatoEntry will not be calculated until they are passed to plato.Init() method
func NewPlatoEntry(name string, calculateMetrics ...MetricFactory) *PlatoEntry {
	pe := &PlatoEntry{
		Rule:       nil,
		Name:       name,
		rtt:        stat.NewSlidingWindow(-1, time.Millisecond*time.Duration(base.DefaultIntervalMs)),
		completion: stat.NewSlidingWindow(-1, time.Millisecond*time.Duration(base.DefaultIntervalMs)),
		blocked:    stat.NewSlidingWindow(-1, time.Millisecond*time.Duration(base.DefaultIntervalMs)),
		error:      stat.NewSlidingWindow(-1, time.Millisecond*time.Duration(base.DefaultIntervalMs)),
	}

	for _, m := range calculateMetrics {
		pe.AddMetric(m)
	}

	return pe
}

//A function create a *PlatoEntry that has metrics of pct99 latency、avg latency、error rate and qps.
//Metrics of this *PlatoEntry will not be calculated in background until they are passed to plato.Init() method
func DefaultEntry(name string) *PlatoEntry {
	return NewPlatoEntry(name, PctRT, AvgRT, ErrRate, Qps)
}

//A function create a *PlatoEntry whose metrics passed by calculatedMetrics will be calculated
func NewEntryCalculated(name string, calculateMetrics ...MetricFactory) *PlatoEntry {
	entry := NewPlatoEntry(name, calculateMetrics...)
	Init([]*PlatoEntry{entry})
	return entry
}

//A function create a *PlatoEntry that its pct99 latency、avg latency、error rate and qps will be calculated in background goroutine..
func DefaultEntryCalculated(name string) *PlatoEntry {
	return NewEntryCalculated(name, PctRT, AvgRT, ErrRate, Qps)
}

type PlatoEntry struct {
	Rule       RuleInterface
	Name       string
	Metrics    map[MetricFactory]*Metric //todo : maybe we can use uintptr as key here
	rtt        *stat.SlidingWindow
	completion *stat.SlidingWindow
	blocked    *stat.SlidingWindow
	error      *stat.SlidingWindow
}

//A helper function to use entry, avoid typo error.
//Error return by r will be used to ReportError
func (pe *PlatoEntry) Run(r func() error) (error, bool) {
	c, b := pe.Entry()
	if !b {
		return ErrRejectByRule, b
	}

	e := r()
	pe.Exit(c)
	if e != nil {
		pe.ReportError(e)
	}
	return e, b
}

func (pe *PlatoEntry) Exit(ctx *EntryCtx) {
	rt := util.CurrentTimeMillis() - ctx.startTime
	pe.rtt.Add(int(rt))
	pe.completion.Add(1)
}

func (pe *PlatoEntry) Entry() (*EntryCtx, bool) {
	ctx := NewCtx(pe)
	ctx.startTime = util.CurrentTimeMillis()

	if pe.Rule == nil {
		return ctx, true
	} else if flag := pe.Rule.Decide(ctx); !flag {
		pe.blocked.Add(1)
		return ctx, false
	} else {
		return ctx, true
	}
}

func (pe *PlatoEntry) ReportError(err error) {
	pe.error.Add(1)
}

func (pe *PlatoEntry) AddMetric(factory MetricFactory) {
	if pe.Metrics == nil {
		return
	}
	pe.Metrics[factory] = factory.CreateMetric(pe)
}
