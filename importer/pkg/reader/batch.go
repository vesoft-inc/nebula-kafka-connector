package reader

import (
	"io"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec"
)

type (
	BatchRecordReader interface {
		source.Sizer
		ReadBatch() (int, spec.Records, error)
	}

	defaultBatchReader struct {
		rr    RecordReader
		batch int
		eof   bool
	}
)

func NewBatchRecordReader(rr RecordReader, batch int) BatchRecordReader {
	if batch <= 0 {
		batch = DefaultBatchSize
	}
	return &defaultBatchReader{
		rr:    rr,
		batch: batch,
	}
}

func (r *defaultBatchReader) Size() (int64, error) {
	return r.rr.Size()
}

func (r *defaultBatchReader) ReadBatch() (int, spec.Records, error) {
	var (
		totalBytes int
		records    = make(spec.Records, 0, r.batch)
	)

	if r.eof {
		return 0, nil, io.EOF
	}

	for i := 0; i < r.batch; i++ {
		n, record, err := r.rr.Read()
		if err != nil {
			if err != io.EOF || totalBytes == 0 {
				return 0, nil, err
			}
			r.eof = true
			break
		}
		totalBytes += n
		records = append(records, record)
	}
	return totalBytes, records, nil
}
