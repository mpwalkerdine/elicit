package elicit

import "testing"
import "fmt"
import "os"

// CurrentContext is a single instance used by all elicit tests
var CurrentContext = &Context{}

// Context stores test machinery and maintains state between specs/scenarios/steps
type Context struct {
	specT     *testing.T
	scenarioT *testing.T
}

// BeginSpecTest registers the start of Spec
func (e *Context) BeginSpecTest(t *testing.T) {
	e.specT = t
}

// BeginScenarioTest registers the start of a Scenario
func (e *Context) BeginScenarioTest(t *testing.T) {
	e.scenarioT = t
}

// RegisterStep maps a Regexpr to a step implementation
func (e *Context) RegisterStep(pattern string, stepFunc interface{}) {
	fmt.Fprintf(os.Stderr, "++ Registering step %T for %q\n", stepFunc, pattern)
}

// RunStep matches the stepText to a registered step implementation and invokes it
func (e *Context) RunStep(stepText string) {
	e.scenarioT.Logf(stepText)
}

// Fail records test failure
func (e *Context) Fail() {
	e.scenarioT.Fail()
}

// Assert that the parameter is true, otherwise fails
func (e *Context) Assert(shouldBeTrue bool) {
	if !shouldBeTrue {
		e.Fail()
	}
}
