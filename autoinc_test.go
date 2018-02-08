// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package autoinc

import (
	"math"
	"sync"
	"testing"
	"time"

	"github.com/issue9/assert"
)

func TestAutoInc_overflow(t *testing.T) {
	a := assert.New(t)

	ai := New(math.MaxInt64-1, 2, 2)
	a.NotNil(ai)
	id, err := ai.ID()
	a.Error(err).Equal(id, 0)
}

func TestAutoInc_ID_1(t *testing.T) {
	a := assert.New(t)

	a.Panic(func() {
		ai := New(0, 0, 2)
		a.Nil(ai)
	})

	// 正规的ai操作
	ai := New(0, 2, 2)
	a.NotNil(ai)
	for i := 0; i < 7; i++ {
		a.Equal(ai.MustID(), i*2)
	}
	ai.Stop()

	// 可以从负数起始
	ai = New(-100, 2, 5)
	a.NotNil(ai)
	for i := 0; i < 7; i++ {
		a.Equal(ai.MustID(), -100+i*2)
	}

	// start,step 双负数
	ai = New(-100, -3, 0)
	a.NotNil(ai)
	for i := 0; i < 7; i++ {
		a.Equal(ai.MustID(), -100+i*-3)
	}

	ai.Reset(99, 2)
	ai.MustID()
	for i := 0; i < 7; i++ {
		a.Equal(ai.MustID(), 99+2+i*2) // 第一个值跳过
	}
}

func TestAutoInc_ID_2(t *testing.T) {
	a := assert.New(t)

	ai := New(2, 2, 2)
	a.NotNil(ai)

	mu := sync.Mutex{}
	mapped := map[int64]bool{}

	fn := func() {
		for i := 0; i < 100; i++ {
			id := ai.MustID()

			mu.Lock()
			_, found := mapped[id]
			a.False(found, "找到重复元素:%v", id)
			mapped[id] = true
			mu.Unlock()
		}
	}

	go fn()
	go fn()
	go fn()
	go fn()
}

func TestAutoInc_Stop(t *testing.T) {
	a := assert.New(t)

	ai := New(0, 1, 2)
	a.NotNil(ai)
	ai.Stop()

	println("stop1")
	for {
		id, err := ai.ID()
		if err != nil {
			break
		}
		println(id)
	}

	ai = New(0, 2, 100)
	a.NotNil(ai)
	time.AfterFunc(20*time.Microsecond, func() {
		ai.Stop()
	})
	println("stop2")
	for {
		id, err := ai.ID()
		if err != nil {
			break
		}
		println(id)
	}
}
