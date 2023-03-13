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

import (
	"errors"
	"sync/atomic"
	"time"
)

type finishableFuture interface {
	Complete
	Future
}

type defaultFutureMsg struct {
	e error
	r interface{}
}

func newDefaultFuture() finishableFuture {
	f := &defaultFuture{resultChan: make(chan defaultFutureMsg, 1)}
	return f
}

type defaultFuture struct {
	resultChan chan defaultFutureMsg
	finished   int32
}

func (d *defaultFuture) Get(duration time.Duration) (interface{}, error) {
	if d.IsCompleted() { //directly return when finished
		r := <-d.resultChan
		return r.r, r.e
	}

	t := time.NewTimer(duration)
	select {
	case r := <-d.resultChan:
		t.Stop()
		return r.r, r.e
	case <-t.C:
		e := errors.New("timeout")
		d.Done(nil, e) //manually done to save resource
		return nil, e
	}
}

func (d *defaultFuture) Done(r interface{}, e error) {
	if d.IsCompleted() {
		return
	}
	//avoid block current goroutine
	select {
	case d.resultChan <- defaultFutureMsg{
		e: e,
		r: r,
	}:
		atomic.SwapInt32(&d.finished, 1)
	default: //do nothing
	}
}

func (d *defaultFuture) IsCompleted() bool {
	return atomic.LoadInt32(&d.finished) == 1
}
