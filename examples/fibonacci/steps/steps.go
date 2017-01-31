package steps

import (
	"mmatt/elicit"
	"mmatt/elicit/examples/fibonacci"
)

// SpuriousType should not affect generation
type SpuriousType struct{}

// StepToRunPhaseEachScenario .
//elicit:step Step to run (before|after) each scenario
func StepToRunPhaseEachScenario(e *elicit.Context, phase string) {
	e.Assert(phase == "before" || phase == "after").IsTrue()
}

// TheNthItemIsX .
//elicit:step The (-?\d+)(?:st|nd|rd|th) item is (-?\d+)
func TheNthItemIsX(e *elicit.Context, n, x int) {
	r := fibonacci.Sequence(n, n)
	e.Assert(r[0]).IsDeepEqual(x)
}

// TheFirstNItemsAreX .
//elicit:step The first (\d+) items are ((?:\d+,\s*)+\d+)
func TheFirstNItemsAreX(e *elicit.Context, n int, x []int) {
	r := fibonacci.Sequence(1, n)
	e.Assert(r).IsDeepEqual(x)
}

// TheFoothToTheBarthItemsAreBaz .
//elicit:step The (\d+)(?:st|nd|rd|th) to the (\d+)(?:st|nd|rd|th) items are ((?:\d+,\s*)+\d+)
func TheFoothToTheBarthItemsAreBaz(e *elicit.Context, m, n int, x []int) {
	r := fibonacci.Sequence(m, n)
	e.Assert(r).IsDeepEqual(x)
}

// StepWithNoPattern should not be registered
func StepWithNoPattern(e *elicit.Context, a string) {

}

// StepWithMismatchedPattern should cause error during generation
//elicit:step The (.*) and the (.*)
func StepWithMismatchedPattern(e *elicit.Context, a string) {

}

// StepButWithReturn should not be registered
func StepButWithReturn(e *elicit.Context) string {
	return "nothing to see here"
}

// SpuriousMethod should not be registered
func (st *SpuriousType) SpuriousMethod(e *elicit.Context) {

}

// SpuriousMethodReturn should not be registered
func (st *SpuriousType) SpuriousMethodReturn(e *elicit.Context) int {
	return 1
}

// SpuriousFunc should not be registered
func SpuriousFunc(s string) {}
