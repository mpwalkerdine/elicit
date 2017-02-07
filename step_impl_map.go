package elicit

import (
	"fmt"
	"regexp"
	"strings"
)

type stepImplMap map[*regexp.Regexp]interface{}

func (sim *stepImplMap) register(pattern string, stepFunc interface{}) {

	pattern = strings.TrimSpace(pattern)
	pattern = ensureCompleteMatch(pattern)

	p, err := regexp.Compile(pattern)

	if err != nil {
		panic(fmt.Sprintf("compiling step regexp %q, %s", pattern, err))
	}

	// TODO(matt) check the pattern captures the correct number of parameters

	(*sim)[p] = stepFunc
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
