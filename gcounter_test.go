// SPDX-License-Idenfier: BSD-2-Clause
// Author: Eishun Kondoh <dreamdiagnosis@gmail.com>

package dt

import (
	"testing"
)

func TestNewGCounter(t *testing.T) {
	c1 := NewGCounter()
	if len(c1.Counters) != 0 {
		t.Error("C1 should be empty")
	}
}

func TestGCounterValue(t *testing.T) {
	c1 := NewGCounter()
	c2 := NewGCounter()
	c1.Counters = map[string]uint{"a": 1, "b": 13, "c": 1}

	if c1.Value() != 15 {
		t.Errorf("Counter should be 15 but %d", c1.Value())
	}

	if c2.Value() != 0 {
		t.Errorf("Counter should be 0 but %d", c2.Value())
	}
}

func TestGCounterIncrement(t *testing.T) {
	c1 := NewGCounter()

	c1.Increment("a")
	c1.Increment("b")
	c1.Increment("a")

	expected := NewGCounter()
	expected.Counters = map[string]uint{"a": 2, "b": 1}

	if !c1.Equal(expected) {
		t.Error("Counter should be same with expected")
	}
}

func TestGCounterMerge(t *testing.T) {
	c1 := NewGCounter()
	c2 := NewGCounter()
	c1.Counters = map[string]uint{"1": 1, "2": 2, "4": 4}
	c2.Counters = map[string]uint{"3": 3, "4": 3}
	c1.Merge([]GCounter{c2})

	expected := NewGCounter()
	expected.Counters = map[string]uint{"1": 1, "2": 2, "3": 3, "4": 4}

	if !c1.Equal(expected) {
		t.Errorf("Counter should be same with expected %#v", c1)
	}
}

func TestGCounterMergeLessLeft(t *testing.T) {
	c1 := NewGCounter()
	c2 := NewGCounter()
	c1.Counters = map[string]uint{"5": 5}
	c2.Counters = map[string]uint{"6": 6, "7": 7}
	c1.Merge([]GCounter{c2})

	expected := NewGCounter()
	expected.Counters = map[string]uint{"5": 5, "6": 6, "7": 7}

	if !c1.Equal(expected) {
		t.Errorf("Counter should be same with expected %#v", c1)
	}
}

func TestGCounterMergeLessRight(t *testing.T) {
	c1 := NewGCounter()
	c2 := NewGCounter()
	c1.Counters = map[string]uint{"6": 6, "7": 7}
	c2.Counters = map[string]uint{"5": 5}
	c1.Merge([]GCounter{c2})

	expected := NewGCounter()
	expected.Counters = map[string]uint{"5": 5, "6": 6, "7": 7}

	if !c1.Equal(expected) {
		t.Errorf("Counter should be same with expected %#v", c1)
	}
}

func TestGCounterUsage(t *testing.T) {
	c1 := NewGCounter()
	c2 := NewGCounter()

	if !c1.Equal(c2) {
		t.Errorf("C1 should be same with C2 %#v", c1)
	}

	c1.IncrementBy("a1", 2)
	c2.Increment("a2")
	c1.Merge([]GCounter{c2})
	c2.IncrementBy("a3", 3)
	c1.Increment("a4")
	c1.Increment("a1")
	c1.Merge([]GCounter{c2})

	expected := NewGCounter()
	expected.Counters = map[string]uint{"a1": 3, "a2": 1, "a3": 3, "a4": 1}

	if !c1.Equal(expected) {
		t.Errorf("Counter should be same with expected %#v", c1)
	}
}
