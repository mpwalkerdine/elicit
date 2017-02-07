package elicit

import (
	"fmt"
	"testing"
)

// Context stores test machinery and maintains state between specs/scenarios/steps
type Context struct {
	specs           specCollection
	stepImpls       stepImplMap
	transforms      transformMap
	data            map[string]interface{}
	log             log
	currentSpec     *spec
	currentScenario *scenario
	currentStep     *step
}

// WithSpecsFolder recursively adds the path to the discovery path of specs
func (ctx *Context) WithSpecsFolder(path string) *Context {
	p := specParser{}
	specs := p.parseSpecFolder(path, ctx)
	ctx.specs = append(ctx.specs, specs...)
	return ctx
}

// WithSteps registers steps from the supplied map of patterns to functions
func (ctx *Context) WithSteps(steps map[string]interface{}) *Context {
	for p, fn := range steps {
		ctx.stepImpls.register(p, fn)
	}
	return ctx
}

// WithTransforms registers step argument transforms from the suppled map of patterns to functions
func (ctx *Context) WithTransforms(txs map[string]StepArgumentTransform) *Context {
	for p, fn := range txs {
		ctx.transforms.register(p, fn)
	}
	return ctx
}

// WithReportPath sets the output path for the execution summary.
// If path is the empty string, the report is written to stdout.
func (ctx *Context) WithReportPath(path string) *Context {
	return ctx
}

// RunTests runs all the discovered specs as tests
func (ctx *Context) RunTests(ctxT *testing.T) *Context {
	allSkipped := true

	for _, spec := range ctx.specs {
		ctxT.Run(spec.path+"/"+spec.name, func(specT *testing.T) {
			spec.runTest(specT)

			if !specT.Skipped() {
				allSkipped = false
			}
		})
	}

	fmt.Println(ctx.log.String())

	if allSkipped {
		ctxT.SkipNow()
	}

	return ctx
}

// Add stores data in the current context (scenario or spec) for use by other steps
func (ctx *Context) Add(key string, value interface{}) {
	ctx.data[key] = value
}

// Get retrieves data for the scenario context, falling back to the spec context
func (ctx *Context) Get(key string) interface{} {
	if value, found := ctx.data[key]; found {
		return value
	}

	return nil
}
