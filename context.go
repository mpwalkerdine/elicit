package elicit

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

// Context stores test machinery and maintains state between specs/scenarios/steps
type Context struct {
	specs          []*spec
	stepImpls      stepImpls
	transforms     transformMap
	beforeSpec     []Hook
	afterSpec      []Hook
	beforeScenario []Hook
	afterScenario  []Hook
	beforeStep     []Hook
	afterStep      []Hook
	unusedSteps    stepImpls
	log            log
}

// WithSpecsFolder recursively adds the path to the discovery path of specs
func (ctx *Context) WithSpecsFolder(path string) *Context {

	if i, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "warning: parsing spec folder %q: %s\n", path, err)
	} else if !i.IsDir() {
		fmt.Fprintf(os.Stderr, "warning: parsing spec folder %q: path is not a directory\n", path)
	} else {
		p := specParser{context: ctx}
		p.parseSpecFolder(path)
	}

	return ctx
}

// WithSteps registers steps from the supplied map of patterns to functions
func (ctx *Context) WithSteps(steps Steps) *Context {
	for p, fn := range steps {
		si := ctx.stepImpls.register(p, fn)
		if si != nil {
			ctx.unusedSteps = append(ctx.unusedSteps, si)
		}
	}
	return ctx
}

// WithTransforms registers step argument transforms from the suppled map of patterns to functions
func (ctx *Context) WithTransforms(txs Transforms) *Context {
	for p, fn := range txs {
		ctx.transforms.register(p, fn)
	}
	return ctx
}

// BeforeSpecs registers a function to be called before each spec
func (ctx *Context) BeforeSpecs(hook Hook) *Context {
	ctx.beforeSpec = append(ctx.beforeSpec, hook)
	return ctx
}

// AfterSpecs registers a function to be called after each spec
func (ctx *Context) AfterSpecs(hook Hook) *Context {
	ctx.afterSpec = append(ctx.afterSpec, hook)
	return ctx
}

// BeforeScenarios registers a function to be called before each scenario
func (ctx *Context) BeforeScenarios(hook Hook) *Context {
	ctx.beforeScenario = append(ctx.beforeScenario, hook)
	return ctx
}

// AfterScenarios registers a function to be called after each scenario
func (ctx *Context) AfterScenarios(hook Hook) *Context {
	ctx.afterScenario = append(ctx.afterScenario, hook)
	return ctx
}

// BeforeSteps registers a function to be called before each step
func (ctx *Context) BeforeSteps(hook Hook) *Context {
	ctx.beforeStep = append(ctx.beforeStep, hook)
	return ctx
}

// AfterSteps registers a function to be called after each step
func (ctx *Context) AfterSteps(hook Hook) *Context {
	ctx.afterStep = append(ctx.afterStep, hook)
	return ctx
}

// RunTests runs all the discovered specs as tests
func (ctx *Context) RunTests(ctxT *testing.T) *Context {
	allSkipped := true

	ctx.validate()

	for _, spec := range ctx.specs {

		var hookErr error
		for _, h := range ctx.beforeSpec {
			if hookErr = h.run(); hookErr != nil {
				spec.skipAllScenarios()
				spec.result = panicked
				fmt.Fprintf(os.Stderr, "panic during before spec hook: %s\n", hookErr)
				break
			}
		}

		ctxT.Run(spec.path+"/"+spec.name, func(specT *testing.T) {
			if hookErr != nil {
				specT.FailNow()
			}

			spec.run(specT)

			if !specT.Skipped() {
				allSkipped = false
			}

			if hookErr == nil {
				for _, h := range ctx.afterSpec {
					if hookErr = h.run(); hookErr != nil {
						spec.result = panicked
						fmt.Fprintf(os.Stderr, "panic during after spec hook: %s\n", hookErr)
						specT.Fail()
						break
					}
				}
			}
		})
	}

	ctx.log.writeToConsole()
	ctx.log.writeToFile()

	if allSkipped {
		ctxT.SkipNow()
	}

	return ctx
}

func (ctx *Context) validate() {
	if len(ctx.specs) == 0 {
		fmt.Fprintln(os.Stderr, "warning: No specifications found. Add a folder containing *.md files with Context.WithSpecsFolder().")
	}

	if len(ctx.stepImpls) == 0 {
		fmt.Fprintln(os.Stderr, "warning: No steps registered. Add some with Context.WithSteps().")
	}

	ctx.checkTransforms()
	ctx.resolveSteps()
}

func (ctx *Context) checkTransforms() {
	for _, impl := range ctx.stepImpls {
		fn := reflect.ValueOf(impl.fn)
		fnSig := fn.Type()

		paramCount, _, _ := ctx.stepImpls.countStepImplParams(fn)

		for p := 1; p < paramCount; p++ {
			pType := fnSig.In(p)

			if len(ctx.transforms[pType]) == 0 {
				fmt.Fprintf(os.Stderr, stepWarnNoTransform, impl, pType)
			}
		}
	}
}

func (ctx *Context) resolveSteps() {
	for _, spec := range ctx.specs {
		for _, scenario := range spec.scenarios {
			for _, step := range scenario.steps {
				ctx.matchStepImpl(step)
			}
		}
	}

	for _, impl := range ctx.unusedSteps {
		fmt.Fprintf(os.Stderr, stepWarnNotUsed, impl)
	}
}

func (ctx *Context) matchStepImpl(s *step) {
	// unresolved parameters count as pending
	if len(s.params) > 0 {
		return
	}

	type candidate struct {
		impl *stepImpl
		call func(t *testing.T)
	}
	candidates := []candidate{}

	for _, impl := range ctx.stepImpls {
		fn := reflect.ValueOf(impl.fn)
		params := impl.regex.FindStringSubmatch(s.text)

		if convertedParams, ok := ctx.transforms.convertParams(s, fn, params); ok {
			call := s.createCall(fn, convertedParams)
			candidates = append(candidates, candidate{impl, call})
		}
	}

	if len(candidates) == 1 {
		for _, c := range candidates {
			for i, us := range ctx.unusedSteps {
				if us == c.impl {
					// remove from log of unused steps
					ctx.unusedSteps = append(ctx.unusedSteps[:i], ctx.unusedSteps[i+1:]...)
					break
				}
			}
			s.impl = c.call
			// We set this to skipped in case it never gets a chance to run
			s.result = skipped
		}
	} else if len(candidates) > 1 {
		warning := fmt.Sprintf(stepWarnAmbiguous, s.text)
		for _, c := range candidates {
			warning += fmt.Sprintf("            - %s\n", c.impl)
		}
		fmt.Fprint(os.Stderr, warning)
	}
}

func (h Hook) run() error {
	return func() (rcvrErr error) {
		defer func() {
			if rcvr := recover(); rcvr != nil {
				if rerr, ok := rcvr.(error); ok {
					rcvrErr = rerr
				} else {
					rcvrErr = fmt.Errorf("%s", rcvr)
				}
			}
		}()

		h()
		return
	}()
}
