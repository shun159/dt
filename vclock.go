// SPDX-License-Idenfier: BSD-2-Clause
// Author: Eishun Kondoh <dreamdiagnosis@gmail.com>

// Dotted Version Vector Sets implementation
// based on http://haslab.uminho.pt/tome/files/dvvset-dais.pdf

package dt

import (
	"sort"
	"time"

	"golang.org/x/xerrors"
)

// Timestamp of individual event
type Dot struct {
	node      string
	counter   uint32
	timestamp uint32
}

// Individual event with timestamp
type DVV struct {
	vector []Dot
}

// For sort. Returns length of vector of the dvv
func (d DVV) Len() int {
	return len(d.vector)
}

// For sort. Returns the result of compare two node
func (d DVV) Less(i, j int) bool {
	return d.vector[i].node < d.vector[j].node
}

// For sort. Swap element
func (d DVV) Swap(i, j int) {
	d.vector[i], d.vector[j] = d.vector[j], d.vector[i]
}

// Instantiate a new DVV struct
func NewDVV() DVV {
	DVV := new(DVV)
	return *DVV
}

//  Return true if the vclock va is a direct descendant of vclock vb, else false.
func (va DVV) Descends(vb DVV) bool {
	acc := true
	for _, dot_b := range vb.vector {
		if dot_a, err := va.GetDot(dot_b.node); err == nil {
			acc = (dot_a.counter >= dot_b.counter) && acc
		} else {
			acc = acc && false
		}
	}
	return acc
}

// true if vclock va strictly dominates vclock vb. Note: ignores timestamps
// In riak it is possible to have vclocks that are identical except for timestamps.
// when two vclocks descend each other, but are not equal, they are concurrent.
// See source comment for more details.
func (va DVV) Dominates(vb DVV) bool {
	// In a same world if two vclocks descend each ether they must be equal.
	// In riak they can descend each other and have different timestamps
	// How? Deleted keys, re-written, then restored is one example.
	return va.Descends(vb) && !vb.Descends(va)
}

// Combine all vclocks in the input list in to their least possible common descendant.
func (v *DVV) Merge(vclocks []DVV) {
	nclock := *v
	for len(vclocks) > 0 {
		vclock := vclocks[0]
		acc := NewDVV()
		for nclock.Len() > 0 || vclock.Len() > 0 {
			if vclock.Len() < 1 {
				acc.vector = append(acc.vector, nclock.vector...)
				break
			}
			if nclock.Len() < 1 {
				acc.vector = append(acc.vector, vclock.vector...)
				break
			} else {
				if vclock.vector[0].node < nclock.vector[0].node {
					acc.vector = append(acc.vector, vclock.vector[0])
					vclock.vector = vclock.vector[1:]
				} else if vclock.vector[0].node > nclock.vector[0].node {
					acc.vector = append(acc.vector, nclock.vector[0])
					nclock.vector = nclock.vector[1:]
				} else {
					var cnt uint32
					var ts uint32
					node := nclock.vector[0].node
					if vclock.vector[0].counter > nclock.vector[0].counter {
						cnt = vclock.vector[0].counter
						ts = vclock.vector[0].timestamp
					} else if vclock.vector[0].counter < nclock.vector[0].counter {
						cnt = nclock.vector[0].counter
						ts = nclock.vector[0].timestamp
					} else {
						if vclock.vector[0].timestamp < nclock.vector[0].timestamp {
							cnt = vclock.vector[0].counter
							ts = nclock.vector[0].timestamp
						} else {
							cnt = vclock.vector[0].counter
							ts = vclock.vector[0].timestamp
						}
					}
					vclock.vector = vclock.vector[1:]
					nclock.vector = nclock.vector[1:]
					ct := Dot{node: node, counter: cnt, timestamp: ts}
					acc.vector = append(acc.vector, ct)
				}
			}
		}
		vclocks = vclocks[1:]
		nclock = acc
	}
	v.vector = nclock.vector
}

// Get the counter value in a DVV set from node
func (v *DVV) GetCounter(id string) uint32 {
	if dot, err := v.GetDot(id); err == nil {
		return dot.counter
	} else {
		return 0
	}
}

// Get the timestamp value in a DVV set from node
func (v *DVV) GetTimestamp(id string) (uint32, error) {
	if dot, err := v.GetDot(id); err == nil {
		return dot.timestamp, nil
	} else {
		return 0, xerrors.New("Not found timestamp of id in vector")
	}
}

// Get the timestamp value in a DVV set from node
func (v *DVV) GetDot(id string) (*Dot, error) {
	for _, dot := range v.vector {
		if dot.node == id {
			return &dot, nil
		}
	}
	return nil, xerrors.New("Not found dot of id in vector")
}

// Increment vclock at Node
func (v *DVV) Increment(node string) {
	ts := uint32(time.Now().Unix())
	for idx, dot := range v.vector {
		if dot.node == node {
			v.vector[idx].counter = v.vector[idx].counter + 1
			v.vector[idx].timestamp = ts
			return
		}
	}
	dot := Dot{node: node, counter: 1, timestamp: ts}
	v.vector = append(v.vector, dot)
}

// Possibly shrink the size of a vclock, depending on current age and size.
func (v *DVV) Prune(now uint32, props map[string]int) {
	// This sort need to be deterministic, to avoid spurious merge conflicts later.
	// We achieve this by using the node ID as secondary key.
	v.sort_by_dot()
	for v.Len() > GetSmallVclock(props) {
		head_time := v.vector[0].timestamp
		is_young := (now - head_time) < uint32(GetYoungVclock(props))
		if is_young {
			return
		} else {
			is_big := v.Len() > GetBigVclock(props)
			is_old := ((now - head_time) > uint32(GetOldVclock(props)))
			if is_big || is_old {
				v.vector = v.vector[1:]
			} else {
				return
			}
		}
	}
}

// ------------------- private functions -------------------

func (v *DVV) sort_by_dot() {
	sort.Slice(v.vector, func(i, j int) bool {
		left := uint64(v.vector[i].counter)<<63 | uint64(v.vector[i].timestamp)
		right := uint64(v.vector[j].counter)<<63 | uint64(v.vector[j].timestamp)
		return left < right
	})
}
