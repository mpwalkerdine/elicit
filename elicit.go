// Package elicit is a native go BDD testing framework using markdown for executable specifications
package elicit

import (
	"flag"
	"fmt"
	"path/filepath"
)

var (
	reportFile = flag.String("elicit.report", "", "Path to save an execution report")
)

// Steps are used to register step implemenations against regex patterns
// The key must successfully compile to a regexp.Regexp and the values
// must be functions obeying the semantics of a step implementation.
// Step Implementations are of the form:
// func(t *testing.T, param string) {}
// The number of additional parameters should match the number of subgroups
// in the pattern used as a key.
type Steps map[string]interface{}

// Transforms are used to register functions which can convert strings
// captured by steps into the types expected by the step implementation.
// The string is a pattern which captured subgroups must match for the transform to be considered.
type Transforms map[string]interface{}

// New creates a new elicit context which stores specs, steps and transforms
func New() *Context {
	ctx := &Context{
		transforms: transformMap{},
	}

	ctx.log.ctx = ctx

	if *reportFile != "" {
		if reportFileAbs, err := filepath.Abs(*reportFile); err != nil {
			panic(fmt.Errorf("determining absolute path for %s: %s", *reportFile, err))
		} else {
			ctx.log.outpath = reportFileAbs
		}
	}

	ctx.transforms.init()

	return ctx
}
