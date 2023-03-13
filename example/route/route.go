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
package  route

import (
	"errors"
	"math/rand"
	"sync/atomic"

	"github.com/bytedance/plato"
)

var (
	me = plato.DefaultEntry("getDataM")
)

type RouteRule struct {
	mRate                     int32
	mRateMin                  int32
	mRateModifyStep           int32
	errRateThreshold          float64
	recoverErrorRateThreshold float64
}

func (r *RouteRule) Decide(ctx *plato.EntryCtx) bool {

	mr := atomic.LoadInt32(&r.mRate)
	errRate := plato.Chain(me, plato.ErrRate)
	if errRate > r.errRateThreshold && mr > r.mRateMin { //Error rate> x% and read master ratio> x%, read master ratio-step
		atomic.CompareAndSwapInt32(&r.mRate, mr, mr-r.mRateModifyStep)
	} else if errRate < r.recoverErrorRateThreshold && mr < 100 { //Error rate <y% and read master ratio <100, read master ratio + step
		atomic.CompareAndSwapInt32(&r.mRate, mr, mr+r.mRateModifyStep)
	}

	mr = atomic.LoadInt32(&r.mRate)
	return rand.Int31n(100) < mr

}

func init() {
	me.Rule = &RouteRule{mRate: 100, mRateMin: 1, mRateModifyStep: 1, errRateThreshold: 0.01, recoverErrorRateThreshold: 0.005}
	plato.Init([]*plato.PlatoEntry{me})
}

func GetData() (v interface{}, e error) {
	e, _ = me.Run(func() error {
		v, e = getDataM()
		return e
	})
	//If getDataM is not executed, then getDataS is executed
	if e == plato.ErrRejectByRule {
		v, e = getDataS()
	}
	return
}

func getDataM() (interface{}, error) {
	return 1, errors.New("master error")
}

func getDataS() (interface{}, error) {
	return 2, nil
}
