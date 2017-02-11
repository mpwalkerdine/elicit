package elicit

import "testing"

type spec struct {
	context   *Context
	path      string
	name      string
	scenarios []*scenario
	tables    []stringTable
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
