package elicit

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// StepArgumentTransform transforms captured groups in the step pattern to a function parameter type
// Note that if the actual string cannot be converted to the target type by the transform, it should return false
type StepArgumentTransform func([]string, reflect.Type) (interface{}, bool)

type transformMap map[*regexp.Regexp]StepArgumentTransform

func (tm *transformMap) init() {
	tm.register(`.*`, func(params []string, t reflect.Type) (interface{}, bool) {
		if t != reflect.TypeOf((*string)(nil)).Elem() {
			return nil, false
		}

		return params[0], true
	})

	tm.register(`-?\d+`, func(params []string, t reflect.Type) (interface{}, bool) {
		if t != reflect.TypeOf((*int)(nil)).Elem() {
			return nil, false
		}

		if t, err := strconv.Atoi(params[0]); err == nil {
			return t, true
		}

		return nil, false
	})

	tm.register(`(?:.+,\s*)*.+`, func(params []string, t reflect.Type) (interface{}, bool) {
		if t.Kind() != reflect.Slice {
			return nil, false
		}

		r := reflect.ValueOf(reflect.New(t).Elem().Interface())

		for _, i := range strings.Split(params[0], ",") {
			i = strings.TrimSpace(i)
			if e, ok := tm.convertParam(i, t.Elem()); ok {
				r = reflect.Append(r, e)
			} else {
				return nil, false
			}
		}

		return r.Interface(), true
	})
}

func (tm *transformMap) register(pattern string, transform StepArgumentTransform) {
	pattern = ensureCompleteMatch(pattern)

	p, err := regexp.Compile(pattern)

	if err != nil {
		panic(fmt.Sprintf("compiling transform regexp %q, %s", pattern, err))
	}

	(*tm)[p] = transform
}

func (tm *transformMap) convertParam(param string, target reflect.Type) (reflect.Value, bool) {
	for regex, tx := range *tm {
		params := regex.FindStringSubmatch(param)
		if params == nil {
			continue
		}

		f := reflect.ValueOf(tx)

		in := []reflect.Value{
			reflect.ValueOf(params),
			reflect.ValueOf(target),
		}

		out := f.Call(in)
		if out[1].Interface().(bool) {
			return reflect.ValueOf(out[0].Interface()), true
		}
	}

	return reflect.Value{}, false
}
