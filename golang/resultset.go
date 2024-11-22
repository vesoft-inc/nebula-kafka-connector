package nebula_ng

import (
	"fmt"
	"io"

	"github.com/vesoft-inc/nebula-ng-tools/golang/internal/decode"
	"github.com/vesoft-inc/nebula-ng-tools/golang/internal/generated_code/v5.0.0/proto/graph"
	"github.com/vesoft-inc/nebula-ng-tools/golang/internal/internal_error"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/types"
)

type resultSet struct {
	index   int
	table   *decode.ResultTable
	summary *graph.Summary
	cursor  []byte
}

type rowData struct {
	resultSet *resultSet
	values    []types.Value
}

func (rs *resultSet) HasNext() bool {
	if rs.table == nil {
		return false
	}
	return rs.index < int(rs.table.NumRecords())
}

func (rs *resultSet) Next() (types.Row, error) {
	if !rs.HasNext() {
		return nil, io.EOF
	}
	row := &rowData{
		resultSet: rs,
		values:    make([]types.Value, 0),
	}
	vs, err := rs.table.Next()
	if err != nil {
		return nil, err
	}
	rs.index++
	for _, v := range vs {
		row.values = append(row.values, v)
	}

	return row, nil
}

func (rs *resultSet) RowSize() int {
	if rs.table == nil {
		return 0
	}
	return int(rs.table.NumRecords())
}

func (rs *resultSet) ColumnTypes() []types.ColumnType {
	if rs.table == nil {
		return nil
	}
	return rs.table.ColumnTypes()
}

func (rs *resultSet) Summary() types.Summary {
	if rs.summary == nil {
		return nil
	}
	return &summary{summary: rs.summary}
}

func (rs *resultSet) Cursor() []byte {
	return rs.cursor
}

func (rs *resultSet) Columns() []string {
	if rs.table == nil {
		return nil
	}
	return rs.table.ColumnNames()
}

func (rs *resultSet) Scan(dsts ...any) error {
	if !rs.HasNext() {
		return io.EOF
	}

	row, err := rs.Next()
	if err != nil {
		return err
	}

	values := row.Values()
	if len(dsts) != len(values) {
		return internal_error.ErrInternal(fmt.Sprintf("scanner length not match values length"))
	}
	for i, dst := range dsts {
		if err := rs.convertValue(values[i], dst); err != nil {
			return err
		}
	}
	return nil
}

func (rs *resultSet) convertValue(src types.Value, dst any) error {
	switch dst.(type) {
	case *int, *uint, *float32, *float64, *string, *bool:
		return rs.convertBasicValue(src, dst)
	}
	if sn, ok := dst.(scanner); ok {
		return sn.scan(src)
	}

	return fmt.Errorf("cannot scan the value, unsupported type: %T", dst)
}

func (rs *resultSet) convertBasicValue(src types.Value, dst any) error {
	if src.IsNull() {
		return fmt.Errorf("value is null")
	}
	switch d := dst.(type) {
	case *int:
		switch src.GetType() {
		case types.ValueTypeInt8:
			v, _ := src.AsInt8()
			*d = int(v)
		case types.ValueTypeInt16:
			v, _ := src.AsInt16()
			*d = int(v)
		case types.ValueTypeInt32:
			v, _ := src.AsInt32()
			*d = int(v)
		case types.ValueTypeInt64:
			v, _ := src.AsInt64()
			*d = int(v)
		case types.ValueTypeUInt8:
			v, _ := src.AsUInt8()
			*d = int(v)
		case types.ValueTypeUInt16:
			v, _ := src.AsUInt16()
			*d = int(v)
		case types.ValueTypeUInt32:
			v, _ := src.AsInt32()
			*d = int(v)
		case types.ValueTypeUInt64:
			v, _ := src.AsInt64()
			*d = int(v)
		default:
			return internal_error.ErrInternal(fmt.Sprintf("value type not match"))
		}
	case *uint:
		switch src.GetType() {
		case types.ValueTypeUInt8:
			v, _ := src.AsUInt8()
			*d = uint(v)
		case types.ValueTypeUInt16:
			v, _ := src.AsUInt16()
			*d = uint(v)
		case types.ValueTypeUInt32:
			v, _ := src.AsInt32()
			*d = uint(v)
		case types.ValueTypeUInt64:
			v, _ := src.AsInt64()
			*d = uint(v)
		default:
			return internal_error.ErrInternal(fmt.Sprintf("value type not match"))
		}
	case *float32:
		switch src.GetType() {
		case types.ValueTypeFloat:
			v, _ := src.AsFloat()
			*d = float32(v)
		default:
			return internal_error.ErrInternal(fmt.Sprintf("value type not match"))
		}
	case *float64:
		switch src.GetType() {
		case types.ValueTypeFloat:
			v, _ := src.AsFloat()
			*d = float64(v)
		case types.ValueTypeDouble:
			v, _ := src.AsDouble()
			*d = float64(v)
		default:
			return internal_error.ErrInternal(fmt.Sprintf("value type not match"))
		}
	case *string:
		switch src.GetType() {
		case types.ValueTypeString:
			v, _ := src.AsString()
			*d = string(v)
		default:
			return internal_error.ErrInternal(fmt.Sprintf("value type not match"))
		}
	case *bool:
		switch src.GetType() {
		case types.ValueTypeBool:
			v, _ := src.AsBool()
			*d = bool(v)
		default:
			return internal_error.ErrInternal(fmt.Sprintf("value type not match"))
		}
	}
	return nil
}

func (rd *rowData) Values() []types.Value {
	return rd.values
}

func (rd *rowData) GetValueByName(name string) (types.Value, error) {
	names := rd.resultSet.Columns()
	var index int = -1
	for i, n := range names {
		if string(n) == name {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, internal_error.ErrInternal(fmt.Sprintf("column %s not found", name))
	}
	return rd.values[index], nil
}

func (rd *rowData) GetValueByIndex(index int) (types.Value, error) {
	if index < 0 || index >= len(rd.values) {
		return nil, internal_error.ErrInternal(fmt.Sprintf("index out of range"))
	}
	return rd.values[index], nil
}
