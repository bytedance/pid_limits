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

import "time"

type Runnable func() error

type Callable func() (interface{}, error)

//Define the behavior of the caller to get the returned result
type Future interface {
	Get(duration time.Duration) (interface{}, error)
}

//Define the behavior of task execution completion
type Complete interface {
	Done(interface{}, error)
	IsCompleted() bool
}

//A coroutine pool
type ExecutorPool interface {
	Run(runnable Runnable) error
	Submit(callable Callable) Future
	Close()
}

//Define the behavior when the pool rejects
type RejectPolicy interface {
	OnReject(callable Callable, f Complete)
	RejectWhenComplete(Complete) bool //whether reject the callable when task is completed(i.e. when timeout)
}
