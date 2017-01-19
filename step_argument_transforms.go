package elicit

import "strconv"
import "reflect"
import "strings"

func stringTransform(s string, t reflect.Type) (interface{}, bool) {
	if t.Kind() != reflect.String {
		return nil, false
	}

	return s, true
}

func intTransform(s string, t reflect.Type) (interface{}, bool) {
	if t != reflect.TypeOf((*int)(nil)).Elem() {
		return nil, false
	}

	if t, err := strconv.Atoi(s); err == nil {
		return t, true
	}

	return nil, false
}

func commaSliceTransform(s string, t reflect.Type) (interface{}, bool) {
	if t.Kind() != reflect.Slice {
		return nil, false
	}

	r := reflect.ValueOf(reflect.New(t).Elem().Interface())

	for _, i := range strings.Split(s, ",") {
		i = strings.TrimSpace(i)
		if e, ok := convertParam(i, t.Elem()); ok {
			r = reflect.Append(r, e)
		} else {
			return nil, false
		}
	}

	return r.Interface(), true
}
