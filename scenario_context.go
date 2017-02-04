package elicit

import "bytes"

// ScenarioContext stores scenario specific information and data. It is reset for each scenario.
type ScenarioContext struct {
	name        string
	currentStep string
	skipped     bool
	failed      bool
	logbuf      bytes.Buffer
}
