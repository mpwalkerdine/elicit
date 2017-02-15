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

type stepImpl struct {
	regex *regexp.Regexp
	fn    interface{}
}

type stepImpls []*stepImpl

const (
	stepWarnPrefix      = "warning: registered step %q => [%v] "
	stepWarnNotFunc     = stepWarnPrefix + "must be a function.\n"
	stepWarnBadRegex    = stepWarnPrefix + "has an invalid regular expression: %s.\n"
	stepWarnFirstParam  = stepWarnPrefix + "has an invalid implementation. The first parameter must be of type *testing.T.\n"
	stepWarnParamCount  = stepWarnPrefix + "captures %d parameter%s but the supplied implementation takes %d.\n"
	stepWarnNoTransform = "warning: registered step %s has a parameter type %q for which no transforms exist.\n"
	stepWarnNotUsed     = "warning: registered step %s is not used.\n"
	stepWarnAmbiguous   = "warning: step %q is ambiguous:\n"
)

var (
	typeTestingT = reflect.TypeOf((*testing.T)(nil))
)

func (s *stepImpl) String() string {
	p := s.regex.String()
	p = strings.TrimLeft(p, "^")
	p = strings.TrimRight(p, "$")
	return fmt.Sprintf("%q => [%v]", p, reflect.TypeOf(s.fn))
}

func (si *stepImpls) register(pattern string, stepFunc interface{}) *stepImpl {
	if r, ok := si.validate(pattern, stepFunc); ok {
		*si = append(*si, &stepImpl{regex: r, fn: stepFunc})
		return (*si)[len(*si)-1]
	}
	return nil
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

func (si *stepImpls) validate(pattern string, impl interface{}) (*regexp.Regexp, bool) {
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

	patternCaptures := regex.NumSubexp()
	if fnSig.NumIn() == 0 || fnSig.In(0) != typeTestingT {
		fmt.Fprintf(os.Stderr, stepWarnFirstParam, pattern, fnSig)
		return nil, false
	}

	// Note paramCount includes the first *testing.T parameter
	if paramCount, _, _ := si.countStepImplParams(fn); paramCount-1 != patternCaptures {
		plural := ""
		if patternCaptures != 1 {
			plural = "s"
		}
		fmt.Fprintf(os.Stderr, stepWarnParamCount, pattern, fnSig, patternCaptures, plural, paramCount-1)
		return nil, false
	}

	return regex, true
}

func (si *stepImpls) countStepImplParams(fn reflect.Value) (params, tables, textBlocks int) {
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
