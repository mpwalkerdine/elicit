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
	switch r {
	case skipped, passed:
		return false
	case pending, failed, panicked:
		return true
	default:
		panic(fmt.Errorf("unknown result: %d", r))
	}
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

func (s *spec) runTest(specT *testing.T) {
	for _, scenario := range s.scenarios {

		for _, h := range s.context.beforeScenario {
			h()
		}

		specT.Run(scenario.name, func(scenarioT *testing.T) {
			scenario.run(scenarioT)
		})

		for _, h := range s.context.afterScenario {
			h()
		}

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
