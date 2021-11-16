// SPDX-License-Idenfier: BSD-2-Clause
// Author: Eishun Kondoh <dreamdiagnosis@gmail.com>

package dt

import (
	"testing"
)

func TestNewGSet(t *testing.T) {
	gset := NewGSet()
	if len(gset.Values()) != 0 {
		t.Errorf("GSet.value should be empty %#v", gset)
	}
}

func TestGSetAdd(t *testing.T) {
	gset := NewGSet()
	gset.Add("value1")
	gset.Add("value2")
	if gset.Values()[1] != "value2" {
		t.Errorf("the 1st element of the gset.value should be 'value2'")
	}
}

func TestGSetExists(t *testing.T) {
	gset := NewGSet()
	gset.Add("value1")
	gset.Add("value2")
	if !gset.Exists("value2") {
		t.Errorf("gset should have 'value2'")
	}

	if gset.Exists("value3") {
		t.Errorf("gset shouldn't have 'value3'")
	}
}

func TestGSetEqual(t *testing.T) {
	gset1 := NewGSet()
	gset1.Add("value1")

	gset2 := NewGSet()
	gset2.Add("value1")

	if !gset1.Equal(gset2) {
		t.Errorf("gset should be equal with gset2")
	}

	gset2.Add("value2")

	if gset1.Equal(gset2) {
		t.Errorf("gset shouldn't be equal with gset2")
	}
}
