package elicit_test

import (
	"flag"
	"io/ioutil"
	"mmatt/elicit"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	reportFile = flag.String("report", "", "Path to save the execution report, otherwise stdout")
)

// TODO(matt) store in scenario context
var tempdir string

func Test(t *testing.T) {
	elicit.New().
		WithReportPath(*reportFile).
		WithSpecsFolder("./specs").
		WithSteps(steps).
		RunTests(t)
}

var steps = map[string]interface{}{}

func init() {

	steps["Create a temporary directory"] =
		func(t *testing.T) {
			var err error
			tempdir, err = ioutil.TempDir("", "elicit_test")

			if err != nil {
				t.Fatalf("creating tempdir: %s", err)
			}
		}

	steps["Create a `(.*)` file:"] =
		func(t *testing.T, filename string, text elicit.TextBlock) {
			outpath := filepath.Join(tempdir, filename)
			ioutil.WriteFile(outpath, []byte(text.Content), 0777)
		}

	steps["Running `go test` will output:"] =
		func(t *testing.T, text elicit.TextBlock) {
			if err := os.Chdir(tempdir); err != nil {
				t.Fatalf("switching to tempdir %s: %s", tempdir, err)
			}

			output, _ := exec.Command("go", "test").CombinedOutput()

			expected, actual := quoteOutput(text.Content), quoteOutput(string(output))
			if !strings.Contains(actual, expected) {
				t.Errorf("\n\nExpected:\n\n%s\n\nto contain:\n\n%s\n", actual, expected)
			}
		}

	steps["Remove the temporary directory"] =
		func(t *testing.T) {
			if err := os.RemoveAll(tempdir); err != nil {
				t.Errorf("removing tempdir %q: %s", tempdir, err)
			}
		}
}

func quoteOutput(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, " ", "·", -1)
	s = strings.Replace(s, "\t", "➟", -1)
	s = "  | " + strings.Join(strings.Split(s, "\n"), "\n  | ")
	return s
}
