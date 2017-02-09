// Package elicit is a native go BDD testing framework using markdown for executable specifications
package elicit

import "regexp"

// New creates a new elicit context which stores specs, steps and transforms
func New() *Context {
	ctx := &Context{
		stepImpls:  map[*regexp.Regexp]interface{}{},
		transforms: map[*regexp.Regexp]StepArgumentTransform{},
	}

	// TODO(matt) discover these automatically
	ctx.transforms.register(`.*`, stringTransform)
	ctx.transforms.register(`-?\d+`, intTransform)
	// TODO(matt) reinstate ctx.transforms.register(`(?:.+,\s*)*.+`, commaSliceTransform)

	return ctx
}
