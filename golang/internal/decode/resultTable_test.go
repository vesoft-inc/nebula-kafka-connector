package decode

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummpBatch struct {
	name string
	rows uint32
}

var _ batcher = &dummpBatch{}

func (b *dummpBatch) numRecords() uint32 {
	return b.rows
}
func (b *dummpBatch) getRowByIndex(index uint32) ([]*nebulaValue, error) {
	if index >= b.rows {
		return nil, io.EOF
	}
	return nil, nil
}

func TestTable(t *testing.T) {
	testcases := []struct {
		name         string
		batches      int
		rowsPerBatch int
		nextTimes    int
		hasErr       bool
	}{
		{"1st", 1, 1, 1, false},
		{"2nd", 1, 1, 2, true},
		{"3rd", 1, 2, 2, false},
		{"4th", 1, 2, 3, true},
		{"5th", 2, 1, 1, false},
		{"6th", 2, 1, 2, false},
		{"7th", 2, 1, 3, true},
		{"8th", 2, 2, 1, false},
		{"9th", 2, 2, 4, false},
		{"10th", 2, 2, 5, true},
		{"11th", 100, 100, 9999, false},
		{"12th", 100, 100, 10000, false},
		{"13th", 100, 100, 10001, true},
	}
	for _, tc := range testcases {
		batches := make([]batcher, tc.batches)
		for i := 0; i < tc.batches; i++ {
			batches[i] = &dummpBatch{name: fmt.Sprintf("%d", i), rows: uint32(tc.rowsPerBatch)}
		}
		tbl := &ResultTable{
			batches:    batches,
			numBatches: tc.batches,
		}
		for i := 0; i < tc.nextTimes-1; i++ {
			_, err := tbl.Next()
			if err != nil {
				t.Fatal(err)
			}
		}
		// the last batch name should be tc.batches-1
		lastBatchName := fmt.Sprintf("%d", tc.batches-1)
		b, ok := tbl.batches[tc.batches-1].(*dummpBatch)
		if !ok {
			t.Fatal("type assertion failed")
		}
		assert.Equal(t, lastBatchName, b.name)
		_, err := tbl.Next()
		if tc.hasErr {
			assert.Error(t, err, tc.name)
		} else {
			assert.NoError(t, err, tc.name)
		}
	}
}
