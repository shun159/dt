// A convergent, replicated, state based grow only set

package dt

import (
	"reflect"
)

type GSet struct {
	Set map[string]bool
}

// Create a new gset
func NewGSet() GSet {
	s := map[string]bool{}
	return GSet{Set: s}
}

// Returns a set of value.
func (gset *GSet) Values() []string {
	acc := []string{}
	for v, _ := range gset.Set {
		acc = append(acc, v)
	}
	return acc
}

// append an element to the set
func (gset *GSet) Add(elem string) {
	gset.Set[elem] = true
}

// return true if the element exists in the set
func (gset GSet) Exists(elem string) bool {
	if _, ok := gset.Set[elem]; ok {
		return true
	}
	return false
}

// Compair Two set for equality
func (gset_a GSet) Equal(gset_b GSet) bool {
	return reflect.DeepEqual(gset_a.Set, gset_b.Set)
}
