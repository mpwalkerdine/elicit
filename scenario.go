package elicit

import "testing"
import "fmt"

type scenario struct {
	context *Context
	spec    *spec
	name    string
	steps   []*step
	tables  []stringTable
	result  result
}

func (s *scenario) runTest(scenarioT *testing.T) result {
	s.result = passed

	if len(s.steps) == 0 {
		s.result = undefined
	}

	for _, step := range s.steps {
		r := step.run(scenarioT)
		s.updateResult(r)
	}

	return s.result
}

func (s *scenario) updateResult(result result) {
	switch result {
	case passed:
		if s.result == notrun {
			s.result = passed
		}
	case undefined, skipped:
		if s.result != failed && s.result != panicked {
			s.result = result
		}
	case failed, panicked:
		s.result = result
	default:
		panic(fmt.Errorf("unrecognized stepResult: %d", result))
	}
}
