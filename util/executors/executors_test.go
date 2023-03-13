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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFixedSizeExecutorsPool(t *testing.T) {
	p := NewFixedSizeExecutorsPool(2, 4)
	f := p.Submit(func() (interface{}, error) {
		return 1, nil
	})
	val, e := f.Get(1 * time.Millisecond)
	assert.Nil(t, e)
	assert.True(t, val.(int) == 1)
	f = p.Submit(func() (interface{}, error) {
		time.Sleep(time.Second)
		return 1, nil
	})
	get, e := f.Get(1 * time.Millisecond)
	assert.Nil(t, get)
	assert.NotNil(t, e)
	t.Log(e)
}

func TestNewFixedSizeExecutorsPool2(t *testing.T) {
	p := NewFixedSizeExecutorsPool(2, -2)
	assert.Nil(t, p, "should be nil")
}

func TestExecutorsPanic(t *testing.T) {
	p := NewFixedSizeExecutorsPool(2, 0)
	time.Sleep(1 * time.Second)
	f := p.Submit(func() (i interface{}, e error) {
		panic("e")
		return nil, nil
	})
	_, e := f.Get(1 * time.Second)
	t.Log(e)
}

//go test -v -race -test.run TestConcurrent
func TestConcurrent(t *testing.T) {
	pool := NewFixedSizeExecutorsPool(2, 4)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var iwg sync.WaitGroup
			iwg.Add(1)
			f := pool.Submit(func() (interface{}, error) {
				defer iwg.Done()
				return 1, nil
			})
			v, _ := f.Get(10 * time.Millisecond)
			t.Log(v)
			iwg.Add(1)
			_ = pool.Run(func() error {
				defer iwg.Done()
				time.Sleep(1 * time.Second)
				t.Log(2)
				return nil
			})
		}()
	}

	wg.Wait()
	time.Sleep(10 * time.Second)
}

func Test_fixedSizeExecutorsPool_Submit(t *testing.T) {
	type fields struct {
		size         int
		rp           RejectPolicy
		callableChan chan executorsMsg
		closeChan    chan struct{}
		closed       int32
	}
	type args struct {
		callable Callable
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Future
	}{
		{
			"test",
			fields{
				size:         10,
				rp:           dropRejectPolicy{},
				callableChan: make(chan executorsMsg, 1000),
				closeChan:    make(chan struct{}),
				closed: int32(0),
			},
			args{callable: func() (interface{}, error) {
				return "ok", nil
			}},
			newDefaultFuture(),
		},
		{
			"test2",
			fields{
				size:         10,
				rp:           dropRejectPolicy{},
				callableChan: make(chan executorsMsg, 1000),
				closeChan:    make(chan struct{}),
				closed: int32(0),
			},
			args{callable: nil},
			newDefaultFuture(),
		},
		{
			"test",
			fields{
				size:         10,
				rp:           dropRejectPolicy{},
				callableChan: make(chan executorsMsg, 1000),
				closeChan:    make(chan struct{}),
				closed: int32(1),
			},
			args{callable: func() (interface{}, error) {
				return "ok", nil
			}},
			newDefaultFuture(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fixedSizeExecutorsPool{
				size:         tt.fields.size,
				rp:           tt.fields.rp,
				callableChan: tt.fields.callableChan,
				closeChan:    tt.fields.closeChan,
				closed:       tt.fields.closed,
			}
			if got := f.Submit(tt.args.callable); false {
				t.Errorf("fixedSizeExecutorsPool.Submit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fixedSizeExecutorsPool_Close(t *testing.T) {
	type fields struct {
		size         int
		rp           RejectPolicy
		callableChan chan executorsMsg
		closeChan    chan struct{}
		closed       int32
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"test",
			fields{
				size:         10,
				rp:           dropRejectPolicy{},
				callableChan: make(chan executorsMsg, 1000),
				closeChan:    make(chan struct{}),
				closed: int32(0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fixedSizeExecutorsPool{
				size:         tt.fields.size,
				rp:           tt.fields.rp,
				callableChan: tt.fields.callableChan,
				closeChan:    tt.fields.closeChan,
				closed:       tt.fields.closed,
			}
			f.Close()
		})
	}
}

func Test_fixedSizeExecutorsPool_isClosed(t *testing.T) {
	type fields struct {
		size         int
		rp           RejectPolicy
		callableChan chan executorsMsg
		closeChan    chan struct{}
		closed       int32
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"test",
			fields{
				size:         10,
				rp:           dropRejectPolicy{},
				callableChan: make(chan executorsMsg, 1000),
				closeChan:    make(chan struct{}),
				closed: int32(0),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fixedSizeExecutorsPool{
				size:         tt.fields.size,
				rp:           tt.fields.rp,
				callableChan: tt.fields.callableChan,
				closeChan:    tt.fields.closeChan,
				closed:       tt.fields.closed,
			}
			if got := f.isClosed(); got != tt.want {
				t.Errorf("fixedSizeExecutorsPool.isClosed() = %v, want %v", got, tt.want)
			}
		})
	}
}
