package elicit

import (
	"fmt"
	"testing"
)

type spec struct {
	context   *Context
	path      string
	name      string
	scenarios []*scenario
	tables    []stringTable
	result    result
}

type result int

const (
	// Note: These are ordered by precedence
	passed result = iota
	skipped
	pending
	failed
	panicked
	numResultTypes
)

func (r result) shouldLog() bool {
	return r > skipped
}

func (r result) String() string {
	switch r {
	case pending:
		return "Pending"
	case skipped:
		return "Skipped"
	case failed:
		return "Failed"
	case panicked:
		return "Panicked"
	case passed:
		return "Passed"
	default:
		panic(fmt.Errorf("unknown result: %d", r))
	}
}

func (s *spec) run(specT *testing.T) {
	for _, scenario := range s.scenarios {

		var hookErr error
		if hookErr = s.context.beforeScenario.run("before scenario"); hookErr != nil {
			scenario.result = panicked
		}

		specT.Run(scenario.name, func(scenarioT *testing.T) {
			if hookErr != nil {
				scenarioT.FailNow()
			}

			s.runScenario(scenarioT, scenario)
		})

		if scenario.result > s.result {
			s.result = scenario.result
		}
	}

	switch s.result {
	case panicked, failed:
		specT.Fail()
	case skipped, pending:
		specT.SkipNow()
	}
}

func (s *spec) runScenario(scenarioT *testing.T, scenario *scenario) {
	// Ensure the after scenario hooks are run regardless of the result
	defer func() {
		if hookErr := s.context.afterScenario.run("after scenario"); hookErr != nil {
			scenario.result = panicked
			scenarioT.FailNow()
		}
	}()

	scenario.run(scenarioT)
}

func (s *spec) skipAllScenarios() {
	for _, scenario := range s.scenarios {
		scenario.result = skipped
	}
}
