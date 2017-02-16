package elicit

import (
	"fmt"
	"os"
	"testing"
)

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

		var hookErr error
		for _, h := range s.context.beforeStep {
			if hookErr = h.run(); hookErr != nil {
				step.result = panicked
				s.result = panicked
				fmt.Fprintf(os.Stderr, "panic during before step hook: %s\n", hookErr)
				scenarioT.Fail()
				break
			}
		}

		if hookErr == nil {
			// Ensure after step hooks execute regardless of what happens in the step
			func() {
				defer func() {
					if step.result > s.result {
						s.result = step.result
					}

					for _, h := range s.context.afterStep {
						if hookErr = h.run(); hookErr != nil {
							step.result = panicked
							s.result = panicked
							fmt.Fprintf(os.Stderr, "panic during after step hook: %s\n", hookErr)
							scenarioT.Fail()
							break
						}
					}
				}()

				step.run(scenarioT)
			}()
		}
	}
}
