package fibonacci_test

import (
	"mmatt/elicit"
	"mmatt/elicit/examples/fibonacci/steps"
	"testing"
)

func Test(t *testing.T) {
	elicit.New().
		WithSpecsFolder("./specs").
		WithSteps(steps.Get()).
		RunTests(t)
}
