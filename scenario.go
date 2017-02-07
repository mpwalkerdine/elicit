package elicit

import "testing"
import "fmt"

type scenario struct {
	context  *Context
	spec     *spec
	name     string
	steps    []step
	tables   []stringTable
	stepsRun int
	result   stepResult
}

func (s *scenario) createStep() *step {
	s.steps = append(s.steps, step{context: s.context, spec: s.spec, scenario: s})
	return &s.steps[len(s.steps)-1]
}

func (s *scenario) runTest(scenarioT *testing.T) {
	s.context.log.scenario(s)
	s.context.currentScenario = s
	s.result = passed

	for _, before := range s.spec.beforeSteps {
		before.scenario = s
		before.run(scenarioT)
		before.scenario = nil
	}

	for _, step := range s.steps {
		step.run(scenarioT)
	}

	for _, after := range s.spec.afterSteps {
		after.scenario = s
		after.run(scenarioT)
		after.scenario = nil
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
