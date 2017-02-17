package elicit

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type log struct {
	ctx       *Context
	buffer    bytes.Buffer
	useColour bool
	outpath   string
}

func (l *log) writeToConsole() {
	l.useColour = true
	l.fillBuffer(false)
	if l.buffer.Len() > 0 {
		fmt.Println(l.buffer.String())
	}
}

func (l *log) writeToFile() {
	if l.outpath == "" {
		return
	}

	if err := os.MkdirAll(filepath.Dir(l.outpath), 0755); err != nil {
		panic(err)
	}

	l.useColour = false
	l.fillBuffer(true)

	if err := ioutil.WriteFile(l.outpath, bytes.TrimSpace(l.buffer.Bytes()), 0755); err != nil {
		panic(err)
	}
}

func (l *log) fillBuffer(forceVerbose bool) {
	l.buffer.Truncate(0)

	verbose := forceVerbose
	if v := flag.Lookup("test.v"); !forceVerbose && v != nil {
		verbose = v.Value.String() == "true"
	}

	for _, spec := range l.ctx.specs {
		if !verbose && !spec.result.shouldLog() {
			continue
		}
		l.writeSpecHeader(spec)

		for _, scenario := range spec.scenarios {
			if !verbose && !scenario.result.shouldLog() {
				continue
			}
			l.writeScenarioHeader(scenario)

			for _, step := range scenario.steps {
				l.writeStepResult(step)
			}

		}
	}
}

func (l *log) writeLn() {
	fmt.Fprintln(&l.buffer)
}

func (l *log) writeSpecHeader(s *spec) {
	name := s.name
	underline := strings.Repeat("=", len(s.name))
	resultCounts := [numResultTypes]int{}

	switch s.result {
	case pending:
		name = l.yellow(name)
		underline = l.yellow(underline)
	case skipped:
		name = l.blue(name)
		underline = l.blue(underline)
	case failed, panicked:
		name = l.red(name)
		underline = l.red(underline)
	}

	for _, scenario := range s.scenarios {
		resultCounts[scenario.result]++
	}

	resultString := ""

	for i, count := range resultCounts {
		if count > 0 {
			resultString += fmt.Sprintf("\n%s: %d", result(i), count)
		}
	}

	fmt.Fprintf(&l.buffer, "\n\n%s\n%s%s\n", name, underline, resultString)
}

func (l *log) writeScenarioHeader(s *scenario) {
	name := s.name
	underline := strings.Repeat("-", len(s.name))

	switch s.result {
	case pending:
		name = l.yellow(name)
		underline = l.yellow(underline)
	case skipped:
		name = l.blue(name)
		underline = l.blue(underline)
	case failed, panicked:
		name = l.red(name)
		underline = l.red(underline)
	}

	fmt.Fprintf(&l.buffer, "\n%s\n%s\n%s\n\n", name, underline, s.result)
}

func (l *log) writeStepResult(s *step) {
	var prefix, suffix string
	text := s.text

	textBlocks := len(s.textBlocks)
	if textBlocks > 0 {
		suffix += strings.Repeat(" ☰", textBlocks)
	}

	tables := len(s.tables)
	if tables > 0 {
		suffix += strings.Repeat(" ☷", tables)
	}

	switch s.result {
	case pending:
		prefix = l.yellow("?")
		text = l.yellow(text)
	case skipped:
		prefix = l.blue("⤹")
		text = l.blue(text)
	case failed:
		prefix = l.red("✘")
		text = l.red(text)
	case panicked:
		prefix = l.red("⚡")
		text = l.red(text)
	case passed:
		prefix = l.green("✓")
	}

	fmt.Fprintf(&l.buffer, "    %s %s%s\n", prefix, text, suffix)

	if s.log.Len() > 0 {
		leftPad := "        "
		stepLog := s.log.String()
		stepLog = strings.TrimSuffix(stepLog, "\n")
		lines := strings.Split(stepLog, "\n")
		stepLog = leftPad + strings.Join(lines, "\n"+leftPad)
		fmt.Fprintln(&l.buffer, stepLog)
	}
}

func (l *log) red(s string) string {
	return l.colour(s, 31)
}

func (l *log) green(s string) string {
	return l.colour(s, 32)
}

func (l *log) yellow(s string) string {
	return l.colour(s, 33)
}

func (l *log) blue(s string) string {
	return l.colour(s, 34)
}

func (l *log) colour(s string, colour int) string {
	if l.useColour {
		s = fmt.Sprintf("\033[%dm%s\033[0m", colour, s)
	}
	return s
}
