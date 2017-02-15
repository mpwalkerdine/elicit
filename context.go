package elicit

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

// Context stores test machinery and maintains state between specs/scenarios/steps
type Context struct {
	specs       []*spec
	stepImpls   stepImpls
	unusedSteps stepImpls
	transforms  transformMap
	log         log
}

// WithSpecsFolder recursively adds the path to the discovery path of specs
func (ctx *Context) WithSpecsFolder(path string) *Context {
	p := specParser{context: ctx}
	p.parseSpecFolder(path)
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

// RunTests runs all the discovered specs as tests
func (ctx *Context) RunTests(ctxT *testing.T) *Context {
	allSkipped := true

	ctx.validate()

	for _, spec := range ctx.specs {
		ctxT.Run(spec.path+"/"+spec.name, func(specT *testing.T) {
			spec.runTest(specT)

			if !specT.Skipped() {
				allSkipped = false
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
		fmt.Fprintln(os.Stderr, "warning: No specifications found. Add a folder containing *.spec files with Context.WithSpecsFolder().")
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
				step.impl = ctx.matchStepImpl(step)
			}
		}
	}

	for _, impl := range ctx.unusedSteps {
		fmt.Fprintf(os.Stderr, stepWarnNotUsed, impl)
	}
}

func (ctx *Context) matchStepImpl(s *step) func(*testing.T) {
	type candidate struct {
		impl  *stepImpl
		match func(t *testing.T)
	}
	candidates := []candidate{}

	for _, impl := range ctx.stepImpls {
		fn := reflect.ValueOf(impl.fn)
		params := impl.regex.FindStringSubmatch(s.text)

		if convertedParams, ok := ctx.transforms.convertParams(s, fn, params); ok {
			match := func(t *testing.T) {
				convertedParams[0] = reflect.ValueOf(t)
				fn.Call(convertedParams)
			}
			candidates = append(candidates, candidate{impl, match})
		}
	}

	if len(candidates) == 1 {
		for _, c := range candidates {
			for i, us := range ctx.unusedSteps {
				if us == c.impl {
					ctx.unusedSteps = append(ctx.unusedSteps[:i], ctx.unusedSteps[i+1:]...)
					break
				}
			}
			return c.match
		}
	}

	if len(candidates) > 1 {
		warning := fmt.Sprintf(stepWarnAmbiguous, s.text)
		for _, c := range candidates {
			warning += fmt.Sprintf("            - %s\n", c.impl)
		}
		fmt.Fprint(os.Stderr, warning)
	}

	return nil
}
