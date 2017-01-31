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
	failMsgFmt = "expected %v, got %v"
)

// Assert takes an actual value and returns an object against which assertions can be made
// in a fluent-style syntax
func (e *Context) Assert(actual interface{}) *Assertion {
	return &Assertion{context: e, actual: actual}
}

// IsTrue checks that the actual value is a bool which is true
func (a *Assertion) IsTrue() {
	if b, isBool := a.actual.(bool); !isBool || !b {
		a.context.Failf(failMsgFmt, true, a.actual)
	}
}

// IsDeepEqual checks that the reflect.DeepEqual(actual, expected)
func (a *Assertion) IsDeepEqual(expected interface{}) {
	if !reflect.DeepEqual(a.actual, expected) {
		a.context.Failf(failMsgFmt, expected, a.actual)
	}
}
