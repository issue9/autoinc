// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// autoinc用于产生唯一ID，可以指定起始ID和步长。
//  ai := autoinc.New(0, 1, 1)
//  for i:=0; i<10; i++ {
//      fmt.Println(ai.ID())
//  }
package autoinc

// AutoInc用于产生唯一ID。
// AutoInc实例一旦声明，就无法关闭，所以并不是很适合短期的服务。
type AutoInc struct {
	start, step int64
	channel     chan int64
}

// 声明一个新的AutoInc实例。
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

var defaultAI = New(1, 1, 100)

// 获取从 1 开始，步长为 1 的自增 ID 值
func ID() int64 {
	return defaultAI.ID()
}
