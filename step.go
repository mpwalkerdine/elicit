package elicit

import (
	"fmt"
	"reflect"
	"testing"
)

type stepResult int

const (
	undefined stepResult = iota
	skipped
	failed
	panicked
	passed
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
	result     stepResult
}

func (s *step) run(scenarioT *testing.T) stepResult {
	skip := (s.scenario.result == failed || s.scenario.result == skipped)

	scenarioT.Run(s.testName(), func(stepT *testing.T) {

		// unresolved parameters count as undefined
		if len(s.params) > 0 {
			s.result = undefined
			stepT.SkipNow()
		} else if impl, found := s.matchStepImpl(stepT); !found {
			s.result = undefined
			stepT.SkipNow()
		} else if !s.force && skip {
			s.result = skipped
			stepT.SkipNow()
		} else {

			defer func() {
				if r := recover(); r != nil {
					s.result = panicked
					stepT.Error(r)
				}
			}()

			impl()

			if skip {
				s.forced = true
			}

			if stepT.Failed() {
				s.result = failed
			} else if stepT.Skipped() {
				s.result = skipped
			} else {
				s.result = passed
			}
		}
	})

	return s.result
}

func (s *step) testName() string {
	return fmt.Sprintf("%d_%s", s.scenario.stepsRun, s.text)
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

	paramCount := f.Type().NumIn()
	tableParamCount := 0
	textBlockParamCount := 0
	tableType := reflect.TypeOf((*Table)(nil)).Elem()
	textBlockType := reflect.TypeOf((*TextBlock)(nil)).Elem()

	for p := paramCount - 1; p >= 0; p-- {
		thisParam := f.Type().In(p)
		if thisParam == tableType {
			paramCount--
			tableParamCount++
		} else if thisParam == textBlockType {
			paramCount--
			textBlockParamCount++
		} else {
			break
		}
	}

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
