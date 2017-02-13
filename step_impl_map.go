package elicit

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"regexp/syntax"
	"strings"
	"testing"
)

type stepImplMap map[*regexp.Regexp]interface{}

const (
	stepWarnPrefix     = "warning: registered step %q => [%v] "
	stepWarnNotFunc    = stepWarnPrefix + "must be a function.\n"
	stepWarnBadRegex   = stepWarnPrefix + "has an invalid regular expression: %s.\n"
	stepWarnFirstParam = stepWarnPrefix + "has an invalid implementation. The first parameter must be of type *testing.T.\n"
	stepWarnParamCount = stepWarnPrefix + "captures %d parameter%s but the supplied implementation takes %d.\n"
	stepWarnParamType  = stepWarnPrefix + "has a parameter type %q for which no transforms exist.\n"
)

func (sim stepImplMap) register(pattern string, stepFunc interface{}) {
	if r, ok := sim.validate(pattern, stepFunc); ok {
		sim[r] = stepFunc
	}
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

func (sim stepImplMap) validate(pattern string, impl interface{}) (*regexp.Regexp, bool) {
	fn := reflect.ValueOf(impl)
	fnSig := fn.Type()

	if fnSig.Kind() != reflect.Func {
		fmt.Fprintf(os.Stderr, stepWarnNotFunc, pattern, fnSig)
		return nil, false
	}

	cleanPattern := strings.TrimSpace(pattern)
	cleanPattern = ensureCompleteMatch(pattern)
	regex, err := regexp.Compile(cleanPattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, stepWarnBadRegex, pattern, fnSig, err.(*syntax.Error).Code)
		return nil, false
	}

	typeTestingT := reflect.TypeOf((*testing.T)(nil))
	patternCaptures := regex.NumSubexp()
	if fnSig.NumIn() == 0 || fnSig.In(0) != typeTestingT {
		fmt.Fprintf(os.Stderr, stepWarnFirstParam, pattern, fnSig)
		return nil, false
	}

	// Note paramCount includes the first *testing.T parameter
	if paramCount, _, _ := sim.countStepImplParams(fn); paramCount-1 != patternCaptures {
		plural := ""
		if patternCaptures != 1 {
			plural = "s"
		}
		fmt.Fprintf(os.Stderr, stepWarnParamCount, pattern, fnSig, patternCaptures, plural, paramCount-1)
		return nil, false
	}

	return regex, true
}

func (sim stepImplMap) countStepImplParams(fn reflect.Value) (params, tables, textBlocks int) {
	tableType := reflect.TypeOf((*Table)(nil)).Elem()
	textBlockType := reflect.TypeOf((*TextBlock)(nil)).Elem()

	params = fn.Type().NumIn()
	for p := params - 1; p >= 0; p-- {
		thisParam := fn.Type().In(p)
		if thisParam == tableType {
			params--
			tables++
		} else if thisParam == textBlockType {
			params--
			textBlocks++
		} else {
			break
		}
	}

	return
}

func (sim stepImplMap) checkTransforms(tm transformMap) {
	for stepPattern, stepImpl := range sim {
		stepFn := reflect.ValueOf(stepImpl)
		stepFnSig := stepFn.Type()

		pCount, _, _ := sim.countStepImplParams(stepFn)

		for p := 1; p < pCount; p++ {
			pType := stepFnSig.In(p)

			if tm[pType] == nil {
				fmt.Fprintf(os.Stderr, stepWarnParamType, stepPattern, stepFnSig, pType)
			}
		}
	}
}
