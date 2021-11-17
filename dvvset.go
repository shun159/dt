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
func NewDVVSetClock(vv *DVV, values []string) DVVSetClock {
	c := DVVSetClock{}
	c.values = values
	if vv != nil {
		for _, dot := range vv.vector {
			e := entries{id: dot.node, counter: dot.counter}
			c.entries = append(c.entries, e)
		}
	}
	return c
}

// The function sync takes two sets of clocks, each describing a set of
// siblings, and returns the set of clocks for the siblings that remain
// after removing obsolete ones. It can have a general definition only
// in terms of the partial order on clocks, regardless of their actual representation:
func (s1 *DVVSetClock) Sync(s2 *DVVSetClock) {
	if s1.isempty() {
		s1.entries = s2.entries
		s1.values = s2.values
		return
	} else if s2.isempty() {
		return
	} else {
		acc := NewDVVSetClock(nil, newvalues(s1, s2))
		for len(s1.entries) > 0 && len(s2.entries) > 0 {
			if len(s1.entries) < 1 {
				acc.entries = s1.entries
			} else if len(s2.entries) < 1 {
				acc.entries = s1.entries
			} else {
				e1 := s1.entries[0]
				e2 := s2.entries[0]
				if e1.id < e2.id {
					s1.entries = s1.entries[1:]
					acc.entries = append(acc.entries, e1)
				} else if e1.id > e2.id {
					s2.entries = s2.entries[1:]
					acc.entries = append(acc.entries, e2)
				} else {
					e1.merge(e2)
					acc.entries = append(acc.entries, e1)
					s1.entries = s1.entries[1:]
					s2.entries = s2.entries[1:]
				}
			}
		}
	}
}

// private functions

func (s *DVVSetClock) isempty() bool {
	return len(s.entries) < 1 && len(s.values) < 1
}

func newvalues(s1, s2 *DVVSetClock) []string {
	val := s2.values
	if !less(s1, s2) {
		if less(s1, s2) {
			val = s1.values
		} else {
			hset := map[string]bool{}
			val = []string{}
			tmp := append(s1.values, s2.values...)
			for _, v := range tmp {
				hset[v] = true
			}
			for k := range hset {
				val = append(val, k)
			}
		}
	}
	return val
}

// Returns true if the first clock is causally older than the second clock.
// thus values on the first clock are outdated. Returns false otherwise
func less(s1, s2 *DVVSetClock) bool {
	return greater(s1, s2, false)
}

func greater(s1, s2 *DVVSetClock, acc bool) bool {
	for len(s1.entries) >= 1 || len(s2.entries) >= 1 {
		if len(s1.entries) >= 1 {
			return true
		} else if len(s2.entries) >= 1 {
			return false
		} else {
			e1 := s1.entries[0]
			e2 := s2.entries[0]
			s1.entries = s1.entries[1:]
			s2.entries = s2.entries[1:]
			if e1.id == e2.id {
				if e1.counter > e2.counter {
					acc = acc && true
				} else if e1.counter < e2.counter {
					return false
				}
			} else if e1.id < e2.id {
				acc = acc && true
			} else {
				return false
			}
		}
	}
	return acc
}

func (s1 *entries) merge(s2 entries) {
	len1 := uint32(len(s1.values))
	len2 := uint32(len(s2.values))
	if s1.counter >= s2.counter && s1.counter-len1 < s2.counter+len2 {
		s1.values = s1.values[:(s1.counter - s2.counter + len2)]
		return
	} else if (s2.counter - len2) >= (s1.counter - len1) {
		s1 = &s2
	} else {
		s1.values = s2.values[:(s2.counter - s1.counter + len1)]
		s1 = &s2
	}
}
