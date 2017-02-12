package elicit

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type log struct {
	ctx     *Context
	buffer  bytes.Buffer
	outpath string
}

func (l *log) writeToFile() {

	for _, spec := range l.ctx.specs {
		l.writeSpecHeader(spec)
		for _, scenario := range spec.scenarios {
			l.writeScenarioHeader(scenario)
			for _, step := range scenario.steps {
				l.writeStepResult(step)
			}
		}
	}

	if l.outpath == "" {
		fmt.Println(l.buffer.String())
	} else {
		os.MkdirAll(filepath.Dir(l.outpath), 0755)
		ioutil.WriteFile(l.outpath, l.buffer.Bytes(), 0755)
	}
}

func (l *log) writeSpecHeader(s *spec) {
	fmt.Fprintf(&l.buffer, "\n\n%s\n%s\n", s.name, strings.Repeat("=", len(s.name)))
}

func (l *log) writeScenarioHeader(s *scenario) {
	f := "%s\n%s"

	switch s.result {
	case undefined:
		f = l.yellow(f)
	case skipped:
		f = l.blue(f)
	case failed, panicked:
		f = l.red(f)
	}

	f = "\n" + f + "\n"

	fmt.Fprintf(&l.buffer, f, s.name, strings.Repeat("-", len(s.name)))
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
	case undefined:
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
		if s.forced {
			prefix = l.green("✔")
			text = l.green(text)
		} else {
			prefix = l.green("✓")
		}
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
	if l.outpath == "" {
		s = fmt.Sprintf("\033[%dm%s\033[0m", colour, s)
	}
	return s
}
