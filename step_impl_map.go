package elicit

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type stepImplMap map[*regexp.Regexp]interface{}

const (
	invalidFirstParam = "warning: The step pattern %q has an invalid implementation." +
		" The first parameter must be of type *testing.T.\n"
	countMismatch = "warning: The step pattern %q captures %d parameter%s but the supplied implementation takes %d.\n"
	noTransform   = "warning: The step pattern %q has a parameter type %q for which no transforms exist.\n"
)

func (sim stepImplMap) register(pattern string, stepFunc interface{}) {

	pattern = strings.TrimSpace(pattern)
	pattern = ensureCompleteMatch(pattern)

	p, err := regexp.Compile(pattern)

	if err != nil {
		panic(fmt.Sprintf("compiling step regexp %q, %s", pattern, err))
	}

	// TODO(matt) check the pattern captures the correct number of parameters

	sim[p] = stepFunc
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

func (sim stepImplMap) validate() {

	tType := reflect.TypeOf((*testing.T)(nil))

	for regex, impl := range sim {
		fn := reflect.ValueOf(impl)
		pattern := strings.TrimRight(strings.TrimLeft(regex.String(), "^"), "$")
		patternCaptures := regex.NumSubexp()

		if fn.Type().NumIn() == 0 || fn.Type().In(0) != tType {
			fmt.Fprintf(os.Stderr, invalidFirstParam, pattern)
			continue
		}

		// Note paramCount includes the first *testing.T parameter
		if paramCount, _, _ := countStepImplParams(fn); paramCount-1 != patternCaptures {
			plural := ""
			if patternCaptures != 1 {
				plural = "s"
			}
			fmt.Fprintf(os.Stderr, countMismatch, pattern, patternCaptures, plural, paramCount-1)
			continue
		}
	}
}

func countStepImplParams(f reflect.Value) (params, tables, textBlocks int) {
	tableType := reflect.TypeOf((*Table)(nil)).Elem()
	textBlockType := reflect.TypeOf((*TextBlock)(nil)).Elem()

	params = f.Type().NumIn()
	for p := params - 1; p >= 0; p-- {
		thisParam := f.Type().In(p)
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
