// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package autoinc

import (
	"testing"

	"github.com/issue9/assert"
)

func TestAutoInc1(t *testing.T) {
	a := assert.New(t)

	// 正规的ai操作
	ai := New(0, 2, 2)
	a.NotNil(ai)
	for i := 0; i < 7; i++ {
		a.Equal(ai.ID(), i*2)
	}

	// 可以从负数起始
	ai = New(-100, 2, 5)
	a.NotNil(ai)
	for i := 0; i < 7; i++ {
		a.Equal(ai.ID(), -100+i*2)
	}

	// start,step双负数
	ai = New(-100, -3, 0)
	a.NotNil(ai)
	for i := 0; i < 7; i++ {
		a.Equal(ai.ID(), -100+i*-3)
	}
}

func TestAutoInc2(t *testing.T) {
	a := assert.New(t)

	ai := New(2, 2, 2)
	a.NotNil(ai)
	mapped := map[int64]bool{}

	fn := func() {
		for i := 0; i < 100; i++ {
			id := ai.ID()
			_, found := mapped[id]
			a.False(found, "找到重复元素:%v", id)
			mapped[id] = true
		}
	}

	go fn()
	go fn()
}
