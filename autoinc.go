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

import (
	"errors"
	"math"
)

// 常用的错误类型
var (
	ErrOverflow = errors.New("数值溢出")
	ErrNotFound = errors.New("未获取得正确的 id 值")
)

// AutoInc 用于产生唯一 ID。
type AutoInc struct {
	start, step int64
	channel     chan int64
	done        chan struct{}

	err   error
	faild bool
}

// New 声明一个新的 AutoInc 实例。
//
// start：起始数值；
// step：步长，可以为负数，但不能为 0；
// bufferSize；缓存的长度。
//
// 如果 step 为 0，会直接 panic
func New(start, step, bufferSize int64) *AutoInc {
	if step == 0 {
		panic("无效的参数 step")
	}

	ai := &AutoInc{
		start:   start,
		step:    step,
		channel: make(chan int64, bufferSize),
		done:    make(chan struct{}),
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
		case ai.channel <- ai.start:
			if (ai.step > 0 && ai.start > 0 && (math.MaxInt64-ai.start) < ai.step) ||
				(ai.step < 0 && ai.start < 0 && (-math.MaxInt64-ai.start) > ai.step) {
				ai.err = ErrOverflow
				// 此时不能关闭 channel，其中依然有值。即不能设置 done，也不能 close(channel)
				return
			}

			ai.start += ai.step
		}
	}
}

// ID 获取 ID 值。
func (ai *AutoInc) ID() (int64, error) {
	if ai.faild {
		return 0, ai.err
	}

	ret, ok := <-ai.channel

	if !ok {
		if ai.err != nil {
			return 0, ai.err
		}

		return 0, ErrNotFound
	}

	// 如果仅判断 channel 是否为空，则有可能是写 channel 的速度没有读 channel 的快。
	if ai.err != nil && len(ai.channel) == 0 {
		ai.faild = true
	}
	return ret, nil
}

// MustID 获取 ID 值，若不成功，则 panic。
func (ai *AutoInc) MustID() int64 {
	id, err := ai.ID()
	if err != nil {
		panic(err)
	}

	return id
}

// Stop 停止计时
func (ai *AutoInc) Stop() {
	ai.done <- struct{}{}
}
