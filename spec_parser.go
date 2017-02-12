package elicit

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"regexp"

	bf "github.com/russross/blackfriday"
)

// implements the blackfriday Renderer interface.
type specParser struct {
	context         *Context
	currentPath     string
	currentSpec     *spec
	currentScenario *scenario
	currentStep     *step
	beforeSteps     []*step
	afterSteps      []*step
	textTarget      *string
	lastText        string
	tableHeaders    []string
	tableRow        []string
	tableRows       stringTable
}

func (p *specParser) parseSpecFolder(directory string) {
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".spec") {
			p.currentPath = path
			p.loadFromFile()
		}
		return nil
	})
}

func (p *specParser) loadFromFile() {
	specText, err := ioutil.ReadFile(p.currentPath)

	if err != nil {
		panic(fmt.Errorf("parsing spec file: %s: %s", p.currentPath, err))
	}

	// Strip out non-step items so they're not parsed
	specText = regexp.MustCompile(`(?m)^[-*] `).ReplaceAllLiteral(specText, []byte{})

	// Parse
	bf.Markdown(specText, p, bf.EXTENSION_TABLES|bf.EXTENSION_FENCED_CODE)

	p.closeSpec()
}

func (p *specParser) createSpec() {
	p.closeSpec()

	p.context.specs = append(p.context.specs, &spec{
		context: p.context,
		path:    p.currentPath,
	})
	p.currentSpec = p.context.specs[len(p.context.specs)-1]
	p.textTarget = &p.currentSpec.name
}

func (p *specParser) createScenario() {
	p.closeScenario()
	p.currentSpec.scenarios = append(p.currentSpec.scenarios, &scenario{
		context: p.context,
		spec:    p.currentSpec,
	})
	p.currentScenario = p.currentSpec.scenarios[len(p.currentSpec.scenarios)-1]
	p.textTarget = &p.currentScenario.name
}

func (p *specParser) createStep() {
	p.closeStep()

	if p.currentScenario != nil {
		p.currentScenario.steps = append(p.currentScenario.steps, &step{
			context:  p.context,
			spec:     p.currentSpec,
			scenario: p.currentScenario,
		})
		p.currentStep = p.currentScenario.steps[len(p.currentScenario.steps)-1]
	} else if len(p.currentSpec.scenarios) == 0 {
		p.beforeSteps = append(p.beforeSteps, &step{
			context: p.context,
			spec:    p.currentSpec,
		})
		p.currentStep = p.beforeSteps[len(p.beforeSteps)-1]
	} else {
		p.afterSteps = append(p.afterSteps, &step{
			context: p.context,
			spec:    p.currentSpec,
		})
		p.currentStep = p.afterSteps[len(p.afterSteps)-1]
	}
	p.textTarget = &p.currentStep.text
}

func (p *specParser) resolveStepParams(s *step) []string {
	table := stringTable{}
	resolved := []string{}
	found := false

	if p.currentScenario != nil {
		table, found = p.findTableWithStepParams(s, p.currentScenario.tables)
	}

	if !found {
		table, found = p.findTableWithStepParams(s, p.currentSpec.tables)
	}

	if !found {
		return resolved
	}

	m := table.columnNameToIndexMap()
	for _, row := range table[1:] {
		text := s.text
		for _, p := range s.params {
			pname := strings.TrimSuffix(strings.TrimPrefix(p, "<"), ">")
			pvalue := row[m[pname]]
			text = strings.Replace(text, p, pvalue, -1)
		}
		resolved = append(resolved, text)
	}

	return resolved
}

func (p *specParser) findTableWithStepParams(s *step, tables []stringTable) (stringTable, bool) {
	for _, t := range tables {
		if t.hasParams(s.params) {
			return t, true
		}
	}
	return nil, false
}

func (p *specParser) closeSpec() {
	if p.currentSpec == nil {
		return
	}

	p.closeScenario()

	for _, scenario := range p.currentSpec.scenarios {
		before := make([]*step, 0, len(p.beforeSteps))
		for _, b := range p.beforeSteps {
			ns := *b
			ns.scenario = scenario
			before = append(before, &ns)
		}
		scenario.steps = append(before, scenario.steps...)

		for _, a := range p.afterSteps {
			ns := *a
			ns.scenario = scenario
			scenario.steps = append(scenario.steps, &ns)
		}
	}

	p.beforeSteps = nil
	p.afterSteps = nil
	p.currentSpec = nil
}

func (p *specParser) closeScenario() {
	if p.currentScenario == nil {
		return
	}

	p.closeStep()
	p.currentScenario = nil
}

func (p *specParser) closeStep() {
	if p.currentStep == nil {
		return
	}

	if len(p.currentStep.params) > 0 {
		s := p.removeLastStep()
		for _, text := range p.resolveStepParams(s) {
			p.createStep()
			p.currentStep.text = text
		}
	}

	p.currentStep = nil
}

// GetFlags not used
func (p *specParser) GetFlags() int {
	return 0
}

// DocumentHeader not used
func (p *specParser) DocumentHeader(out *bytes.Buffer) {
}

// BlockCode not used
func (p *specParser) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	if p.currentStep != nil {
		p.currentStep.textBlocks = append(p.currentStep.textBlocks, TextBlock{Language: lang, Content: string(text[:])})
	}
}

// TitleBlock not used
func (p *specParser) TitleBlock(out *bytes.Buffer, text []byte) {
}

