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
	//
	// 所有的 ID 值是提交保存在 channel 中的，并不是一发生数据溢出了，
	// 马上 ID() 就不可用了。而是要等待 ID() 将 channel 中的值取完，
	// 才会出返回错误信息。
	err    error
	errVal int64
	faild  bool
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
	if ai.faild {
		return 0, ai.err
	}

	ret, ok := <-ai.channel

	if !ok {
		return 0, ErrNotFound
	}

	// NOTE: 不能只通过判断 errVal 与 id 的值来确定是否出错，两都有可能是 0。
	// ai.errVal 值本身是可用的。可以正常返回。
	if ai.err != nil && ai.errVal == ret {
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
