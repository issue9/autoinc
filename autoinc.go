// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package autoinc

// AutoInc用于产生唯一ID
type AutoInc struct {
	start, step int64
	channel     chan int64
}

// 声明一个新的AutoInc实例
// start：起始数值；step：步长；bufferSize；缓存的长度。
func New(start, step, bufferSize int64) *AutoInc {
	ret := &AutoInc{
		start:   start,
		step:    step,
		channel: make(chan int64, bufferSize),
	}

	go func() {
		for i := ret.start; true; i += ret.step {
			ret.channel <- i
		}
	}()

	return ret
}

// 获取ID值
func (ai *AutoInc) ID() int64 {
	return <-ai.channel
}
