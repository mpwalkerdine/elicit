package elicit

import "testing"
import "fmt"

type scenario struct {
	context  *Context
	spec     *spec
	name     string
	steps    []*step
	tables   []stringTable
	stepsRun int
	result   stepResult
}

func (s *scenario) runTest(scenarioT *testing.T) {
	s.result = passed

	for _, step := range s.steps {
		s.stepsRun++
		r := step.run(scenarioT)
		s.updateResult(r)
	}

	switch s.result {
	case skipped:
		scenarioT.SkipNow()
	case failed:
		scenarioT.Fail()
	}
}

func (s *scenario) updateResult(result stepResult) {
	switch result {
	case passed:
	case undefined, skipped:
		if s.result != failed {
			s.result = skipped
		}
	case failed, panicked:
		s.result = failed
	default:
		panic(fmt.Errorf("unrecognized stepResult: %d", result))
	}
}
