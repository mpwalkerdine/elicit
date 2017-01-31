package elicit

import (
	"fmt"
	"regexp"
	"strings"

	"bytes"
	"reflect"
)

func init() {
	e := CurrentContext

	e.RegisterTransform(`.*`, stringTransform)
	e.RegisterTransform(`-?\d+`, intTransform)
	e.RegisterTransform(`(?:.+,\s*)*.+`, commaSliceTransform)
}

// CurrentContext is a single instance used by all elicit tests
var CurrentContext = &Context{
	steps:      make(map[*regexp.Regexp]interface{}),
	transforms: make(map[*regexp.Regexp]StepArgumentTransform),
}

// StepArgumentTransform transforms captured groups in the step pattern to a function parameter type
type StepArgumentTransform func(string, reflect.Type) (interface{}, bool)

// Result represents the outcome of a scenario test
type Result int

const (
	// Passed means the the scenario passed
	Passed Result = iota
	// Skipped means the scenario was skipped, normally due to undefined steps
	Skipped
	// Failed means the scenario failed, either due to a failed assertion or an error
	Failed
)

// Context stores test machinery and maintains state between specs/scenarios/steps
type Context struct {
	steps        map[*regexp.Regexp]interface{}
	transforms   map[*regexp.Regexp]StepArgumentTransform
	specName     string
	scenarioName string
	current      string
	skipped      bool
	failed       bool
	logbuf       bytes.Buffer
}

// BeginSpecTest registers the start of Spec
func (e *Context) BeginSpecTest(name string) {
	e.specName = name
	e.logSpecStart()
}

// BeginScenarioTest registers the start of a Scenario
func (e *Context) BeginScenarioTest(name string) {
	e.scenarioName = name
	e.skipped = false
	e.failed = false
	e.logbuf = bytes.Buffer{}
	e.logScenarioStart()
}

// EndScenarioTest marks the end of a scenario and signals the outcome
func (e *Context) EndScenarioTest() (Result, string) {

	log := string(e.logbuf.Bytes())

	if e.failed {
		return Failed, log
	}

	if e.skipped {
		return Skipped, log
	}

	return Passed, log
}

// EndSpecTest marks the end of a spec
func (e *Context) EndSpecTest() {

}

// RegisterStep maps a Regexpr to a step implementation
func (e *Context) RegisterStep(pattern string, stepFunc interface{}) {

	pattern = strings.TrimSpace(pattern)
	pattern = ensureCompleteMatch(pattern)

	p, err := regexp.Compile(pattern)

	if err != nil {
		panic(fmt.Sprintf("compiling step regexp %q, %s", pattern, err))
	}

	e.steps[p] = stepFunc
}

// RegisterTransform registers a function which will be used when matching step implementation parameters
// Note that if the actual string cannot be converted to the target type by the transform, it should return false
func (e *Context) RegisterTransform(pattern string, transform StepArgumentTransform) {
	pattern = ensureCompleteMatch(pattern)

	p, err := regexp.Compile(pattern)

	if err != nil {
		panic(fmt.Sprintf("compiling transform regexp %q, %s", pattern, err))
	}

	e.transforms[p] = transform
}

func ensureCompleteMatch(pattern string) string {
	if !strings.HasPrefix(pattern, "^") {
		pattern = "^" + pattern
	}

	if !strings.HasSuffix(pattern, "$") {
		pattern = pattern + "$"
	}

	return pattern
}

// RunStep matches stepText to a registered step implementation and invokes it
func (e *Context) RunStep(stepText string) {
	e.current = stepText

	defer func() {
		if r := recover(); r != nil {
			e.Failf("panic during step execution: %s", r)
		}
	}()

	for regex, fn := range e.steps {
		f := reflect.ValueOf(fn)
		params := regex.FindStringSubmatch(stepText)

		if in, ok := convertParams(f, params); ok {

			if !e.skipped && !e.failed {
				f.Call(in)
			} else {
				e.Skip("")
			}

			if !e.skipped && !e.failed {
				e.stepPassed()
			}

			return
		}
	}

	e.stepNotFound()
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
	for regex, tx := range CurrentContext.transforms {
		params := regex.FindStringSubmatch(s)
		if params == nil || len(params) != 1 {
			continue
		}

		f := reflect.ValueOf(tx)
		fTyp := f.Type()

		in := make([]reflect.Value, fTyp.NumIn())
		in[0] = reflect.ValueOf(s)
		in[1] = reflect.ValueOf(target)

		out := f.Call(in)
		if out[1].Interface().(bool) {
			return reflect.ValueOf(out[0].Interface()), true
		}
	}

	return reflect.Value{}, false
}

// Skip skips the current step execution and all subsequent steps
func (e *Context) Skip(format string, args ...interface{}) {
	e.logStepResult("⤹", format, args...)
	e.skipped = true
}

func (e *Context) stepNotFound() {
	e.logStepResult("?", "")
	e.skipped = true
}

func (e *Context) stepPassed() {
	e.logStepResult("✓", "")
}

// Fail records test step failure
func (e *Context) Fail() {
	e.Failf("")
}

// Failf logs the supplied message and records test step failure
func (e *Context) Failf(format string, args ...interface{}) {
	e.logStepResult("✘", format, args...)
	e.failed = true
}

func (e *Context) logSpecStart() {
	fmt.Fprintf(&e.logbuf, "%s\n", e.specName)
}

func (e *Context) logScenarioStart() {
	fmt.Fprintf(&e.logbuf, "\n\t%s\n", e.scenarioName)
}

func (e *Context) logStepResult(prefix, format string, args ...interface{}) {
	if len(format) > 0 {
		format = fmt.Sprintf("\t\t%s %s\t(%s)\n", prefix, e.current, format)
		fmt.Fprintf(&e.logbuf, format, args...)
	} else {
		fmt.Fprintf(&e.logbuf, "\t\t%s %s\n", prefix, e.current)
	}
}
