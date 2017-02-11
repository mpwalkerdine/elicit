// Package elicit is a native go BDD testing framework using markdown for executable specifications
package elicit

import "regexp"

// New creates a new elicit context which stores specs, steps and transforms
func New() *Context {
	ctx := &Context{
		stepImpls:  map[*regexp.Regexp]interface{}{},
		transforms: map[*regexp.Regexp]StepArgumentTransform{},
	}

	ctx.log.ctx = ctx

	ctx.transforms.init()

	return ctx
}
