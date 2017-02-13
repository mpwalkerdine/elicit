package elicit

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"regexp/syntax"
	"strconv"
	"strings"
)

type transformMap map[reflect.Type]map[*regexp.Regexp]interface{}

const (
	txWarnPrefix    = "warning: registered transform %q => [%v] "
	txWarnNotFunc   = txWarnPrefix + "must be a function.\n"
	txWarnBadRegex  = txWarnPrefix + "has an invalid regular expression: %s.\n"
	txWarnParamType = txWarnPrefix + "must take one argument of type []string.\n"
	txWarnReturn    = txWarnPrefix + "must return precisely one value.\n"
)

func (tm transformMap) init() {
	tm.register(`.*`, func(params []string) string {
		return params[0]
	})

	tm.register(`-?\d+`, func(params []string) int {
		i, err := strconv.Atoi(params[0])

		if err != nil {
			panic(fmt.Errorf("converting %q to int: %s", params[0], err))
		}

		return i
	})

	tm.register(`(?:.+,\s*)*.+`, func(params []string) []string {
		ss := []string{}

		for _, s := range strings.Split(params[0], ",") {
			s = strings.TrimSpace(s)
			ss = append(ss, s)
		}

		return ss
	})

	tm.register(`(?:-?\d+,\s*)*-?\d+`, func(params []string) []int {
		si := []int{}

		for _, s := range strings.Split(params[0], ",") {
			s = strings.TrimSpace(s)
			i, err := strconv.Atoi(s)
			if err != nil {
				panic(fmt.Errorf("converting %q to int: %s", s, err))
			}
			si = append(si, i)
		}

		return si
	})
}

func (tm transformMap) register(pattern string, transform interface{}) {
	if regex, typ, ok := tm.validate(pattern, transform); ok {
		if tm[typ] == nil {
			tm[typ] = map[*regexp.Regexp]interface{}{}
		}

		tm[typ][regex] = transform
	}
}

func (tm transformMap) validate(pattern string, transform interface{}) (*regexp.Regexp, reflect.Type, bool) {
	fn := reflect.ValueOf(transform)
	fnSig := fn.Type()

	if fnSig.Kind() != reflect.Func {
		fmt.Fprintf(os.Stderr, txWarnNotFunc, pattern, fnSig)
		return nil, nil, false
	}

	cleanPattern := ensureCompleteMatch(pattern)
	regex, err := regexp.Compile(cleanPattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, txWarnBadRegex, pattern, fnSig, err.(*syntax.Error).Code)
		return nil, nil, false
	}

	stringSliceType := reflect.TypeOf((*[]string)(nil)).Elem()
	if fnSig.NumIn() != 1 || fnSig.In(0) != stringSliceType {
		fmt.Fprintf(os.Stderr, txWarnParamType, pattern, fnSig)
		return nil, nil, false
	}

	if fnSig.NumOut() != 1 {
		fmt.Fprintf(os.Stderr, txWarnReturn, pattern, fnSig)
		return nil, nil, false
	}

	typ := fnSig.Out(0)

	return regex, typ, true
}

func (tm transformMap) convertParam(param string, target reflect.Type) (reflect.Value, bool) {
	for regex, tx := range tm[target] {
		params := regex.FindStringSubmatch(param)
		if params == nil {
			continue
		}

		f := reflect.ValueOf(tx)

		in := []reflect.Value{
			reflect.ValueOf(params),
		}

		out := f.Call(in)
		return reflect.ValueOf(out[0].Interface()), true
	}

	return reflect.Value{}, false
}