// BlockQuote not used
func (p *specParser) BlockQuote(out *bytes.Buffer, text []byte) {

}

// BlockHtml not used
func (p *specParser) BlockHtml(out *bytes.Buffer, text []byte) {
}

// Header creates test hierarchy
func (p *specParser) Header(out *bytes.Buffer, text func() bool, level int, id string) {

	switch level {
	case 1:
		p.createSpec()
	case 2:
		p.createScenario()
	}

	marker := out.Len()

	if !text() {
		out.Truncate(marker)
		return
	}

	p.textTarget = nil
}

// HRule escapes from the current scenario (i.e. subsequent steps appear in parent "scope")
func (p *specParser) HRule(out *bytes.Buffer) {
	p.closeScenario()
}

// List wraps test steps (there's no way to specify an empty one)
func (p *specParser) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()

	p.createStep()

	if !text() {
		out.Truncate(marker)
	}

	p.removeLastStep()
	p.textTarget = nil
}

// ListItem creates a test step
func (p *specParser) ListItem(out *bytes.Buffer, text []byte, flags int) {
	p.createStep()
}

func (p *specParser) removeLastStep() *step {
	var steps *[]*step
	if p.currentScenario != nil {
		steps = &p.currentScenario.steps
	} else if len(p.currentSpec.scenarios) == 0 {
		steps = &p.beforeSteps
	} else {
		steps = &p.afterSteps
	}
	lastStep := (*steps)[len(*steps)-1]
	*steps = (*steps)[:len(*steps)-1]

	if len(*steps) == 0 {
		p.currentStep = nil
	} else {
		p.currentStep = (*steps)[len(*steps)-1]
	}

	return lastStep
}

// Paragraph text prevents association of tables and code blocks with step
func (p *specParser) Paragraph(out *bytes.Buffer, text func() bool) {
	p.closeStep()
	marker := out.Len()

	if !text() {
		out.Truncate(marker)
	}
}

// Table adds the constructed table to the active context
func (p *specParser) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	if p.currentStep != nil {
		p.currentStep.tables = append(p.currentStep.tables, p.tableRows)
	} else if p.currentScenario != nil {
		p.currentScenario.tables = append(p.currentScenario.tables, p.tableRows)
	} else {
		p.currentSpec.tables = append(p.currentSpec.tables, p.tableRows)
	}

	p.tableRows = nil
}

// TableRow saves the current row
func (p *specParser) TableRow(out *bytes.Buffer, text []byte) {
	if len(p.tableRow) > 0 {
		p.tableRows = append(p.tableRows, p.tableRow)
	}
	p.tableRow = nil
}

// TableHeaderCell defines a column in a table
func (p *specParser) TableHeaderCell(out *bytes.Buffer, text []byte, align int) {
	p.tableRow = append(p.tableRow, p.lastText)
}

// TableCell adds a cell to the current row
func (p *specParser) TableCell(out *bytes.Buffer, text []byte, align int) {
	p.tableRow = append(p.tableRow, p.lastText)
}

// Footnotes not used
func (p *specParser) Footnotes(out *bytes.Buffer, text func() bool) {
	marker := out.Len()

	if !text() {
		out.Truncate(marker)
	}
}

// FootnoteItem not used
func (p *specParser) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
}

// AutoLink output plaintext
func (p *specParser) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	p.NormalText(out, link)
}

// CodeSpan output plaintext
func (p *specParser) CodeSpan(out *bytes.Buffer, text []byte) {
	s := "`" + string(text) + "`"
	p.WriteText(s)
}

// DoubleEmphasis output plaintext
func (p *specParser) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	p.NormalText(out, text)
}

// Emphasis output plaintext
func (p *specParser) Emphasis(out *bytes.Buffer, text []byte) {
	p.NormalText(out, text)
	if p.currentStep != nil && p.textTarget == &p.currentStep.text {
		p.currentStep.force = true
	}
}

// Image not used
func (p *specParser) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
}

// LineBreak not used
func (p *specParser) LineBreak(out *bytes.Buffer) {
}

// Link contents written as plaintext
func (p *specParser) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	p.NormalText(out, content)
}

// RawHtmlTag represents a table-derived parameter
func (p *specParser) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	p.NormalText(out, tag)
	if p.currentStep != nil && p.textTarget == &p.currentStep.text {
		p.currentStep.params = append(p.currentStep.params, p.lastText)
	}
}

// TripleEmphasis outputs plaintext
func (p *specParser) TripleEmphasis(out *bytes.Buffer, text []byte) {
	p.NormalText(out, text)
}

// StrikeThrough output as plaintext
func (p *specParser) StrikeThrough(out *bytes.Buffer, text []byte) {
	p.NormalText(out, text)
}

// FootnoteRef not used
func (p *specParser) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
}

// Entity output as plaintext
func (p *specParser) Entity(out *bytes.Buffer, entity []byte) {
	p.NormalText(out, entity)
}

// NormalText output as plaintext
func (p *specParser) NormalText(out *bytes.Buffer, text []byte) {
	p.WriteText(string(text[:]))
}

func (p *specParser) WriteText(s string) {
	p.lastText = strings.Replace(s, "\n", " ", -1)

	if len(strings.TrimSpace(p.lastText)) == 0 {
		return
	}

	if p.textTarget != nil {
		*p.textTarget += p.lastText
	}
}

// DocumentFooter not used
func (p *specParser) DocumentFooter(out *bytes.Buffer) {
}
