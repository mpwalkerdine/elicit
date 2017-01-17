// GENERATED BY ELICIT - DO NOT EDIT

package specs_test

import (
	"mmatt/elicit"
	"testing"
)

func Test_Example_Sub_Spec(t *testing.T) {
	e := elicit.CurrentContext

	e.BeginSpecTest(t)

	t.Run("A Scenario", func(t *testing.T) {
		e.BeginScenarioTest(t)

		e.RunStep("Something Here")

		e.RunStep("Test Step 1")
		e.RunStep("Test Step 2")
	})
}
