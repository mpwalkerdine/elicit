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

func (s *scenario) run(scenarioT *testing.T) {
	if len(s.steps) == 0 {
		s.result = pending
	}

	for _, step := range s.steps {
		for _, h := range s.context.beforeStep {
			h()
		}

		step.run(scenarioT)

		for _, h := range s.context.afterStep {
			h()
		}
	}
}
