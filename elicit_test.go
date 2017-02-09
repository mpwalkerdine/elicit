package elicit_test

import (
	"flag"
	"io/ioutil"
	"mmatt/elicit"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	reportFile = flag.String("report", "", "Path to save the execution report, otherwise stdout")
)

var tempdir string

func Test(t *testing.T) {
	elicit.New().
		WithReportPath(*reportFile).
		WithSpecsFolder("./specs").
		WithSteps(steps).
		WithTransforms(transforms).
		RunTests(t)
}

var steps = map[string]interface{}{}
var transforms = map[string]elicit.StepArgumentTransform{}

func init() {

	steps["Create a temporary directory"] = createTempDir

	steps["Create a temporary environment"] =
		func(t *testing.T) {
			createTempDir(t)
			createFile(t, "specs_test.go", testfile)
		}

	steps["Create a `(.*)` file:"] =
		func(t *testing.T, filename string, text elicit.TextBlock) {
			createFile(t, filename, text.Content)
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

func createTempDir(t *testing.T) {
	var err error
	tempdir, err = ioutil.TempDir("", "elicit_test")

	if err != nil {
		t.Fatalf("creating tempdir: %s", err)
	}
}

func createFile(t *testing.T, filename, contents string) {
	outpath := filepath.Join(tempdir, filename)
	ioutil.WriteFile(outpath, []byte(contents), 0777)
}

func quoteOutput(s string) string {
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`\033\[\d+(;\d+)?m`).ReplaceAllString(s, "")
	s = strings.Replace(s, " ", "·", -1)
	s = strings.Replace(s, "\t", "➟", -1)
	s = "  | " + strings.Join(strings.Split(s, "\n"), "\n  | ")
	return s
}

const testfile = `
package elicit_test

import (
    "mmatt/elicit"
    "testing"
)

func Test(t *testing.T) {
    elicit.New().
        WithSpecsFolder(".").
        WithSteps(steps).
        WithTransforms(transforms).
        RunTests(t)
}

var steps = map[string]interface{}{}
var transforms = map[string]elicit.StepArgumentTransform{}
`
