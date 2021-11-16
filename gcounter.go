// A G-counter CRDT. A G-Counter is a Grow-only counter.

package dt

import (
	"reflect"
)

type GCounter struct {
	Counters map[string]uint
}

// Create a new gcounter
func NewGCounter() GCounter {
	c := map[string]uint{}
	counter := GCounter{Counters: c}
	return counter
}

// The single total value of a gcounter
func (counter GCounter) Value() uint {
	var acc uint = 0
	for _, v := range counter.Counters {
		acc = acc + v
	}
	return acc
}

// Compare two counter for equality
func (counter_a GCounter) Equal(counter_b GCounter) bool {
	return reflect.DeepEqual(counter_a, counter_b)
}

// Combine all counters in the input list uinto a map
func (counter_a *GCounter) Merge(counters []GCounter) {
	acc := map[string]uint{}
	nodes := counter_a.AllNode(counters)
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		for j := 0; j < len(counters); j++ {
			counter_b := counters[j]
			cnt_a := counter_a.Counters[node]
			cnt_b := counter_b.Counters[node]
			cnt_max := cnt_a
			if cnt_max < cnt_b {
				cnt_max = cnt_b
			}
			acc[node] = cnt_max
		}
	}
	counter_a.Counters = acc
}

// increment counter for the node by 1
func (counter *GCounter) Increment(node string) {
	counter.IncrementBy(node, 1)
}

// perform the increment
func (counter *GCounter) IncrementBy(node string, amount uint) {
	if val, ok := counter.Counters[node]; ok {
		counter.Counters[node] = val + amount
		return
	}
	counter.Counters[node] = amount
}

// Returns the list of all nodes that have ever incremneted counter
func (counter GCounter) AllNode(counters []GCounter) []string {
	counters = append(counters, counter)
	hset := map[string]bool{}
	var nodes []string
	for i := 0; i < len(counters); i++ {
		for node := range counters[i].Counters {
			hset[node] = true
		}
	}

	for n := range hset {
		nodes = append(nodes, n)
	}

	return nodes
}
