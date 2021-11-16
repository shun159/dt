// SPDX-License-Idenfier: BSD-2-Clause
// Author: Eishun Kondoh <dreamdiagnosis@gmail.com>

package dt

// sorted by id
type entries struct {
	id      string
	counter uint32
	values  []string
}

type DVVSetClock struct {
	entries []entries
	values  []string
}

// Construct a new clock set without causal history
// and receives a list of values that goes to the anonymous list.
func DVVSetClockNew() DVVSetClock {
	return DVVSetClock{}
}

func DVVSetClockNewWithValue(values []string) DVVSetClock {
	c := DVVSetClock{values: values}
	return c
}

// Construct a new clock set with causal history
// of the given version vector clock, and receives a list of values that goes to the anonymous list
// the version vector clock should be a direct result of Join
func DVVSetClockNewWithVV(vv DVV, values []string) DVVSetClock {
	c := DVVSetClockNew()
	c.values = values
	for _, dot := range vv.vector {
		e := entries{id: dot.node, counter: dot.counter}
		c.entries = append(c.entries, e)
	}
	return c
}
