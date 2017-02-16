package elicit

import (
	"fmt"
	"os"
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

func (s *spec) run(specT *testing.T) {
	for _, scenario := range s.scenarios {

		var hookErr error
		for _, h := range s.context.beforeScenario {
			if hookErr = h.run(); hookErr != nil {
				scenario.result = panicked
				fmt.Fprintf(os.Stderr, "panic during before scenario hook: %s\n", hookErr)
				break
			}
		}

		specT.Run(scenario.name, func(scenarioT *testing.T) {
			if hookErr != nil {
				scenarioT.FailNow()
			}

			// Ensure the after scenario hooks are run regardless of the result
			func() {
				defer func() {
					if hookErr == nil {

						for _, h := range s.context.afterScenario {
							if hookErr = h.run(); hookErr != nil {
								scenario.result = panicked
								fmt.Fprintf(os.Stderr, "panic during after scenario hook: %s\n", hookErr)
								scenarioT.FailNow()
								break
							}
						}
					}
				}()

				scenario.run(scenarioT)
			}()
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

func (s *spec) skipAllScenarios() {
	for _, scenario := range s.scenarios {
		scenario.result = skipped
	}
}
