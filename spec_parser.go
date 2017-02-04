package elicit

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	bf "github.com/russross/blackfriday"
)

// implements the blackfriday Renderer interface.
type elicitSpecRenderer struct {
	spec         *specDef
	scenario     *scenarioDef
	step         *stepDef
	textTarget   *string
	lastText     string
	tableHeaders []string
	tableRow     []string
	tableRows    stringTable
}

type stringTable [][]string

type specDef struct {
	Path        string
	Name        string
	BeforeSteps []stepDef
	Scenarios   []scenarioDef
	AfterSteps  []stepDef
	Tables      []stringTable
}

type scenarioDef struct {
	Spec   *specDef
	Name   string
	Steps  []stepDef
	Tables []stringTable
}

type stepDef struct {
	Spec     *specDef
	Scenario *scenarioDef
	Text     string
	Params   []string
	Tables   []stringTable
}

func (specs *specCollection) parseSpecFolder(directory string) {
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".spec") {
			*specs = append(*specs, specDef{})
			spec := &(*specs)[len(*specs)-1]
			spec.loadFromFile(path)
		}
		return nil
	})
}

func (spec *specDef) loadFromFile(specFilePath string) *specDef {
	specText, err := ioutil.ReadFile(specFilePath)

	if err != nil {
		log.Fatalf("parsing spec file: %s: %s", specFilePath, err)
	}

	spec.Path = specFilePath

	r := &elicitSpecRenderer{spec: spec}
	md := bf.Markdown(specText, r, bf.EXTENSION_TABLES)

	fmt.Printf(string(md))

	return spec
}

// GetFlags not used
func (e *elicitSpecRenderer) GetFlags() int {
	return 0
}

// DocumentHeader not used
func (e *elicitSpecRenderer) DocumentHeader(out *bytes.Buffer) {
}

// BlockCode not used
func (e *elicitSpecRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
}

// TitleBlock not used
func (e *elicitSpecRenderer) TitleBlock(out *bytes.Buffer, text []byte) {
}

// BlockQuote not used
func (e *elicitSpecRenderer) BlockQuote(out *bytes.Buffer, text []byte) {
}

// BlockHtml not used
func (e *elicitSpecRenderer) BlockHtml(out *bytes.Buffer, text []byte) {
}

// Header creates test hierarchy
func (e *elicitSpecRenderer) Header(out *bytes.Buffer, text func() bool, level int, id string) {

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

func (s *specDef) createScenario() *scenarioDef {
	s.Scenarios = append(s.Scenarios, scenarioDef{Spec: s})
	return &s.Scenarios[len(s.Scenarios)-1]
}

// HRule escapes from the current scenario (i.e. subsequent steps appear in parent "scope")
func (e *elicitSpecRenderer) HRule(out *bytes.Buffer) {
	e.scenario = nil
	e.step = nil
}

// List wraps test steps (there's no way to specify an empty one)
func (e *elicitSpecRenderer) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()

	e.addStepToCurrentContext()

	if !text() {
		out.Truncate(marker)
	}

	e.removeLastStep()
	e.textTarget = nil
}

// ListItem creates a test step
func (e *elicitSpecRenderer) ListItem(out *bytes.Buffer, text []byte, flags int) {
	e.addStepToCurrentContext()
}

func (e *elicitSpecRenderer) addStepToCurrentContext() {
	if e.scenario != nil {
		e.step = e.scenario.createStep()
	} else if len(e.spec.Scenarios) == 0 {
		e.step = e.spec.createBeforeStep()
	} else {
		e.step = e.spec.createAfterStep()
	}
	e.textTarget = &e.step.Text
}

func (s *scenarioDef) createStep() *stepDef {
	s.Steps = append(s.Steps, stepDef{Spec: s.Spec, Scenario: s})
	return &s.Steps[len(s.Steps)-1]
}

func (s *specDef) createBeforeStep() *stepDef {
	s.BeforeSteps = append(s.BeforeSteps, stepDef{Spec: s})
	return &s.BeforeSteps[len(s.BeforeSteps)-1]
}

func (s *specDef) createAfterStep() *stepDef {
	s.AfterSteps = append(s.AfterSteps, stepDef{Spec: s})
	return &s.AfterSteps[len(s.AfterSteps)-1]
}

