// SPDX-License-Idenfier: BSD-2-Clause
// Author: Eishun Kondoh <dreamdiagnosis@gmail.com>

package dt

import (
	"testing"
)

func TestDVVSetNew(t *testing.T) {
	c1 := NewDVVSetClock(nil, []string{})
	if !(len(c1.entries) < 1) {
		t.Error("Initial dvvset must be have empty set")
	}

	c2 := NewDVVSetClock(nil, []string{"a", "b", "c"})
	if len(c2.values) != 3 {
		t.Error("Initial dvvset must have empty set with 3 values")
	}

	if len(c2.entries) != 0 {
		t.Error("Initial dvvset must have empty entries")
	}

	v1 := NewDVV()
	v1.vector = []Dot{Dot{node: "1", counter: 1, timestamp: 1}}
	c3 := NewDVVSetClock(&v1, []string{"a", "b", "c"})

	if len(c3.entries) != 1 {
		t.Error("Initial dvvset must have a entry")
	}

	if len(c3.values) != 3 {
		t.Error("Initial dvvset must have a entry")
	}
}
