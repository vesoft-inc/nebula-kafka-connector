package reader

import (
	"bufio"
	"encoding/csv"
	"io"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec"
)

type (
	csvReader struct {
		s  source.Source
		rr *remainingReader
		br *bufio.Reader
		cr *csv.Reader
	}

	remainingReader struct {
		io.Reader
		remaining int
	}
)

func NewCSVReader(s source.Source) RecordReader {
	rr := &remainingReader{Reader: s}
	br := bufio.NewReader(rr)
	cr := csv.NewReader(br)

	if c := s.Config(); c != nil && c.CSV != nil {
		if chars := []rune(c.CSV.Delimiter); len(chars) > 0 {
			cr.Comma = chars[0]
		}
	}

	return &csvReader{
		s:  s,
		rr: rr,
		br: br,
		cr: cr,
	}
}

func (r *csvReader) Size() (int64, error) {
	return r.s.Size()
}

func (r *csvReader) Read() (int, spec.Record, error) {
	record, err := r.cr.Read()
	return r.rr.Take(r.br.Buffered()), record, err
}

func (r *remainingReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.remaining += n
	return n, err
}

func (r *remainingReader) Take(buffered int) (n int) {
	n, r.remaining = r.remaining-buffered, buffered
	return n
}
