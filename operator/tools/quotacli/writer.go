package main

import (
	"fmt"
	"io"

	"k8s.io/cli-runtime/pkg/printers"
)

type flusher interface {
	Flush()
}

// PrefixWriter can write text at various indentation levels.
type PrefixWriter interface {
	// Write writes text with the specified indentation level.
	Write(level int, format string, a ...interface{})
	// WriteLine writes an entire line with no indentation level.
	WriteLine(a ...interface{})
	// Flush forces indentation to be reset.
	Flush()
}

// prefixWriter implements PrefixWriter
type prefixWriter struct {
	out io.Writer
}

var _ PrefixWriter = &prefixWriter{}

func NewPrefixWriter(out io.Writer) PrefixWriter {
	return &prefixWriter{out: out}
}

func (pw *prefixWriter) Write(level int, format string, a ...interface{}) {
	levelSpace := "  "
	prefix := ""
	for i := 0; i < level; i++ {
		prefix += levelSpace
	}
	output := fmt.Sprintf(prefix+format, a...)
	printers.WriteEscaped(pw.out, output)
}

func (pw *prefixWriter) WriteLine(a ...interface{}) {
	output := fmt.Sprintln(a...)
	printers.WriteEscaped(pw.out, output)
}

func (pw *prefixWriter) Flush() {
	if f, ok := pw.out.(flusher); ok {
		f.Flush()
	}
}
