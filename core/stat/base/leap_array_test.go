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
package  base

import (
	"fmt"
	"testing"
)
func test(i int) {
	var arr [10]int
	//优先使用错误拦截 在错误出现之前进行拦截 在错误出现后进行错误捕获
	//错误拦截必须配合defer使用  通过匿名函数使用
	//defer func() {
	//	//恢复程序的控制权
	//	err := recover()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}()

	arr[i] = 123 //err panic
	fmt.Println(arr)
}
func TestA(t *testing.T) {
	defer func() {
		//恢复程序的控制权
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	i := 10
	test(i)
	fmt.Println("hello world")
}
