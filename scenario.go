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

		if s.result != passed {
			scenarioT.SkipNow()
		}

		if hookErr := s.context.beforeStep.run("before step"); hookErr != nil {
			s.result = panicked
			step.result = panicked
			scenarioT.Fail()
			return
		}

		s.runStep(scenarioT, step)
	}
}
func (s *scenario) runStep(scenarioT *testing.T, step *step) {
	// Ensure after step hooks execute regardless of what happens in the step
	defer func() {
		if step.result > s.result {
			s.result = step.result
		}

		if hookErr := s.context.afterStep.run("after step"); hookErr != nil {
			s.result = panicked
			step.result = panicked
			scenarioT.Fail()
		}
	}()

	step.run(scenarioT)
}
