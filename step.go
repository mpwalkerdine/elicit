package elicit

import (
	"bytes"
	"io"
	"os"
	"testing"
)

type step struct {
	context    *Context
	spec       *spec
	scenario   *scenario
	text       string
	params     []string
	tables     []stringTable
	textBlocks []TextBlock
	impl       func(*testing.T)
	force      bool
	forced     bool
	result     result
	log        bytes.Buffer
}

func (s *step) run(scenarioT *testing.T) {
	defer s.restoreStdout(s.redirectStdout())

	skip := (s.scenario.result != passed)

	scenarioT.Run("", func(stepT *testing.T) {
		defer func() {
			if rcvr := recover(); rcvr != nil {
				s.result = panicked
				stepT.Error(rcvr)
			}

			// Don't overwrite existing result
			if s.result == notrun {
				if stepT.Failed() {
					s.result = failed
				} else if stepT.Skipped() {
					s.result = skipped
				} else {
					s.result = passed
				}
			}
		}()

		// unresolved parameters count as undefined
		if len(s.params) > 0 {
			s.result = pending
			stepT.Skip("unresolved parameters:", s.params)
		} else if s.impl == nil {
			s.result = pending
			stepT.Skip("no matching step implementation")
		} else if !s.force && skip {
			stepT.SkipNow()
		} else {
			if skip {
				s.forced = true
			}
			s.impl(stepT)
		}
	})
}

func (s *step) redirectStdout() (*os.File, chan bool) {
	stdout := os.Stdout

	r, w, err := os.Pipe()

	if err != nil {
		return stdout, nil
	}

	waitChan := make(chan bool)
	go func() {
		// This will continue until w is closed
		io.Copy(&s.log, r)

		// Signal that copying has been completed
		waitChan <- true
	}()

	os.Stdout = w

	return stdout, waitChan
}

func (s *step) restoreStdout(stdout *os.File, waitChan chan bool) {
	w := os.Stdout
	os.Stdout = stdout
	if w != stdout {
		w.Close()
		<-waitChan
	}
}
