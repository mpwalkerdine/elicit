package elicit

import (
	"bytes"
	"io"
	"os"
	"reflect"
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
	result     result
	log        bytes.Buffer
}

func (s *step) run(scenarioT *testing.T) {
	defer s.restoreStdout(s.redirectStdout())

	defer func() {
		if s.result > s.scenario.result {
			s.scenario.result = s.result
		}
	}()

	if s.impl == nil {
		s.result = pending
		scenarioT.SkipNow()
	} else if s.scenario.result != passed {
		s.result = skipped
		scenarioT.SkipNow()
	} else {
		s.impl(scenarioT)
	}
}

func (s *step) createCall(fn reflect.Value, params []reflect.Value) func(*testing.T) {
	return func(t *testing.T) {
		defer func() {
			if rcvr := recover(); rcvr != nil {
				s.result = panicked
				t.Error(rcvr)
			} else if t.Failed() {
				s.result = failed
			} else if t.Skipped() {
				s.result = skipped
			} else {
				s.result = passed
			}
		}()

		params[0] = reflect.ValueOf(t)
		fn.Call(params)
	}
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
