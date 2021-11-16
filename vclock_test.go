// SPDX-License-Idenfier: BSD-2-Clause
// Author: Eishun Kondoh <dreamdiagnosis@gmail.com>

package dt

import (
	"reflect"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	a := NewDVV()
	b := NewDVV()

	a.Increment("a")
	b.Increment("b")

	if !a.Descends(a) {
		t.Error("A should be descendant of self")
	}

	if !b.Descends(b) {
		t.Error("B should be descendant of self")
	}

	if a.Descends(b) {
		t.Error("A shouldn't be descendant of B")
	}

	a.Increment("a")
	a.Merge([]DVV{b})
	a.Increment("c")

	if !a.Descends(a) {
		t.Error("A should be descendant of self")
	}

	if !a.Descends(b) {
		t.Error("A should be descendant of self")
	}

	if b.Descends(a) {
		t.Error("B shouldn't be descendant of A")
	}
}

func TestDVVSmallPrune(t *testing.T) {
	now_ts := uint32(time.Now().Unix())
	old_ts := now_ts - 32000000
	props := map[string]int{"small_vclock": 4}
	vclock := NewDVV()
	vclock.vector = []Dot{
		Dot{node: "1",
			counter:   1,
			timestamp: old_ts},
		Dot{node: "2",
			counter:   1,
			timestamp: old_ts},
		Dot{node: "3",
			counter:   1,
			timestamp: old_ts},
	}
	vclock.Prune(now_ts, props)
	if !(vclock.Len() == 3) {
		t.Error("vclock with less entries small_vclocks will be untouched")
	}
}

func TestDVVYoungPrune(t *testing.T) {
	now_ts := uint32(time.Now().Unix())
	new_ts := now_ts - 1
	props := map[string]int{"small_vclock": 1, "young_vclock": 1000}
	vclock := NewDVV()
	vclock.vector = []Dot{
		Dot{node: "1",
			counter:   1,
			timestamp: new_ts},
		Dot{node: "2",
			counter:   1,
			timestamp: new_ts},
		Dot{node: "3",
			counter:   1,
			timestamp: new_ts},
	}
	vclock.Prune(now_ts, props)
	if !(vclock.Len() == 3) {
		t.Error("vclock with less entries small_vclocks will be untouched")
	}
}

func TestDVVBigPrune(t *testing.T) {
	// vclock not preserved by small or young will be pruned down to
	// no larger than big_vclock entries
	now_ts := uint32(time.Now().Unix())
	new_ts := now_ts - 1000
	props := map[string]int{
		"small_vclock": 1,
		"young_vclock": 1,
		"big_vclock":   2,
		"old_vclock":   100000,
	}
	vclock := NewDVV()
	vclock.vector = []Dot{
		Dot{node: "1",
			counter:   1,
			timestamp: new_ts},
		Dot{node: "2",
			counter:   1,
			timestamp: new_ts},
		Dot{node: "3",
			counter:   1,
			timestamp: new_ts},
	}
	vclock.Prune(now_ts, props)
	if !(vclock.Len() == 2) {
		t.Error("vclock should be 2")
	}
}

func TestDVVOldPrune(t *testing.T) {
	// vclock not preserved by small or young will be pruned down to
	// no larger than big_vclock and no entries more than old_vclock ago
	now_ts := uint32(time.Now().Unix())
	new_ts := now_ts - 1000
	old_ts := now_ts - 100000
	props := map[string]int{
		"small_vclock": 1,
		"young_vclock": 1,
		"big_vclock":   2,
		"old_vclock":   10000,
	}
	vclock := NewDVV()
	vclock.vector = []Dot{
		Dot{node: "1",
			counter:   1,
			timestamp: new_ts},
		Dot{node: "2",
			counter:   1,
			timestamp: old_ts},
		Dot{node: "3",
			counter:   1,
			timestamp: old_ts},
	}
	vclock.Prune(now_ts, props)
	if !(vclock.Len() == 1) {
		t.Error("vclock should be 1")
	}
}

func TestDVVPruneOrder(t *testing.T) {
	// vclock with two nodes of the same timestamp will be pruned down
	// to the same node
	now_ts := uint32(time.Now().Unix())
	old_ts := now_ts - 100000
	props := map[string]int{
		"small_vclock": 1,
		"young_vclock": 1,
		"big_vclock":   2,
		"old_vclock":   10000,
	}
	vclock1 := NewDVV()
	vclock1.vector = []Dot{
		Dot{node: "1",
			counter:   1,
			timestamp: old_ts},
		Dot{node: "2",
			counter:   2,
			timestamp: old_ts},
	}

	vclock2 := NewDVV()
	vclock2.vector = []Dot{
		Dot{node: "2",
			counter:   2,
			timestamp: old_ts},
		Dot{node: "1",
			counter:   1,
			timestamp: old_ts},
	}

	vclock1.Prune(now_ts, props)
	vclock2.Prune(now_ts, props)
	if !(vclock1.vector[0] == vclock2.vector[0]) {
		t.Error("vclock should be same")
	}
}

