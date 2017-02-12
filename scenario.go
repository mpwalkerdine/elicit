package elicit

import "testing"

type scenario struct {
	context *Context
	spec    *spec
	name    string
	steps   []*step
	tables  []stringTable
	result  result
}

func (s *scenario) runTest(scenarioT *testing.T) {
	s.result = passed

	if len(s.steps) == 0 {
		s.result = undefined
	}

	for _, step := range s.steps {
		step.run(scenarioT)
		if step.result > s.result {
			s.result = step.result
		}
	}
}
