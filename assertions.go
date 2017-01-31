package elicit

import (
	"reflect"
)

// Assertion wraps an actual value for fluent-style assertions against expected ones
type Assertion struct {
	context *Context
	actual  interface{}
}

const (
	notEqualFmt  = "expected %v, got %v"
	notInListFmt = "expected one of %v, got %v"
)

// Assert takes an actual value and returns an object against which assertions can be made
// in a fluent-style syntax
func (e *Context) Assert(actual interface{}) *Assertion {
	return &Assertion{context: e, actual: actual}
}

// IsTrue checks that the actual value is a bool which is true
func (a *Assertion) IsTrue() {
	if b, isBool := a.actual.(bool); !isBool || !b {
		a.context.Failf(notEqualFmt, true, a.actual)
	}
}

// IsDeepEqual checks that the reflect.DeepEqual(actual, expected)
func (a *Assertion) IsDeepEqual(expected interface{}) {
	if !reflect.DeepEqual(a.actual, expected) {
		a.context.Failf(notEqualFmt, expected, a.actual)
	}
}

// IsIn checks that the actual value is contained in the list of supplied values
func (a *Assertion) IsIn(list ...interface{}) {
	for _, e := range list {
		if reflect.DeepEqual(a.actual, e) {
			return
		}
	}
	a.context.Failf(notInListFmt, list, a.actual)
}
