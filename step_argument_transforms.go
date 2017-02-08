package elicit

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// StepArgumentTransform transforms captured groups in the step pattern to a function parameter type
// Note that if the actual string cannot be converted to the target type by the transform, it should return false
type StepArgumentTransform func([]string, reflect.Type) (interface{}, bool)

type transformMap map[*regexp.Regexp]StepArgumentTransform

func (tm *transformMap) register(pattern string, transform StepArgumentTransform) {
	pattern = ensureCompleteMatch(pattern)

	p, err := regexp.Compile(pattern)

	if err != nil {
		panic(fmt.Sprintf("compiling transform regexp %q, %s", pattern, err))
	}

	(*tm)[p] = transform
}

func stringTransform(params []string, t reflect.Type) (interface{}, bool) {
	if t != reflect.TypeOf((*string)(nil)).Elem() {
		return nil, false
	}

	return params[0], true
}

func intTransform(params []string, t reflect.Type) (interface{}, bool) {
	if t != reflect.TypeOf((*int)(nil)).Elem() {
		return nil, false
	}

	if t, err := strconv.Atoi(params[0]); err == nil {
		return t, true
	}

	return nil, false
}

// func commaSliceTransform(ctx *Context, s string, t reflect.Type) (interface{}, bool) {
// 	if t.Kind() != reflect.Slice {
// 		return nil, false
// 	}
//
// 	r := reflect.ValueOf(reflect.New(t).Elem().Interface())
//
// 	for _, i := range strings.Split(s, ",") {
// 		i = strings.TrimSpace(i)
// 		if e, ok := ctx.convertParam(i, t.Elem()); ok {
// 			r = reflect.Append(r, e)
// 		} else {
// 			return nil, false
// 		}
// 	}
//
// 	return r.Interface(), true
// }
