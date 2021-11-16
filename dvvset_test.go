// SPDX-License-Idenfier: BSD-2-Clause
// Author: Eishun Kondoh <dreamdiagnosis@gmail.com>

package dt

import (
	"testing"
)

func TestDVVSetNew(t *testing.T) {
	c1 := DVVSetClockNew()
	if !(len(c1.entries) < 1) {
		t.Error("Initial dvvset must be have empty set")
	}

	c2 := DVVSetClockNewWithValue([]string{"a", "b", "c"})
	if len(c2.values) != 3 {
		t.Error("Initial dvvset must be have empty set with 3 values")
	}

	v1 := NewDVV()
	v1.vector = []Dot{
		Dot{
			node:      "1",
			counter:   1,
			timestamp: 1,
		},
	}
	c3 := DVVSetClockNewWithVV(v1, []string{"a", "b", "c"})
	if len(c3.entries) != 1 {
		t.Error("Initial dvvset must be have a entry")
	}
}
