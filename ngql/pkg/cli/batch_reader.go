package cli

import (
	"bufio"
	"io"
	"strings"
)

// BatchReader interface defines methods for handling non-interactive input
type BatchReader interface {
	// ReadInput reads input and returns the processed line content
	// Returns: processed line content, whether EOF is encountered, whether more input is needed, error information
	ReadInput() (string, bool, bool, error)
	// Close releases associated resources
	Close()
}

// batchReader implements the BatchReader interface
type batchReader struct {
	reader    *bufio.Reader
	closer    io.Closer
	multiLine struct {
		buffer             string // Accumulated input
		inTripleQuotes     bool   // Whether in triple quote mode
		tripleQuoteType    string // Triple quote type (''' or """)
		inLineContinuation bool   // Whether in line continuation mode
	}
}

// NewBatchReader creates a new BatchReader instance
func NewBatchReader(input io.ReadCloser) BatchReader {
	return &batchReader{
		reader: bufio.NewReader(input),
		closer: input,
	}
}

// ReadInput implements the BatchReader interface
func (r *batchReader) ReadInput() (string, bool, bool, error) {
	for {
		input, err := r.readLine()
		if err != nil {
			if err == io.EOF {
				// Return any remaining buffered input at EOF
				if result := r.multiLine.buffer; result != "" {
					r.multiLine.buffer = ""
					return result, true, false, nil
				}
				return "", true, false, nil
			}
			return "", false, false, err
		}

		result, needMore := r.processLine(input)
		if !needMore {
			return result, false, false, nil
		}
	}
}

// Close implements the BatchReader interface
func (r *batchReader) Close() {
	if r.closer != nil {
		r.closer.Close()
	}
}

// readLine reads the raw line content
func (r *batchReader) readLine() (string, error) {
	var (
		isPartial bool  = true
		err       error = nil
		line, ln  []byte
	)
	for isPartial && err == nil {
		line, isPartial, err = r.reader.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

// processLine processes the input line, handling multi-line input logic
func (r *batchReader) processLine(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)

	if !r.multiLine.inTripleQuotes {
		// Check for triple quote start
		if trimmed == "'''" || trimmed == `"""` {
			r.multiLine.inTripleQuotes = true
			r.multiLine.tripleQuoteType = trimmed
			r.multiLine.buffer = ""
			return "", true
		}
		// Check for backslash line continuation
		if strings.HasSuffix(input, "\\") {
			r.multiLine.inLineContinuation = true
			r.multiLine.buffer += strings.TrimRight(input[:len(input)-1], " \t") + " "
			return "", true
		}
		// Regular line or end of continuation
		if r.multiLine.inLineContinuation {
			r.multiLine.buffer += input
			r.multiLine.inLineContinuation = false
		} else {
			r.multiLine.buffer = input
		}
		result := r.multiLine.buffer
		r.multiLine.buffer = ""
		return result, false
	}

	// Triple quote mode
	if trimmed == r.multiLine.tripleQuoteType {
		// End of triple quote
		r.multiLine.inTripleQuotes = false
		r.multiLine.tripleQuoteType = ""
		result := r.multiLine.buffer
		r.multiLine.buffer = ""
		return result, false
	}

	// Continue collecting lines in triple quote mode
	if strings.HasSuffix(input, "\\") {
		r.multiLine.buffer += strings.TrimRight(input[:len(input)-1], " \t") + " "
	} else {
		r.multiLine.buffer += input + "\n"
	}
	return "", true
}
