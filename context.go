package elicit

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"bytes"
	"reflect"
)

// Context stores test machinery and maintains state between specs/scenarios/steps
type Context struct {
	specs           specCollection
	stepImpls       stepImplMap
	transforms      transformMap
	currentSpec     *SpecContext
	currentScenario *ScenarioContext
}

// StepArgumentTransform transforms captured groups in the step pattern to a function parameter type
// Note that if the actual string cannot be converted to the target type by the transform, it should return false
type StepArgumentTransform func(*Context, string, reflect.Type) (interface{}, bool)

type specCollection []specDef
type stepImplMap map[*regexp.Regexp]interface{}
type transformMap map[*regexp.Regexp]StepArgumentTransform

// New creates a new elicit context which stores specs, steps and transforms
func New() *Context {
	ctx := &Context{
		stepImpls:  make(map[*regexp.Regexp]interface{}),
		transforms: make(map[*regexp.Regexp]StepArgumentTransform),
	}

	// TODO discover these automatically
	ctx.registerTransform(`.*`, stringTransform)
	ctx.registerTransform(`-?\d+`, intTransform)
	ctx.registerTransform(`(?:.+,\s*)*.+`, commaSliceTransform)

	return ctx
}

// WithSpecsFolder recursively adds the path to the discovery path of specs
func (ctx *Context) WithSpecsFolder(path string) *Context {
	ctx.specs.parseSpecFolder(path)
	return ctx
}

// WithSteps registers steps from the supplied map of patterns to functions
func (ctx *Context) WithSteps(steps map[string]interface{}) *Context {
	for p, fn := range steps {
		ctx.registerStep(p, fn)
	}
	return ctx
}

