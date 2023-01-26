// SPDX-License-Identifier: MIT

// Package autoinc 用于产生自增 ID
//
//	ai := autoinc.New(0, 1, 1)
//	ctx, cancel = context.WithCancel(context.Background())
//	defer cancel()
//
//	go ai.Serve(ctx)
//	for i:=0; i<10; i++ {
//	    fmt.Println(ai.ID())
//	}
package autoinc

import (
	"context"
	"errors"
	"math"
)

var errOverflow = errors.New("溢出")

// ErrOverflow 表示自增 ID 溢出了
func ErrOverflow() error { return errOverflow }

type AutoInc struct {
	start, step int64
	channel     chan int64
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

	return &AutoInc{
		start:   start,
		step:    step,
		channel: make(chan int64, bufferSize),
	}
}

// Serve 运行该服务
//
// 这是个阻塞方法，只有此方法运行之后，[AutoInc.ID] 等才有返回值。
//
// 当前自增项超过最大值时，会返回 [ErrOverflow] 错误。
func (ai *AutoInc) Serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			close(ai.channel)
			return ctx.Err()
		case ai.channel <- ai.start:
			if (ai.step > 0 && ai.start > 0 && (math.MaxInt64-ai.start) < ai.step) ||
				(ai.step < 0 && ai.start < 0 && (math.MinInt64-ai.start) > ai.step) {
				close(ai.channel)
				return ErrOverflow()
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
