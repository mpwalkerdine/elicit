package elicit

import "testing"

type spec struct {
	context     *Context
	path        string
	name        string
	beforeSteps []*step
	scenarios   []*scenario
	afterSteps  []*step
	tables      []stringTable
}

func (s *spec) createScenario() *scenario {
	s.scenarios = append(s.scenarios, &scenario{context: s.context, spec: s})
	return s.scenarios[len(s.scenarios)-1]
}

func (s *spec) createBeforeStep() *step {
	s.beforeSteps = append(s.beforeSteps, &step{context: s.context, spec: s})
	return s.beforeSteps[len(s.beforeSteps)-1]
}

func (s *spec) createAfterStep() *step {
	s.afterSteps = append(s.afterSteps, &step{context: s.context, spec: s})
	return s.afterSteps[len(s.afterSteps)-1]
}

func (s *spec) runTest(specT *testing.T) {
	s.context.log.spec(s)
	s.context.currentSpec = s

	allSkipped := true

	for _, scenario := range s.scenarios {
		specT.Run(scenario.name, func(scenarioT *testing.T) {
			scenario.runTest(scenarioT)
			if !scenarioT.Skipped() {
				allSkipped = false
			}
		})
	}

	if allSkipped {
		specT.SkipNow()
	}
}
