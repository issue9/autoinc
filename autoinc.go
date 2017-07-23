// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package autoinc 用于产生唯一 ID，可以指定起始 ID 和步长。
//  ai := autoinc.New(0, 1, 1)
//  for i:=0; i<10; i++ {
//      fmt.Println(ai.ID())
//  }
//
//  ai.Stop()
package autoinc

// AutoInc 用于产生唯一 ID。
type AutoInc struct {
	start, step int64
	channel     chan int64
	done        chan struct{}
}

// New 声明一个新的 AutoInc 实例。
//
// start：起始数值；step：步长；bufferSize；缓存的长度。
func New(start, step, bufferSize int64) *AutoInc {
	ret := &AutoInc{
		start:   start,
		step:    step,
		channel: make(chan int64, bufferSize),
		done:    make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-ret.done:
				close(ret.channel)
				return
			case ret.channel <- ret.start:
				ret.start += ret.step
			}
		}
	}()

	return ret
}

// ID 获取 ID 值。若已经调用 Stop，则之后的 ID 值不保证正确。
func (ai *AutoInc) ID() (int64, bool) {
	ret, ok := <-ai.channel
	return ret, ok
}

// MustID 获取 ID 值，若不成功，则返回零值。
func (ai *AutoInc) MustID() int64 {
	return <-ai.channel
}

// Stop 停止计时
func (ai *AutoInc) Stop() {
	ai.done <- struct{}{}
}

// Reset 重置 start 和 step
func (ai *AutoInc) Reset(start, step int64) {
	ai.start = start
	ai.step = step
}