func (e *elicitSpecRenderer) removeLastStep() {
	var steps *[]stepDef
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
func (e *elicitSpecRenderer) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()

	if !text() {
		out.Truncate(marker)
	}
}

// Table adds the constructed table to the active context
func (e *elicitSpecRenderer) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
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
func (e *elicitSpecRenderer) TableRow(out *bytes.Buffer, text []byte) {
	if len(e.tableRow) > 0 {
		e.tableRows = append(e.tableRows, e.tableRow)
	}
	e.tableRow = nil
}

// TableHeaderCell defines a column in a table
func (e *elicitSpecRenderer) TableHeaderCell(out *bytes.Buffer, text []byte, align int) {
	e.tableRow = append(e.tableRow, e.lastText)
}

// TableCell adds a cell to the current row
func (e *elicitSpecRenderer) TableCell(out *bytes.Buffer, text []byte, align int) {
	e.tableRow = append(e.tableRow, e.lastText)
}

// Footnotes not used
func (e *elicitSpecRenderer) Footnotes(out *bytes.Buffer, text func() bool) {
	marker := out.Len()

	if !text() {
		out.Truncate(marker)
	}
}

// FootnoteItem not used
func (e *elicitSpecRenderer) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
}

// AutoLink output plaintext
func (e *elicitSpecRenderer) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	e.NormalText(out, link)
}

// CodeSpan output plaintext
func (e *elicitSpecRenderer) CodeSpan(out *bytes.Buffer, text []byte) {
	e.NormalText(out, text)
}

// DoubleEmphasis output plaintext
func (e *elicitSpecRenderer) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	e.NormalText(out, text)
}

// Emphasis output plaintext
func (e *elicitSpecRenderer) Emphasis(out *bytes.Buffer, text []byte) {
	e.NormalText(out, text)
}

// Image not used
func (e *elicitSpecRenderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
}

// LineBreak not used
func (e *elicitSpecRenderer) LineBreak(out *bytes.Buffer) {
}

// Link contents written as plaintext
func (e *elicitSpecRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	e.NormalText(out, content)
}

// RawHtmlTag represents a table-derived parameter
func (e *elicitSpecRenderer) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	e.NormalText(out, tag)
	if e.textTarget == &e.step.Text {
		e.step.Params = append(e.step.Params, e.lastText)
	}
}

// TripleEmphasis outputs plaintext
func (e *elicitSpecRenderer) TripleEmphasis(out *bytes.Buffer, text []byte) {
	e.NormalText(out, text)
}

// StrikeThrough output as plaintext
func (e *elicitSpecRenderer) StrikeThrough(out *bytes.Buffer, text []byte) {
	e.NormalText(out, text)
}

// FootnoteRef not used
func (e *elicitSpecRenderer) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
}

// Entity output as plaintext
func (e *elicitSpecRenderer) Entity(out *bytes.Buffer, entity []byte) {
	e.NormalText(out, entity)
}

// NormalText output as plaintext
func (e *elicitSpecRenderer) NormalText(out *bytes.Buffer, text []byte) {
	e.lastText = string(text[:])
	if len(strings.TrimSpace(e.lastText)) == 0 {
		return
	}
	if e.textTarget != nil {
		*e.textTarget += e.lastText
	}
}

// DocumentFooter not used
func (e *elicitSpecRenderer) DocumentFooter(out *bytes.Buffer) {
}

func (s *stepDef) resolveStepParams() []string {
	if len(s.Params) == 0 {
		return []string{s.Text}
	}

	table := stringTable{}
	resolved := []string{}
	found := false

	if s.Scenario != nil {
		table, found = s.findTableWithParams(s.Scenario.Tables)
	}

	if !found {
		table, found = s.findTableWithParams(s.Spec.Tables)
	}

	if found {
		m := table.columnNameToIndexMap()
		for _, row := range table[1:] {
			text := s.Text
			for _, p := range s.Params {
				pname := strings.TrimSuffix(strings.TrimPrefix(p, "<"), ">")
				pvalue := row[m[pname]]
				text = strings.Replace(text, p, pvalue, -1)
			}
			resolved = append(resolved, text)
		}
	}

	return resolved
}

func (s *stepDef) findTableWithParams(tables []stringTable) (stringTable, bool) {
	for _, t := range tables {
		if t.hasParams(s.Params) {
			return t, true
		}
	}
	return nil, false
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
