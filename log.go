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
	fmt.Fprintf(&l.buffer, "\n%s\n%s\n", s.name, strings.Repeat("=", len(s.name)))
}

func (l *log) scenario(s *scenario) {
	fmt.Fprintf(&l.buffer, "\n%s\n%s\n", s.name, strings.Repeat("-", len(s.name)))
}

func (l *log) step(s *step, text string) {
	var prefix, suffix string

	tables := len(s.tables)
	textBlocks := len(s.textBlocks)

	switch s.result {
	case undefined:
		prefix = "?"
	case skipped:
		prefix = "⤹"
	case failed:
		prefix = "✘"
	case panicked:
		prefix = "⚠"
	case passed:
		if s.forced {
			prefix = "✔"
		} else {
			prefix = "✓"
		}
	}

	if textBlocks > 0 {
		suffix += strings.Repeat(" ☰", textBlocks)
	}

	if tables > 0 {
		suffix += strings.Repeat(" ☷", tables)
	}

	fmt.Fprintf(&l.buffer, "  %s %s%s\n", prefix, text, suffix)
}
