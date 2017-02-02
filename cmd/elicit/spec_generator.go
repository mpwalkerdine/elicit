package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

type specGenerator struct {
	buf       bytes.Buffer
	specsRoot string
	spec      *specFile
	outPkg    string
}

func (g *specGenerator) write(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *specGenerator) writeln(format string, args ...interface{}) {
	g.write(format+"\n", args...)
}

func (g *specGenerator) generate() {

	g.writeHeaders()

	for _, scenario := range g.spec.Scenarios {
		g.writeScenarioTestRun(scenario)
	}

	g.writeFooters()

	g.writeSpecTestToFile()
}

func (g *specGenerator) writeHeaders() {
	g.writeln("// GENERATED BY ELICIT - DO NOT EDIT\n")

	pkgName := filepath.Base(g.outPkg)
	g.writeln("package %s", pkgName)

	g.writeln("import (")
	g.writeln("%q", "mmatt/elicit")
	g.writeln("%q", "testing")
	g.writeln(")")

	g.writeln("func Test_%s(t *testing.T) {", escapeIdentifier(g.spec.Name))
	g.writeln("e := elicit.CurrentContext\n")
	g.writeln("e.BeginSpecTest(%q)", g.spec.Name)
}

func (g *specGenerator) writeScenarioTestRun(scenario scenario) {
	g.writeln("\nt.Run(%q, func(t *testing.T) {", scenario.Name)
	g.writeln("e.BeginScenarioTest(%q)", scenario.Name)

	g.doubleSpace(len(g.spec.BeforeSteps) > 0)
	for _, before := range g.spec.BeforeSteps {
		g.writeln("e.RunStep(%q)", before.Text)
	}

	g.doubleSpace(len(scenario.Steps) > 0)
	for _, step := range scenario.Steps {
		g.write("e.RunStep(%q,", step.Text)
		for _, stepTable := range step.Tables {
			g.writeStepTable(stepTable)
		}
		g.writeln(")")
	}

	g.doubleSpace(len(g.spec.AfterSteps) > 0)
	for _, after := range g.spec.AfterSteps {
		g.writeln("e.RunStep(%q)", after.Text)
	}

	g.writeScenarioEndBlock()
	g.writeln("})")
}

func (g *specGenerator) writeStepTable(table [][]string) {
	g.writeln("\n[][]string{")
	g.writeTableRows(table)
	g.writeln("},")
}

func (g *specGenerator) writeTableRows(table [][]string) {
	for _, row := range table {
		g.write("[]string{")
		for _, v := range row {
			g.write("%q, ", v)
		}
		g.writeln("},")
	}
}

func (g *specGenerator) writeScenarioEndBlock() {
	g.writeln("\nif r, l := e.EndScenarioTest(); r == elicit.Failed {")
	g.writeln("t.Errorf(l)")
	g.writeln("} else if r == elicit.Skipped {")
	g.writeln("t.Skipf(l)")
	g.writeln("}")
}

func (g *specGenerator) writeFooters() {
	g.writeln("}")
}

func (g *specGenerator) doubleSpace(condition bool) {
	if condition {
		g.writeln("\n")
	}
}

func escapeIdentifier(raw string) string {
	p, err := regexp.Compile(`[^a-zA-Z0-9_]`)

	if err != nil {
		log.Fatalf("parsing regex: %s", err)
	}

	return p.ReplaceAllLiteralString(raw, "_")
}

func (g *specGenerator) writeSpecTestToFile() {
	specRelPath, err := filepath.Rel(g.specsRoot, g.spec.Path)
	if err != nil {
		log.Fatalf("determing relative path from %q to %q: %s", g.specsRoot, g.spec.Path, err)
	}

	trimmedName := strings.TrimSuffix(strings.TrimLeft(specRelPath, "./\\"), ".spec")
	flattenedName := strings.Replace(trimmedName, string(filepath.Separator), "_", -1)
	outname := fmt.Sprintf("%s_test.go", flattenedName)
	outpath := filepath.Join(g.outPkg, outname)

	src := formatSource(g.buf.Bytes())

	if err := ioutil.WriteFile(outpath, src, 0644); err != nil {
		log.Fatalf("writing output to %q: %s", outpath, err)
	}
}

func formatSource(text []byte) []byte {
	src, err := format.Source(text)
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return text
	}
	return src
}
