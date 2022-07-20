// SPDX-License-Identifier: MIT

// Package autoinc 用于产生唯一 ID
//
//  ai := autoinc.New(0, 1, 1)
//  for i:=0; i<10; i++ {
//      fmt.Println(ai.ID())
//  }
//
//  ai.Stop()
package autoinc

import "math"

// AutoInc 用于产生唯一 ID
type AutoInc struct {
	start, step int64
	channel     chan int64
	done        chan struct{}
}

// New 声明一个新的 AutoInc 实例
//
// start：起始数值；
// step：步长，可以为负数，但不能为 0；
// bufferSize：缓存的长度。
//
// 如果 step 为 0，会直接 panic
func New(start, step int64, bufferSize int) *AutoInc {
	if step == 0 {
		panic("无效的参数 step")
	}

	ai := &AutoInc{
		start:   start,
		step:    step,
		channel: make(chan int64, bufferSize),
		done:    make(chan struct{}, 1),
	}

	go ai.generator()

	return ai
}

func (ai *AutoInc) generator() {
	for {
		select {
		case <-ai.done:
			close(ai.channel)
			return
		case ai.channel <- ai.start: // 在 channel 未满之前，此条一直有效
			if (ai.step > 0 && ai.start > 0 && (math.MaxInt64-ai.start) < ai.step) ||
				(ai.step < 0 && ai.start < 0 && (math.MinInt64-ai.start) > ai.step) {
				close(ai.channel)
				return
			}

			ai.start += ai.step
		}
	}
}

// ID 获取 ID 值
//
// 第二个参数若返回 false，表示当前的 ID 值已经失效。
func (ai *AutoInc) ID() (int64, bool) {
	ret, ok := <-ai.channel
	return ret, ok
}

// MustID 获取 ID 值
//
// 与 ID() 的不同在于，出错时会直接 panic。
func (ai *AutoInc) MustID() int64 {
	id, ok := ai.ID()
	if !ok {
		panic("当前已经停止分发新的 ID")
	}

	return id
}

// Stop 停止服务
//
// NOTE: 多次调用，会造成死锁。
func (ai *AutoInc) Stop() { ai.done <- struct{}{} }
