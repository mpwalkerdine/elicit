package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"strings"

	bf "github.com/russross/blackfriday"
)

type specFile struct {
	Path        string
	Name        string
	BeforeSteps []step
	Scenarios   []scenario
	AfterSteps  []step
	Tables      [][][]string
}

type scenario struct {
	Name   string
	Steps  []step
	Tables [][][]string
}

type step struct {
	Text   string
	Tables [][][]string
}

func (s *specFile) createScenario() *scenario {
	s.Scenarios = append(s.Scenarios, scenario{})
	return &s.Scenarios[len(s.Scenarios)-1]
}

func (s *specFile) createBeforeStep() *step {
	s.BeforeSteps = append(s.BeforeSteps, step{})
	return &s.BeforeSteps[len(s.BeforeSteps)-1]
}

func (s *specFile) createAfterStep() *step {
	s.AfterSteps = append(s.AfterSteps, step{})
	return &s.AfterSteps[len(s.AfterSteps)-1]
}

func (s *scenario) createStep() *step {
	s.Steps = append(s.Steps, step{})
	return &s.Steps[len(s.Steps)-1]
}

func parseSpecFile(specFilePath string) *specFile {
	specText, err := ioutil.ReadFile(specFilePath)

	if err != nil {
		log.Fatalf("parsing spec file: %s: %s", specFilePath, err)
	}

	spec := &specFile{
		Path: specFilePath,
	}

	r := &elicitTest{spec: spec}
	md := bf.Markdown(specText, r, bf.EXTENSION_TABLES)

	fmt.Printf(string(md))

	return spec
}

// elicitTest is a type that implements the blackfriday Renderer interface.
type elicitTest struct {
	spec         *specFile
	scenario     *scenario
	step         *step
	textTarget   *string
	lastText     string
	tableHeaders []string
	tableRow     []string
	tableRows    [][]string
}

// GetFlags not used
func (e *elicitTest) GetFlags() int {
	return 0
}

// DocumentHeader not used
func (e *elicitTest) DocumentHeader(out *bytes.Buffer) {
}

// BlockCode not used
func (e *elicitTest) BlockCode(out *bytes.Buffer, text []byte, lang string) {
}

// TitleBlock not used
func (e *elicitTest) TitleBlock(out *bytes.Buffer, text []byte) {
}

// BlockQuote not used
func (e *elicitTest) BlockQuote(out *bytes.Buffer, text []byte) {
}

// BlockHtml not used
func (e *elicitTest) BlockHtml(out *bytes.Buffer, text []byte) {
}

// Header creates test hierarchy
func (e *elicitTest) Header(out *bytes.Buffer, text func() bool, level int, id string) {

	switch level {
	case 1: // Spec Name
		e.step = nil
		e.scenario = nil
		e.textTarget = &e.spec.Name
	case 2:
		e.step = nil
		e.scenario = e.spec.createScenario()
		e.textTarget = &e.scenario.Name
	}

	marker := out.Len()

	if !text() {
		out.Truncate(marker)
		return
	}

	e.textTarget = nil
}

// HRule escapes from the current scenario (i.e. subsequent steps appear in parent "scope")
func (e *elicitTest) HRule(out *bytes.Buffer) {
	e.scenario = nil
	e.step = nil
}

// List wraps test steps (there's no way to specify an empty one)
func (e *elicitTest) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()

	e.addStepToCurrentContext()

	if !text() {
		out.Truncate(marker)
	}

	e.removeLastStep()
	e.textTarget = nil
}

// ListItem creates a test step
func (e *elicitTest) ListItem(out *bytes.Buffer, text []byte, flags int) {
	e.addStepToCurrentContext()
}

func (e *elicitTest) addStepToCurrentContext() {
	if e.scenario != nil {
		e.step = e.scenario.createStep()
	} else if len(e.spec.Scenarios) == 0 {
		e.step = e.spec.createBeforeStep()
	} else {
		e.step = e.spec.createAfterStep()
	}
	e.textTarget = &e.step.Text
}

func (e *elicitTest) removeLastStep() {
	var steps *[]step
	if e.scenario != nil {
		steps = &e.scenario.Steps
	} else if len(e.spec.Scenarios) == 0 {
		steps = &e.spec.BeforeSteps
	} else {
		steps = &e.spec.AfterSteps
	}
	*steps = (*steps)[:len(*steps)-1]
	e.step = &(*steps)[len(*steps)-1]
}

// Paragraph not used
func (e *elicitTest) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()

	if !text() {
		out.Truncate(marker)
	}
}

// Table adds the constructed table to the active context
func (e *elicitTest) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	if e.step != nil {
		e.step.Tables = append(e.step.Tables, e.tableRows)
	} else if e.scenario != nil {
		e.scenario.Tables = append(e.scenario.Tables, e.tableRows)
	} else {
		e.spec.Tables = append(e.spec.Tables, e.tableRows)
	}

	e.tableRows = nil
}

// TableRow saves the current row
func (e *elicitTest) TableRow(out *bytes.Buffer, text []byte) {
	if len(e.tableRow) > 0 {
		e.tableRows = append(e.tableRows, e.tableRow)
	}
	e.tableRow = nil
}

// TableHeaderCell defines a column in a table
func (e *elicitTest) TableHeaderCell(out *bytes.Buffer, text []byte, align int) {
	e.tableRow = append(e.tableRow, e.lastText)
}

// TableCell adds a cell to the current row
func (e *elicitTest) TableCell(out *bytes.Buffer, text []byte, align int) {
	e.tableRow = append(e.tableRow, e.lastText)
}

// Footnotes not used
func (e *elicitTest) Footnotes(out *bytes.Buffer, text func() bool) {
	marker := out.Len()

	if !text() {
		out.Truncate(marker)
	}
}

// FootnoteItem not used
func (e *elicitTest) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
}

// AutoLink output plaintext
func (e *elicitTest) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	e.NormalText(out, link)
}

// CodeSpan output plaintext
func (e *elicitTest) CodeSpan(out *bytes.Buffer, text []byte) {
	e.NormalText(out, text)
}

// DoubleEmphasis output plaintext
func (e *elicitTest) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	e.NormalText(out, text)
}

// Emphasis output plaintext
func (e *elicitTest) Emphasis(out *bytes.Buffer, text []byte) {
	e.NormalText(out, text)
}

// Image not used
func (e *elicitTest) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
}

// LineBreak not used
func (e *elicitTest) LineBreak(out *bytes.Buffer) {
}

// Link contents written as plaintext
func (e *elicitTest) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	e.NormalText(out, content)
}

// RawHtmlTag written as plaintext
func (e *elicitTest) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	e.NormalText(out, tag)
}

// TripleEmphasis outputs plaintext
func (e *elicitTest) TripleEmphasis(out *bytes.Buffer, text []byte) {
	e.NormalText(out, text)
}

// StrikeThrough output as plaintext
func (e *elicitTest) StrikeThrough(out *bytes.Buffer, text []byte) {
	e.NormalText(out, text)
}

// FootnoteRef not used
func (e *elicitTest) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
}

// Entity output as plaintext
func (e *elicitTest) Entity(out *bytes.Buffer, entity []byte) {
	e.NormalText(out, entity)
}

// NormalText output as plaintext
func (e *elicitTest) NormalText(out *bytes.Buffer, text []byte) {
	e.lastText = string(text[:])
	if len(strings.TrimSpace(e.lastText)) == 0 {
		return
	}
	if e.textTarget != nil {
		*e.textTarget += e.lastText
	}
}

// DocumentFooter not used
func (e *elicitTest) DocumentFooter(out *bytes.Buffer) {
}
