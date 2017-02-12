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
	notrun result = iota
	undefined
	skipped
	failed
	panicked
	passed
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
		r := notrun
		specT.Run(scenario.name, func(scenarioT *testing.T) {
			r = scenario.runTest(scenarioT)

			if r != skipped {
				allSkipped = false
			}

			if r == skipped {
				scenarioT.SkipNow()
			} else if r == failed {
				scenarioT.Fail()
			}
		})
		s.updateResult(r)
	}

	if allSkipped {
		specT.SkipNow()
	}
}

func (s *spec) updateResult(result result) {
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
