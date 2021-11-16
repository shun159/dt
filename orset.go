// A convergent, replicated, state based observe remove set

package dt

import (
	"errors"
	"math"
	"math/rand"
)

type Token struct {
	Tag       int
	IsRemoved bool
}

type ORSet struct {
	Set map[string]Token
}

// Create a new ORSet
func NewOrset() ORSet {
	s := map[string]Token{}
	return ORSet{Set: s}
}

// Return values which hasn't removed only
func (orset ORSet) Value() []string {
	acc := []string{}
	for value, token := range orset.Set {
		if !token.IsRemoved {
			acc = append(acc, value)
		}
	}
	return acc
}

// Return values which has removed
func (orset ORSet) RemovedValue() []string {
	acc := []string{}
	for value, token := range orset.Set {
		if token.IsRemoved {
			acc = append(acc, value)
		}
	}
	return acc
}

func (orset ORSet) Token(elem string) (Token, error) {
	if token, ok := orset.Set[elem]; ok {
		return token, nil
	}
	return Token{}, errors.New("the elem doesn't exist")
}

func (orset ORSet) Lookup(elem string) bool {
	if _, ok := orset.Set[elem]; ok {
		return true
	}
	return false
}

// private functions

func unique() int {
	s := rand.NewSource(int64(rand.Int63()))
	r := rand.New(s)
	return r.Intn(math.MaxInt64)
}
