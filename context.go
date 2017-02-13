package elicit

import "testing"
import "fmt"
import "os"

// Context stores test machinery and maintains state between specs/scenarios/steps
type Context struct {
	specs      []*spec
	stepImpls  stepImplMap
	transforms transformMap
	log        log
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
		ctx.stepImpls.register(p, fn)
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

	ctx.stepImpls.checkTransforms(ctx.transforms)
}
