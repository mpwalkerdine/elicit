package elicit

import (
	"bytes"
	"fmt"
	"strings"
)

type log struct {
	buffer bytes.Buffer
}

func (l *log) String() string {
	return l.buffer.String()
}

func (l *log) spec(s *spec) {
	fmt.Fprintf(&l.buffer, "\n\n%s\n%s\n", s.name, strings.Repeat("=", len(s.name)))
}

func (l *log) scenario(s *scenario) {
	fmt.Fprintf(&l.buffer, "\n%s\n%s\n", s.name, strings.Repeat("-", len(s.name)))
}

func (l *log) step(s *step, text string) {
	var prefix, suffix string

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
	case skipped:
		prefix = l.yellow("⤹")
	case failed:
		prefix = l.red("✘")
	case panicked:
		prefix = l.red("⚡")
	case passed:
		if s.forced {
			prefix = l.green("✔")
		} else {
			prefix = l.green("✓\033[0m")
		}
	}

	fmt.Fprintf(&l.buffer, "  %s %s%s\033[0m\n", prefix, text, suffix)
}

func (l *log) red(s string) string {
	return fmt.Sprintf("\033[1;31m%s", s)
}

func (l *log) green(s string) string {
	return fmt.Sprintf("\033[0;32m%s", s)
}

func (l *log) yellow(s string) string {
	return fmt.Sprintf("\033[0;33m%s", s)
}
