package elicit

import (
	"fmt"
	"reflect"
	"strings"
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

func (s *step) run(scenarioT *testing.T) {
	s.context.currentStep = s
	skip := (s.scenario.result == failed || s.scenario.result == skipped)

	for _, text := range s.resolveStepParams(scenarioT) {
		s.scenario.stepsRun++

		scenarioT.Run(s.testName(text), func(stepT *testing.T) {
			impl, found := s.matchStepImpl(stepT, text)

			if !found {
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

		s.scenario.updateResult(s.result)
		s.context.log.step(s, text)
	}
}

func (s *step) testName(text string) string {
	return fmt.Sprintf("%d_%s", s.scenario.stepsRun, text)
}

func (s *step) matchStepImpl(t *testing.T, text string) (func(), bool) {
	for regex, impl := range s.context.stepImpls {
		fn := reflect.ValueOf(impl)
		params := regex.FindStringSubmatch(text)

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

			if t, ok := s.convertParam(param, pt); ok {
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

func (s *step) convertParam(param string, target reflect.Type) (reflect.Value, bool) {
	for regex, tx := range s.context.transforms {
		params := regex.FindStringSubmatch(param)
		if params == nil {
			continue
		}

		f := reflect.ValueOf(tx)

		in := []reflect.Value{
			reflect.ValueOf(params),
			reflect.ValueOf(target),
		}

		out := f.Call(in)
		if out[1].Interface().(bool) {
			return reflect.ValueOf(out[0].Interface()), true
		}
	}

	return reflect.Value{}, false
}

// TODO(matt) consider resolving these during parsing rather than execution?
func (s *step) resolveStepParams(scenarioT *testing.T) []string {
	if len(s.params) == 0 {
		return []string{s.text}
	}

	table := stringTable{}
	resolved := []string{}
	found := false

	if s.scenario != nil {
		table, found = s.findTableWithParams(s.scenario.tables)
	}

	if !found {
		table, found = s.findTableWithParams(s.spec.tables)
	}

	if !found {
		s.scenario.stepsRun++
		scenarioT.Run(s.testName(s.text), func(stepT *testing.T) {
			s.result = undefined
			stepT.SkipNow()
		})
		s.scenario.updateResult(s.result)
		s.context.log.step(s, s.text)
		return resolved
	}

	m := table.columnNameToIndexMap()
	for _, row := range table[1:] {
		text := s.text
		for _, p := range s.params {
			pname := strings.TrimSuffix(strings.TrimPrefix(p, "<"), ">")
			pvalue := row[m[pname]]
			text = strings.Replace(text, p, pvalue, -1)
		}
		resolved = append(resolved, text)
	}

	return resolved
}

func (s *step) findTableWithParams(tables []stringTable) (stringTable, bool) {
	for _, t := range tables {
		if t.hasParams(s.params) {
			return t, true
		}
	}
	return nil, false
}
