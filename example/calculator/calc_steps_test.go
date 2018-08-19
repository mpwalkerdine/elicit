package calculator_test

import (
	"strconv"
	"testing"

	"github.com/mpwalkerdine/elicit"
	sut "github.com/mpwalkerdine/elicit/example/calculator"
)

type operation struct {
	left   int
	symbol rune
	right  int
}

var steps = elicit.Steps{}
var transforms = elicit.Transforms{}

func init() {
	transforms["`(\\d+)\\s*([+-])\\s*(\\d+)`"] =
		func(params []string) operation {
			left, _ := strconv.Atoi(params[1])
			op := []rune(params[2])[0]
			right, _ := strconv.Atoi(params[3])

			return operation{
				left:   left,
				symbol: op,
				right:  right,
			}
		}

	steps["When (.+) is entered the answer is `(\\d+)`"] =
		func(t *testing.T, op operation, want int) {
			calc := sut.Calculator{}
			var got int
			switch op.symbol {
			case '+':
				got = calc.Add(op.left, op.right)
			case '-':
				got = calc.Sub(op.left, op.right)
			}
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		}

}
