package elicit_test

import (
	"fmt"
	"io/ioutil"
	"mmatt/elicit"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var tempdir string

func Test(t *testing.T) {
	elicit.New().
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

	steps["Create (?:a step definition|step definitions):"] =
		func(t *testing.T, text elicit.TextBlock) {
			createFile(t, "steps_test.go", fmt.Sprintf(stepFileFmt, "", text.Content))
		}

	steps["Create (?:a step definition|step definitions) using (.+):"] =
		func(t *testing.T, imports []string, text elicit.TextBlock) {
			createFile(t, "steps_test.go", fmt.Sprintf(stepFileFmt, strings.Join(imports, "\n"), text.Content))
		}

	steps["Running `(go test.*)` will output:"] =
		func(t *testing.T, command string, text elicit.TextBlock) {
			if err := os.Chdir(tempdir); err != nil {
				t.Fatalf("switching to tempdir %s: %s", tempdir, err)
			}

			parts := strings.Split(command, " ")
			output, _ := exec.Command(parts[0], parts[1:]...).CombinedOutput()

			expected, actual := quoteOutput(text.Content), quoteOutput(string(output))
			if !strings.Contains(actual, expected) {
				t.Errorf("\n\nExpected:\n\n%s\n\nto contain:\n\n%s\n", actual, expected)
			}
		}

	steps["`(.+)` will contain:"] =
		func(t *testing.T, filename string, text elicit.TextBlock) {
			path := filepath.Join(tempdir, filename)

			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Error(filename, err)
			}

			if contents, err := ioutil.ReadFile(path); err != nil {
				t.Error("reading", filename, err)
			} else {
				actual := string(contents)
				expected := strings.TrimSpace(text.Content)
				if actual != expected {
					t.Errorf("\n\nExpected:\n\n%s\n\nto equal:\n\n%s\n", quoteOutput(actual), quoteOutput(expected))
				}
			}
		}

	steps["Remove the temporary directory"] =
		func(t *testing.T) {
			t.Log("removing", tempdir)
			if err := os.RemoveAll(tempdir); err != nil {
				t.Errorf("removing tempdir %q: %s", tempdir, err)
			}
		}
}

func createTempDir(t *testing.T) {
	t.Log("creating tempdir...")
	var err error
	tempdir, err = ioutil.TempDir("", "elicit_test")

	if err != nil {
		t.Fatalf("creating tempdir: %s", err)
	}

	t.Log("tempdir:", tempdir)
}

func createFile(t *testing.T, filename, contents string) {
	if tempdir == "" {
		t.Fatal("creating file: tempdir not set")
	}

	t.Log("writing file:", filename, "\n", quoteOutput(contents))

	outpath := filepath.Join(tempdir, filename)
	ioutil.WriteFile(outpath, []byte(contents), 0777)
}

func quoteOutput(s string) string {
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`\033\[\d+(;\d+)?m`).ReplaceAllString(s, "")
	s = strings.Replace(s, " ", "·", -1)
	s = strings.Replace(s, "\t", "  ➟ ", -1)
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

const stepFileFmt = `
package elicit_test

import (
	"testing"
	%s
)

func init() {
	%s
}
`
