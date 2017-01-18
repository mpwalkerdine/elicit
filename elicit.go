package elicit

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"reflect"
	"strconv"
)

// CurrentContext is a single instance used by all elicit tests
var CurrentContext = &Context{
	steps: make(map[*regexp.Regexp]interface{}),
}

// Context stores test machinery and maintains state between specs/scenarios/steps
type Context struct {
	steps     map[*regexp.Regexp]interface{}
	specT     *testing.T
	scenarioT *testing.T
}

// BeginSpecTest registers the start of Spec
func (e *Context) BeginSpecTest(t *testing.T) {
	e.specT = t
}

// BeginScenarioTest registers the start of a Scenario
func (e *Context) BeginScenarioTest(t *testing.T) {
	e.scenarioT = t
}

// RegisterStep maps a Regexpr to a step implementation
func (e *Context) RegisterStep(pattern string, stepFunc interface{}) {

	pattern = strings.TrimSpace(pattern)

	if !strings.HasPrefix(pattern, "^") {
		pattern = "^" + pattern
	}

	if !strings.HasSuffix(pattern, "$") {
		pattern = pattern + "$"
	}

	p, err := regexp.Compile(pattern)

	if err != nil {
		panic(fmt.Sprintf("compiling regexp %q, %s", pattern, err))
	}

	e.steps[p] = stepFunc
}

// RunStep matches stepText to a registered step implementation and invokes it
func (e *Context) RunStep(stepText string) {
	for regex, fn := range e.steps {
		f := reflect.ValueOf(fn)
		params := regex.FindStringSubmatch(stepText)

		if in, ok := convertParams(f, params); ok {
			f.Call(in)
			return
		}
	}
}

func convertParams(f reflect.Value, stringParams []string) ([]reflect.Value, bool) {

	if stringParams == nil || len(stringParams) != f.Type().NumIn() {
		return nil, false
	}

	c := make([]reflect.Value, len(stringParams))
	for i, param := range stringParams {
		if i == 0 {
			if f.Type().In(0) != reflect.TypeOf(CurrentContext) {
				return nil, false
			}
			c[i] = reflect.ValueOf(CurrentContext)
		} else {
			pt := f.Type().In(i)

			if t, ok := convertParam(param, pt); ok {
				c[i] = t
			} else {
				return nil, false
			}
		}
	}

	return c, true
}

func convertParam(s string, target reflect.Type) (reflect.Value, bool) {
	t := reflect.New(target).Elem().Interface()

	switch t.(type) {
	case int:
		if t, err := strconv.Atoi(s); err == nil {
			return reflect.ValueOf(t), true
		}
	}

	return reflect.Value{}, false
}

// Fail records test failure
func (e *Context) Fail() {
	e.scenarioT.Fail()
}

// Assert that the parameter is true, otherwise fails
func (e *Context) Assert(shouldBeTrue bool) {
	if !shouldBeTrue {
		e.Fail()
	}
}
