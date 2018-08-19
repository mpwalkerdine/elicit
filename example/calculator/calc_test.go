package calculator_test

import (
	"testing"

	"github.com/mpwalkerdine/elicit"
)

func Test(t *testing.T) {
	elicit.New().
		WithSpecsFolder("./specs").
		WithTransforms(transforms).
		WithSteps(steps).
		RunTests(t)
}
