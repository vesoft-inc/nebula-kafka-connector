package nebula_ng

import (
	"fmt"
	"io"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/graph"
)

type resultSet struct {
	index   int
	result  *graph.ResultTable
	summary *graph.Summary
	cursor  []byte
}

type rowData struct {
	resultSet *resultSet
	values    []Value
}

func (rs *resultSet) HasNext() bool {
	if rs.result == nil {
		return false
	}
	return rs.index < len(rs.result.Records)
}

func (rs *resultSet) Next() (Row, error) {
	if !rs.HasNext() {
		return nil, io.EOF
	}

	values := rs.result.GetRecords()[rs.index].Values
	row := &rowData{
		resultSet: rs,
		values:    make([]Value, 0, len(values)),
	}

	for _, v := range values {
		row.values = append(row.values, &grpcValue{data: v})
	}

	rs.index++
	return row, nil
}

func (rs *resultSet) RowSize() int {
	if rs.result == nil {
		return 0
	}
	return len(rs.result.GetRecords())
}

func (rs *resultSet) ColumnTypes() []ColumnType {
	return nil
}

func (rs *resultSet) Summary() Summary {
	if rs.summary == nil {
		return nil
	}
	return &summary{summary: rs.summary}
}

func (rs *resultSet) Cursor() []byte {
	return rs.cursor
}

func (rs *resultSet) Columns() []string {
	if rs.result == nil {
		return nil
	}
	names := rs.result.ColumnNames
	var cols []string
	for _, name := range names {
		cols = append(cols, string(name))
	}
	return cols
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
		return errInternel(fmt.Sprintf("scanner length not match values length"))
	}
	for i, dst := range dsts {
		if err := rs.convertValue(values[i], dst); err != nil {
			return err
		}
	}
	return nil
}

func (rs *resultSet) convertValue(src Value, dst any) error {
	switch dst.(type) {
	case *int, *uint, *float32, *float64, *string, *bool:
		return rs.convertBasicValue(src, dst)
	}
	if sn, ok := dst.(scanner); ok {
		return sn.scan(src)
	}

	return fmt.Errorf("cannot scan the value, unsupported type: %T", dst)
}

func (rs *resultSet) convertBasicValue(src Value, dst any) error {
	if src.IsNull() {
		return fmt.Errorf("value is null")
	}
	switch d := dst.(type) {
	case *int:
		switch src.GetType() {
		case ValueTypeInt8:
			v, _ := src.AsInt8()
			*d = int(v)
		case ValueTypeInt16:
			v, _ := src.AsInt16()
			*d = int(v)
		case ValueTypeInt32:
			v, _ := src.AsInt32()
			*d = int(v)
		case ValueTypeInt64:
			v, _ := src.AsInt64()
			*d = int(v)
		case ValueTypeUInt8:
			v, _ := src.AsUInt8()
			*d = int(v)
		case ValueTypeUInt16:
			v, _ := src.AsUInt16()
			*d = int(v)
		case ValueTypeUInt32:
			v, _ := src.AsInt32()
			*d = int(v)
		case ValueTypeUInt64:
			v, _ := src.AsInt64()
			*d = int(v)
		default:
			return errInternel(fmt.Sprintf("value type not match"))
		}
	case *uint:
		switch src.GetType() {
		case ValueTypeUInt8:
			v, _ := src.AsUInt8()
			*d = uint(v)
		case ValueTypeUInt16:
			v, _ := src.AsUInt16()
			*d = uint(v)
		case ValueTypeUInt32:
			v, _ := src.AsInt32()
			*d = uint(v)
		case ValueTypeUInt64:
			v, _ := src.AsInt64()
			*d = uint(v)
		default:
			return errInternel(fmt.Sprintf("value type not match"))
		}
	case *float32:
		switch src.GetType() {
		case ValueTypeFloat:
			v, _ := src.AsFloat()
			*d = float32(v)
		default:
			return errInternel(fmt.Sprintf("value type not match"))
		}
	case *float64:
		switch src.GetType() {
		case ValueTypeFloat:
			v, _ := src.AsFloat()
			*d = float64(v)
		case ValueTypeDouble:
			v, _ := src.AsDouble()
			*d = float64(v)
		default:
			return errInternel(fmt.Sprintf("value type not match"))
		}
	case *string:
		switch src.GetType() {
		case ValueTypeString:
			v, _ := src.AsString()
			*d = string(v)
		default:
			return errInternel(fmt.Sprintf("value type not match"))
		}
	case *bool:
		switch src.GetType() {
		case ValueTypeBool:
			v, _ := src.AsBool()
			*d = bool(v)
		default:
			return errInternel(fmt.Sprintf("value type not match"))
		}
	}
	return nil
}

func (rd *rowData) Values() []Value {
	return rd.values
}

func (rd *rowData) GetValueByName(name string) (Value, error) {
	names := rd.resultSet.Columns()
	var index int = -1
	for i, n := range names {
		if string(n) == name {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, errInternel(fmt.Sprintf("column %s not found", name))
	}
	return rd.values[index], nil
}

func (rd *rowData) GetValueByIndex(index int) (Value, error) {
	if index < 0 || index >= len(rd.values) {
		return nil, errInternel(fmt.Sprintf("index out of range"))
	}
	return rd.values[index], nil
}
