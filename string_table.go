package elicit

import "strings"

type stringTable [][]string

func (t *stringTable) hasColumn(cname string) bool {
	for _, c := range (*t)[0] {
		if c == cname {
			return true
		}
	}
	return false
}

func (t *stringTable) columnNameToIndexMap() map[string]int {
	m := make(map[string]int, len((*t)[0]))
	for i, c := range (*t)[0] {
		m[c] = i
	}
	return m
}

func (t *stringTable) hasParams(params []string) bool {
	for _, p := range params {
		pname := strings.TrimSuffix(strings.TrimPrefix(p, "<"), ">")
		if !t.hasColumn(pname) {
			return false
		}
	}
	return true
}
