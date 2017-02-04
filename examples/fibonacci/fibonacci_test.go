package fibonacci_test

import (
	"mmatt/elicit"
	"mmatt/elicit/examples/fibonacci"
	"testing"
)

func Test(t *testing.T) {
	elicit.New().
		WithSpecsFolder("./specs").
		WithSteps(steps).
		RunTests(t)
}

// TODO(matt) find a way of getting a reflected version of a function from its name at runtime
var steps = map[string]interface{}{
	`Step to run (before|after) each scenario`:                                                 stepToRunPhaseEachScenario,
	`The (-?\d+)(?:st|nd|rd|th) item is (-?\d+)`:                                               theNthItemIsX,
	`The first (\d+) items are ((?:\d+,\s*)+\d+)`:                                              theFirstNItemsAreX,
	`The (-?\d+)(?:st|nd|rd|th) to the (-?\d+)(?:st|nd|rd|th) items are ((?:-?\d+,\s*)+-?\d+)`: theFoothToTheBarthItemsAreBaz,
	`This step takes a table parameter`:                                                        thisStepTakesATableParameter,
	`The (.*) and the (.*)`:                                                                    stepWithMismatchedPattern,
}

func stepToRunPhaseEachScenario(e *elicit.Context, phase string) {
	e.Assert(phase).IsIn("before", "after")
}

func theNthItemIsX(e *elicit.Context, n, x int) {
	r := fibonacci.Sequence(n, n)
	e.Assert(r[0]).IsDeepEqual(x)
}

func theFirstNItemsAreX(e *elicit.Context, n int, x []int) {
	r := fibonacci.Sequence(1, n)
	e.Assert(r).IsDeepEqual(x)
}

func theFoothToTheBarthItemsAreBaz(e *elicit.Context, m, n int, x []int) {
	r := fibonacci.Sequence(m, n)
	e.Assert(r).IsDeepEqual(x)
}

func thisStepTakesATableParameter(e *elicit.Context, table elicit.Table) {
	e.Assert(table.Columns).IsNotEmpty()
	e.Assert(table.Rows).IsNotEmpty()
}

func stepWithMismatchedPattern(e *elicit.Context, a string) {

}