func (ctx *Context) registerStep(pattern string, stepFunc interface{}) {

	pattern = strings.TrimSpace(pattern)
	pattern = ensureCompleteMatch(pattern)

	p, err := regexp.Compile(pattern)

	if err != nil {
		panic(fmt.Sprintf("compiling step regexp %q, %s", pattern, err))
	}

	// TODO(matt) check the pattern captures the correct number of parameters

	ctx.stepImpls[p] = stepFunc
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

// WithTransforms registers step argument transforms from the suppled map of patterns to functions
func (ctx *Context) WithTransforms(txs map[string]StepArgumentTransform) *Context {
	for p, fn := range txs {
		ctx.registerTransform(p, fn)
	}
	return ctx
}

func (ctx *Context) registerTransform(pattern string, transform StepArgumentTransform) {
	pattern = ensureCompleteMatch(pattern)

	p, err := regexp.Compile(pattern)

	if err != nil {
		panic(fmt.Sprintf("compiling transform regexp %q, %s", pattern, err))
	}

	ctx.transforms[p] = transform
}

// RunTests runs all the discovered specs as tests
func (ctx *Context) RunTests(t *testing.T) {
	for _, spec := range ctx.specs {
		ctx.runSpecTest(t, spec)
	}
}

func (ctx *Context) runSpecTest(t *testing.T, spec specDef) {
	ctx.currentSpec = &SpecContext{
		name: spec.Name,
	}

	t.Run(spec.Name, func(t *testing.T) {
		for _, scenario := range spec.Scenarios {
			ctx.runScenarioTest(t, scenario)
		}
	})
}

func (ctx *Context) runScenarioTest(t *testing.T, scenario scenarioDef) {
	ctx.currentScenario = &ScenarioContext{
		name:    scenario.Name,
		skipped: false,
		failed:  false,
		logbuf:  bytes.Buffer{},
	}
	ctx.logScenarioStart()

	t.Run(scenario.Name, func(t *testing.T) {
		for _, b := range scenario.Spec.BeforeSteps {
			ctx.runStep(t, b)
		}

		for _, s := range scenario.Steps {
			ctx.runStep(t, s)
		}

		for _, a := range scenario.Spec.AfterSteps {
			ctx.runStep(t, a)
		}

		log := string(ctx.currentScenario.logbuf.Bytes())
		if ctx.currentScenario.failed {
			t.Errorf(log)
		} else if ctx.currentScenario.skipped {
			t.Skipf(log)
		}
	})
}

func (ctx *Context) runStep(t *testing.T, step stepDef) {
	ctx.currentScenario.currentStep = step.Text

	defer func() {
		if r := recover(); r != nil {
			ctx.Failf("panic during step execution: %s", r)
		}
	}()

	for regex, fn := range ctx.stepImpls {
		f := reflect.ValueOf(fn)
		params := regex.FindStringSubmatch(step.Text)

		if in, ok := ctx.convertParams(f, params, step.Tables); ok {

			if !ctx.currentScenario.skipped && !ctx.currentScenario.failed {
				f.Call(in)
			} else {
				ctx.Skip("")
			}

			if !ctx.currentScenario.skipped && !ctx.currentScenario.failed {
				ctx.stepPassed()
			}

			return
		}
	}

	ctx.stepNotFound()
}

func (ctx *Context) convertParams(f reflect.Value, stringParams []string, tables []stringTable) ([]reflect.Value, bool) {

	if stringParams == nil {
		return nil, false
	}

	paramCount := f.Type().NumIn()
	tableParamCount := 0
	tableType := reflect.TypeOf((*Table)(nil)).Elem()

	for p := paramCount - 1; p >= 0; p-- {
		thisParam := f.Type().In(p)
		if thisParam == tableType {
			paramCount--
			tableParamCount++
		} else {
			break
		}
	}

	if len(stringParams) != paramCount || tableParamCount != len(tables) {
		return nil, false
	}

	c := make([]reflect.Value, len(stringParams))
	for i, param := range stringParams {
		if i == 0 {
			if f.Type().In(0) != reflect.TypeOf(ctx) {
				return nil, false
			}
			c[i] = reflect.ValueOf(ctx)
		} else {
			pt := f.Type().In(i)

			if t, ok := ctx.convertParam(param, pt); ok {
				c[i] = t
			} else {
				return nil, false
			}
		}
	}

	for _, t := range tables {
		c = append(c, reflect.ValueOf(makeTable(t)))
	}

	return c, true
}

func (ctx *Context) convertParam(s string, target reflect.Type) (reflect.Value, bool) {
	for regex, tx := range ctx.transforms {
		params := regex.FindStringSubmatch(s)
		if params == nil || len(params) != 1 {
			continue
		}

		f := reflect.ValueOf(tx)
		fTyp := f.Type()

		in := make([]reflect.Value, fTyp.NumIn())
		in[0] = reflect.ValueOf(ctx)
		in[1] = reflect.ValueOf(s)
		in[2] = reflect.ValueOf(target)

		out := f.Call(in)
		if out[1].Interface().(bool) {
			return reflect.ValueOf(out[0].Interface()), true
		}
	}

	return reflect.Value{}, false
}

// Skip skips the current step execution and all subsequent steps
func (ctx *Context) Skip(format string, args ...interface{}) {
	ctx.logStepResult("⤹", format, args...)
	ctx.currentScenario.skipped = true
}

func (ctx *Context) stepNotFound() {
	ctx.logStepResult("?", "")
	ctx.currentScenario.skipped = true
}

func (ctx *Context) stepPassed() {
	ctx.logStepResult("✓", "")
}

// Fail records test step failure
func (ctx *Context) Fail() {
	ctx.Failf("")
}

// Failf logs the supplied message and records test step failure
func (ctx *Context) Failf(format string, args ...interface{}) {
	ctx.logStepResult("✘", format, args...)
	ctx.currentScenario.failed = true
}

func (ctx *Context) logScenarioStart() {
	fmt.Fprintf(&ctx.currentScenario.logbuf, "\n\t%s\n", ctx.currentScenario.name)
}

func (ctx *Context) logStepResult(prefix, format string, args ...interface{}) {
	if len(format) > 0 {
		format = fmt.Sprintf("\t\t%s %s\t(%s)\n", prefix, ctx.currentScenario.currentStep, format)
		fmt.Fprintf(&ctx.currentScenario.logbuf, format, args...)
	} else {
		fmt.Fprintf(&ctx.currentScenario.logbuf, "\t\t%s %s\n", prefix, ctx.currentScenario.currentStep)
	}
}
