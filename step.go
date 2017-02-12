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
			s.result = undefined
			stepT.Skip("unresolved parameters:", s.params)
		} else if impl, found := s.matchStepImpl(stepT); !found {
			s.result = undefined
			stepT.Skip("no matching step implementation")
		} else if !s.force && skip {
			stepT.SkipNow()
		} else {
			if skip {
				s.forced = true
			}
			impl()
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

func (s *step) matchStepImpl(t *testing.T) (func(), bool) {
	for regex, impl := range s.context.stepImpls {
		fn := reflect.ValueOf(impl)
		params := regex.FindStringSubmatch(s.text)

		if convertedParams, ok := s.convertParams(t, fn, params); ok {
			return func() {
				fn.Call(convertedParams)
			}, true
		}
	}

	return nil, false
}

func (s *step) convertParams(t *testing.T, f reflect.Value, stringParams []string) ([]reflect.Value, bool) {

	if stringParams == nil {
		return nil, false
	}

	paramCount, tableParamCount, textBlockParamCount := countStepImplParams(f)

	if len(stringParams) != paramCount || tableParamCount != len(s.tables) || textBlockParamCount != len(s.textBlocks) {
		return nil, false
	}

	c := make([]reflect.Value, len(stringParams))
	for i, param := range stringParams {
		if i == 0 {
			if f.Type().In(0) != reflect.TypeOf(t) {
				return nil, false
			}
			c[i] = reflect.ValueOf(t)
		} else {
			pt := f.Type().In(i)

			if t, ok := s.context.transforms.convertParam(param, pt); ok {
				c[i] = t
			} else {
				return nil, false
			}
		}
	}

	for _, t := range s.tables {
		c = append(c, reflect.ValueOf(makeTable(t)))
	}

	for _, tb := range s.textBlocks {
		c = append(c, reflect.ValueOf(tb))
	}

	return c, true
}
