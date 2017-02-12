// Package elicit is a native go BDD testing framework using markdown for executable specifications
package elicit

import (
	"flag"
	"regexp"
)

var (
	reportFile = flag.String("elicit.report", "", "Path to save an execution report")
)

// New creates a new elicit context which stores specs, steps and transforms
func New() *Context {
	ctx := &Context{
		stepImpls:  map[*regexp.Regexp]interface{}{},
		transforms: map[*regexp.Regexp]StepArgumentTransform{},
	}

	ctx.log.ctx = ctx
	ctx.log.outpath = *reportFile

	ctx.transforms.init()

	return ctx
}
