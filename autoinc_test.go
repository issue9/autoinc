// SPDX-FileCopyrightText: 2015-2024 caixw
//
// SPDX-License-Identifier: MIT

package autoinc

import (
	"context"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/issue9/assert/v4"
)

func TestAutoInc_overflow(t *testing.T) {
	a := assert.New(t, false)

	ai := New(math.MaxInt64-1, 2, 4)
	a.NotNil(ai).
		Go(func(a *assert.Assertion) {
			a.ErrorIs(ai.Serve(context.Background()), ErrOverflow())
		}).
		Wait(500 * time.Microsecond) // 保证 ai.Serve 执行完成

	id, ok := ai.ID()
	a.True(ok).Equal(math.MaxInt64-1, id)

	// 不存在第二条数据
	id, ok = ai.ID()
	a.False(ok).Equal(id, 0)

	id, ok = ai.ID()
	a.False(ok).Equal(id, 0)

	a.Panic(func() {
		id = ai.MustID()
		a.Equal(id, 0)
	})
}

func TestAutoInc_ID_1(t *testing.T) {
	a := assert.New(t, false)

	a.Panic(func() {
		ai := New(0, 0, 2)
		a.Nil(ai)
	})

	// 正规的 ai 操作
	ctx, cancel := context.WithCancel(context.Background())
	ai := New(0, 2, 2)
	a.NotNil(ai).
		Go(func(*assert.Assertion) { ai.Serve(ctx) }).
		Wait(500 * time.Microsecond) // 保证 ai.Serve 执行完成
	for i := 0; i < 7; i++ {
		a.Equal(ai.MustID(), i*2)
	}
	cancel()

	// 可以从负数起始
	ctx, cancel = context.WithCancel(context.Background())
	ai = New(-100, 2, 5)
	a.NotNil(ai).
		Go(func(*assert.Assertion) { ai.Serve(ctx) }).
		Wait(500 * time.Microsecond) // 保证 ai.Serve 执行完成
	for i := 0; i < 7; i++ {
		a.Equal(ai.MustID(), -100+i*2)
	}
	cancel()

	// start,step 双负数
	ctx, cancel = context.WithCancel(context.Background())
	ai = New(-100, -3, 0)
	a.NotNil(ai).
		Go(func(a *assert.Assertion) { ai.Serve(ctx) }).
		Wait(500 * time.Microsecond) // 保证 ai.Serve 执行完成
	for i := 0; i < 7; i++ {
		a.Equal(ai.MustID(), -100+i*-3)
	}
	cancel()
}

func TestAutoInc_ID_2(t *testing.T) {
	a := assert.New(t, false)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ai := New(2, 2, 2)
	a.NotNil(ai).
		Go(func(a *assert.Assertion) { ai.Serve(ctx) }).
		Wait(500 * time.Microsecond) // 保证 ai.Serve 执行完成

	mu := sync.Mutex{}
	mapped := map[int64]bool{}
	wg := &sync.WaitGroup{}

	fn := func() {
		for i := 0; i < 100; i++ {
			id := ai.MustID()

			mu.Lock()
			_, found := mapped[id]
			a.False(found, "找到重复元素:%v", id)
			mapped[id] = true
			mu.Unlock()
		}
		wg.Done()
	}

	wg.Add(4)
	go fn()
	go fn()
	go fn()
	go fn()
	wg.Wait()
}
