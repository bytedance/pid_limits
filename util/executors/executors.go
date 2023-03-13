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
	"log"
	"runtime"
	"runtime/debug"
	"sync/atomic"
)

type executorsMsg struct {
	c Callable
	f finishableFuture
}

func NewFixedSizeExecutorsPool(workerCount int, queueSize int) ExecutorPool {
	if workerCount <= 0 || queueSize < 0 {
		return nil
	}

	ret := &fixedSizeExecutorsPool{
		size:         workerCount,
		rp:           dropRejectPolicy{},
		callableChan: make(chan executorsMsg, queueSize),
		closeChan:    make(chan struct{}),
	}
	ret.start()
	return ret
}

type fixedSizeExecutorsPool struct {
	size         int //worker goroutines count
	rp           RejectPolicy
	callableChan chan executorsMsg
	closeChan    chan struct{}
	closed       int32
}

func (f *fixedSizeExecutorsPool) Run(runnable Runnable) error {
	//wrap as a callable
	c := func() (interface{}, error) {
		return nil, runnable()
	}
	future := f.Submit(c)
	//return if error occurs
	if cf, ok := future.(finishableFuture); ok && cf.IsCompleted() {
		_, e := future.Get(0)
		return e
	}
	return nil
}

func (f *fixedSizeExecutorsPool) Submit(callable Callable) Future {
	future := newDefaultFuture()
	if callable == nil {
		future.Done(nil, errors.New("callable is nil"))
		return future
	}

	if f.isClosed() {
		future.Done(nil, errors.New("executors pool is closed"))
		return future
	}

	m := executorsMsg{
		c: callable,
		f: future,
	}
	select {
	case f.callableChan <- m:
		return future
	default: //when reject
		f.rp.OnReject(callable, future)
		return future
	}
}

func (f *fixedSizeExecutorsPool) Close() {
	succ := atomic.CompareAndSwapInt32(&f.closed, 0, 1)
	if !succ { //avoid close closed channel
		return
	}
	close(f.closeChan)
}

func (f *fixedSizeExecutorsPool) isClosed() bool { //whether pool is closed
	return atomic.LoadInt32(&f.closed) == 1
}

func (f *fixedSizeExecutorsPool) start() {
	for i := 0; i < f.size; i++ {
		go func() {
			f.loop()
		}()
	}
}

func (f *fixedSizeExecutorsPool) loop() {
	for {
		select {
		case <-f.closeChan:
			return
		default:
			func() {
				var m executorsMsg
				defer func() {
					if e := recover(); e != nil {
						log.Printf("error: executors panic : error %v, stack : %v", e, string(debug.Stack()))
						if m.f != nil {
							m.f.Done(nil, errors.New("panic occur"))
						}
					}
				}()
				m = <-f.callableChan
				if m.f.IsCompleted() && f.rp.RejectWhenComplete(m.f) { //do nothing when future is finished
					return
				}
				c, e := m.c()
				m.f.Done(c, e)
			}()
		}
		runtime.Gosched() //release p for other goroutines
	}
}
