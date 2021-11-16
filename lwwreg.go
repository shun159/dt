// A state-based last-writer-wins register

package dt

import (
	"time"
)

type LWWReg struct {
	Value     string
	Timestamp int64
}

// Create a new empty lwwreg
func NewLWWReg() LWWReg {
	return LWWReg{}
}

// Assign a value to the lwwreg associating the update with time ts
func (reg *LWWReg) Assign(value string) {
	reg.AssignTS(value, time.Now().UnixNano())
}

func (reg *LWWReg) AssignTS(value string, ts int64) {
	if ts > reg.Timestamp {
		reg.Value = value
		reg.Timestamp = ts
	}
}

// Merge two lwwreg to a single lwwreg. this is the least upper bound
// function described in the literature
func (reg_a *LWWReg) Merge(reg_b LWWReg) {
	if reg_b.Timestamp > reg_a.Timestamp {
		reg_a.Timestamp = reg_b.Timestamp
		reg_a.Value = reg_b.Value
	} else {
		return
	}
}

// Are two lwwreg s structurally equal? this is not value equality.
// Two regsiters might represent the value armchair and not be equal().
// Equality here is that both registers contain the same value and timestamp
func (reg_a LWWReg) Equal(reg_b LWWReg) bool {
	return reg_a.Value == reg_b.Value && reg_a.Timestamp == reg_b.Timestamp
}
