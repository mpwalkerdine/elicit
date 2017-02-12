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
	notrun result = iota
	passed
	skipped
	undefined
	failed
	panicked
)

func (r result) shouldLog() bool {
	switch r {
	case skipped, passed:
		return false
	case undefined, failed, panicked:
		return true
	default:
		panic(fmt.Errorf("unknown result: %d", r))
	}
}

func (r result) String() string {
	switch r {
	case undefined:
		return "PENDING"
	case skipped:
		return "SKIP"
	case failed:
		return "FAIL"
	case panicked:
		return "PANIC"
	case passed:
		return "PASS"
	default:
		panic(fmt.Errorf("unknown result: %d", r))
	}
}

type results []result

func (r results) Len() int           { return len(r) }
func (r results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r results) Less(i, j int) bool { return r[i] < r[j] }

func (s *spec) runTest(specT *testing.T) {
	allSkipped := true

	for _, scenario := range s.scenarios {
		specT.Run(scenario.name, func(scenarioT *testing.T) {
			scenario.runTest(scenarioT)

			switch scenario.result {
			case skipped, undefined:
				scenarioT.SkipNow()
			case failed, panicked:
				allSkipped = false
				scenarioT.Fail()
			default:
				allSkipped = false
			}
		})
		if scenario.result > s.result {
			s.result = scenario.result
		}
	}

	if allSkipped {
		specT.SkipNow()
	}
}
