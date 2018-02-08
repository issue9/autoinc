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

	// 记录错误信息
	//
	// 自增有可能触发溢出等错误，一旦发生，则需要记录该错误，
	// 以及触发该错误时的 ID 值。
	err    error
	errVal int64
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
				if (ret.step > 0 && ret.start > 0 && (math.MaxInt64-ret.start) < ret.step) ||
					(ret.step < 0 && ret.start < 0 && (-math.MaxInt64-ret.start) > ret.step) {
					ret.err = ErrOverflow
					ret.errVal = ret.start
					return
				}

				ret.start += ret.step
			}
		} // end for
	}()

	return ret
}

// ID 获取 ID 值。
func (ai *AutoInc) ID() (int64, error) {
	ret, ok := <-ai.channel

	// NOTE: 不能只通过判断 errVal 与 id 的值来确定是否出错，两都有可能是 0
	if ai.err != nil && ai.errVal == ret {
		return 0, ai.err
	}

	if !ok {
		return 0, ErrNotFound
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

// Reset 重置整个计数器。
//
// 计数器将根据新的参数重新运行。但是已经产生的数据不会回收。
func (ai *AutoInc) Reset(start, step int64) {
	if step == 0 {
		panic("参数 step 不能为 0")
	}

	ai.start = start
	ai.step = step
	ai.err = nil
	ai.errVal = 0
}
