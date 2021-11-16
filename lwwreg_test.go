// SPDX-License-Idenfier: BSD-2-Clause
// Author: Eishun Kondoh <dreamdiagnosis@gmail.com>

package dt

import (
	"testing"
)

func TestLWWRegNew(t *testing.T) {
	lww1 := NewLWWReg()
	if len(lww1.Value) != 0 {
		t.Error("LWW1.value should be empty")
	}
}

func TestLWWRegValue(t *testing.T) {
	value := "The rain in spain falls mainly on the plane"
	lww1 := NewLWWReg()
	lww1.Value = value
	lww1.Timestamp = 19090

	if lww1.Value != value {
		t.Error("LWW1.value should be same with the value")
	}

	lww2 := NewLWWReg()
	if len(lww2.Value) != 0 {
		t.Error("LWW2.value should be empty")
	}
}

func TestLWWRegAssign1(t *testing.T) {
	lww1 := NewLWWReg()
	lww1.AssignTS("value1", 2)
	lww1.AssignTS("value0", 1)
	if lww1.Value != "value1" {
		t.Error("LWW1.value should be value1")
	}
	lww1.AssignTS("value2", 3)
	if lww1.Value != "value2" {
		t.Error("LWW1.value should be value2")
	}
}

func TestLWWRegAssign2(t *testing.T) {
	lww1 := NewLWWReg()
	lww1.Assign("value2")
	lww1.Assign("value1")
	if lww1.Value != "value1" {
		t.Errorf("LWW1.value should be value1 %#v", lww1)
	}
}

func TestLWWRegMerge(t *testing.T) {
	lww1 := NewLWWReg()
	lww2 := NewLWWReg()
	lww1.Assign("old_value")
	lww2.Assign("new_value")
	lww1.Merge(lww2)
	if lww1.Value != "new_value" {
		t.Errorf("LWW1.value should be New_value %#v %#v", lww1, lww2)
	}

	lww2.Merge(lww1)
	if lww2.Value != "new_value" {
		t.Error("LWW1.value should be New_value")
	}
}

func TestLWWRegEqual(t *testing.T) {
	lww1 := LWWReg{Value: "value1", Timestamp: 1000}
	lww2 := LWWReg{Value: "value1", Timestamp: 1000}
	lww3 := LWWReg{Value: "value1", Timestamp: 1001}
	lww4 := LWWReg{Value: "value2", Timestamp: 1000}
	if lww1.Equal(lww3) {
		t.Error("LWW1 shouldn't be equal lww3")
	}

	if !lww1.Equal(lww2) {
		t.Error("LWW1 should be equal lww2")
	}

	if lww4.Equal(lww1) {
		t.Error("LWW4 should be equal lww2")
	}
}