func TestDVVAccessors(t *testing.T) {
	vclock := NewDVV()
	vclock.vector = []Dot{
		Dot{node: "1",
			counter:   1,
			timestamp: 1},
		Dot{node: "2",
			counter:   2,
			timestamp: 2},
	}

	if !(vclock.GetCounter("1") == 1) {
		t.Error("GetCounter should return 1")
	}
	if !(vclock.GetCounter("2") == 2) {
		t.Error("GetCounter should return 2")
	}
	if !(vclock.GetCounter("3") == 0) {
		t.Error("GetCounter should return 0")
	}
	if val, _ := vclock.GetTimestamp("1"); !(val == 1) {
		t.Error("GetTimestamp should return 1")
	}
	if val, _ := vclock.GetTimestamp("2"); !(val == 2) {
		t.Error("GetTimestamp should return 2")
	}
	if val, _ := vclock.GetTimestamp("3"); !(val == 0) {
		t.Error("GetTimestamp should return 0")
	}
}

func TestDVVMerge(t *testing.T) {
	vc1 := NewDVV()
	vc1.vector = []Dot{
		Dot{node: "1",
			counter:   1,
			timestamp: 1},
		Dot{node: "2",
			counter:   2,
			timestamp: 2},
		Dot{node: "4",
			counter:   4,
			timestamp: 4},
	}

	vc2 := NewDVV()
	vc2.vector = []Dot{
		Dot{node: "3",
			counter:   3,
			timestamp: 3},
		Dot{node: "4",
			counter:   3,
			timestamp: 3},
	}
	vc1.Merge([]DVV{vc2})

	expect := NewDVV()
	expect.vector = []Dot{
		Dot{node: "1",
			counter:   1,
			timestamp: 1},
		Dot{node: "2",
			counter:   2,
			timestamp: 2},
		Dot{node: "3",
			counter:   3,
			timestamp: 3},
		Dot{node: "4",
			counter:   4,
			timestamp: 4},
	}
	if !(reflect.DeepEqual(vc1, expect)) {
		t.Error("vc1 should be merged as expected array")
	}
}

func TestDVVLeftMerge(t *testing.T) {
	vc1 := NewDVV()
	vc1.vector = []Dot{
		Dot{node: "5",
			counter:   5,
			timestamp: 5},
	}

	vc2 := NewDVV()
	vc2.vector = []Dot{
		Dot{node: "6",
			counter:   6,
			timestamp: 6},
		Dot{node: "7",
			counter:   7,
			timestamp: 7},
	}
	vc1.Merge([]DVV{vc2})

	expect := NewDVV()
	expect.vector = []Dot{
		Dot{node: "5",
			counter:   5,
			timestamp: 5},
		Dot{node: "6",
			counter:   6,
			timestamp: 6},
		Dot{node: "7",
			counter:   7,
			timestamp: 7},
	}
	if !(reflect.DeepEqual(vc1, expect)) {
		t.Error("vc1 should be merged as expected array")
	}
}

func TestDVVRightMerge(t *testing.T) {
	vc1 := NewDVV()
	vc1.vector = []Dot{
		Dot{node: "6",
			counter:   6,
			timestamp: 6},
		Dot{node: "7",
			counter:   7,
			timestamp: 7},
	}

	vc2 := NewDVV()
	vc2.vector = []Dot{
		Dot{node: "5",
			counter:   5,
			timestamp: 5},
	}

	vc1.Merge([]DVV{vc2})

	expect := NewDVV()
	expect.vector = []Dot{
		Dot{node: "5",
			counter:   5,
			timestamp: 5},
		Dot{node: "6",
			counter:   6,
			timestamp: 6},
		Dot{node: "7",
			counter:   7,
			timestamp: 7},
	}
	if !(reflect.DeepEqual(vc1, expect)) {
		t.Error("vc1 should be merged as expected array")
	}
}

func TestDVVEqualMerge(t *testing.T) {
	vc1 := NewDVV()
	vc1.vector = []Dot{
		Dot{node: "1",
			counter:   1,
			timestamp: 2},
		Dot{node: "2",
			counter:   1,
			timestamp: 4},
	}

	vc2 := NewDVV()
	vc2.vector = []Dot{
		Dot{node: "1",
			counter:   1,
			timestamp: 3},
		Dot{node: "3",
			counter:   1,
			timestamp: 5},
	}

	vc1.Merge([]DVV{vc2})

	expect := NewDVV()
	expect.vector = []Dot{
		Dot{node: "1",
			counter:   1,
			timestamp: 3},
		Dot{node: "2",
			counter:   1,
			timestamp: 4},
		Dot{node: "3",
			counter:   1,
			timestamp: 5},
	}
	if !(reflect.DeepEqual(vc1, expect)) {
		t.Error("vc1 should be merged as expected array")
	}
}

func TestDVVGetEntry(t *testing.T) {
	vc := NewDVV()
	vc.Increment("a")
	vc.Increment("b")
	vc.Increment("c")
	vc.Increment("a")
	if val, _ := vc.GetDot("a"); !(val.counter == 2) {
		t.Error("should be 2")
	}

	if val, _ := vc.GetDot("b"); !(val.counter == 1) {
		t.Error("should be 1")
	}

	if val, _ := vc.GetDot("c"); !(val.counter == 1) {
		t.Error("should be 1")
	}
}
