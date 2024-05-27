package nebula_ng

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Bool bool
type Int8 int8
type Int16 int16
type Int32 int32
type Int64 int64
type UInt8 uint8
type UInt16 uint16
type UInt32 uint32
type UInt64 uint64
type Float float32
type Double float64
type String string
type ValueType uint8

type (
	Value interface {
		String() string
		GetType() ValueType
		IsNull() bool
		AsBool() (Bool, error)
		AsInt8() (Int8, error)
		AsInt16() (Int16, error)
		AsInt32() (Int32, error)
		AsInt64() (Int64, error)
		AsUInt8() (UInt8, error)
		AsUInt16() (UInt16, error)
		AsUInt32() (UInt32, error)
		AsUInt64() (UInt64, error)
		AsFloat() (Float, error)
		AsDouble() (Double, error)
		AsString() (String, error)
		AsList() (List, error)
		AsRecord() (Record, error)
		AsDuration() (Duration, error)
		AsLocalTime() (LocalTime, error)
		AsLocalDatetime() (LocalDatetime, error)
		AsDate() (Date, error)
		AsZonedDatetime() (ZonedDatetime, error)
		AsZonedTime() (ZonedTime, error)
		AsNode() (Node, error)
		AsEdge() (Edge, error)
		AsPath() (Path, error)
	}
	List interface {
		String() string
		GetValues() []Value
		Size() int
	}
	Record interface {
		String() string
		GetValues() map[string]Value
	}

	Duration interface {
		IsMonthBased() bool
		String() string
		GetYear() int32
		GetMonth() int32
		GetDay() int32
		GetMinute() int32
		GetSecond() int32
		GetMicrosecond() int32
	}

	LocalDatetime interface {
		String() string
		Date
		LocalTime
	}

	LocalTime interface {
		String() string
		GetHour() uint32
		GetMinute() uint32
		GetSec() uint32
		GetMicrosec() uint32
	}
	ZonedDatetime interface {
		LocalDatetime
		TimeZone
		Time() *time.Time
	}
	ZonedTime interface {
		LocalTime
		TimeZone
	}

	TimeZone interface {
		//GetOffset return the time zone offset in seconds
		GetOffset() int
	}

	Date interface {
		String() string
		GetYear() int32
		GetMonth() uint32
		GetDay() uint32
	}

	Node interface {
		String() string
		GetProperties() map[string]Value
		GetGraph() string
		GetType() string
		GetLabels() []string
		GetId() int64
	}
	Edge interface {
		String() string
		GetProperties() map[string]Value
		GetSrcId() int64
		GetDstId() int64
		GetGraph() string
		GetType() string
		GetLabels() []string
		GetRank() int64
		IsDirected() bool
	}
	Path interface {
		String() string
		GetValues() []Value
	}

	mapValue map[string]Value
)

const (
	ValueTypeNull ValueType = iota
	ValueTypeBool
	ValueTypeInt8
	ValueTypeInt16
	ValueTypeInt32
	ValueTypeInt64
	ValueTypeUInt8
	ValueTypeUInt16
	ValueTypeUInt32
	ValueTypeUInt64
	ValueTypeFloat
	ValueTypeDouble
	ValueTypeString
	ValueTypeList
	ValueTypeRecord
	ValueTypeNode
	ValueTypeEdge
	ValueTypePath
	ValueTypeDuration
	ValueTypeDate
	ValueTypeLocalTime
	ValueTypeLocalDateTime
	ValueTypeZonedTime
	ValueTypeZonedDateTime
)

func (vt ValueType) String() string {
	switch vt {
	case ValueTypeNull:
		return "NULL"
	case ValueTypeBool:
		return "BOOL"
	case ValueTypeInt8:
		return "INT8"
	case ValueTypeInt16:
		return "INT16"
	case ValueTypeInt32:
		return "INT32"
	case ValueTypeInt64:
		return "INT64"
	case ValueTypeUInt8:
		return "UINT8"
	case ValueTypeUInt16:
		return "UINT16"
	case ValueTypeUInt32:
		return "UINT32"
	case ValueTypeUInt64:
		return "UINT64"
	case ValueTypeFloat:
		return "FLOAT"
	case ValueTypeDouble:
		return "DOUBLE"
	case ValueTypeString:
		return "STRING"
	case ValueTypeDuration:
		return "DURATION"
	case ValueTypeDate:
		return "DATE"
	case ValueTypeLocalTime:
		return "LOCALTIME"
	case ValueTypeLocalDateTime:
		return "LOCALDATETIME"
	case ValueTypeZonedTime:
		return "ZONEDTIME"
	case ValueTypeZonedDateTime:
		return "ZONEDDATETIME"
	case ValueTypeList:
		return "LIST"
	case ValueTypeRecord:
		return "RECORD"
	case ValueTypeNode:
		return "NODE"
	case ValueTypeEdge:
		return "EDGE"
	case ValueTypePath:
		return "PATH"
	}
	return "UNKNOWN"
}

// EmptyValue just for test
type EmptyValue struct{}

var _ Value = &EmptyValue{}

func (nv *EmptyValue) String() string {
	return ""
}

func (nv *EmptyValue) GetType() ValueType {
	return ValueTypeNull
}

func (nv *EmptyValue) IsNull() bool {
	return true
}

func (nv *EmptyValue) AsBool() (Bool, error) {
	return false, nil
}

func (nv *EmptyValue) AsInt8() (Int8, error) {
	return 0, nil
}

func (nv *EmptyValue) AsInt16() (Int16, error) {
	return 0, nil
}

func (nv *EmptyValue) AsInt32() (Int32, error) {
	return 0, nil
}

func (nv *EmptyValue) AsInt64() (Int64, error) {
	return 0, nil
}

func (nv *EmptyValue) AsUInt8() (UInt8, error) {
	return 0, nil
}

func (nv *EmptyValue) AsUInt16() (UInt16, error) {
	return 0, nil
}

func (nv *EmptyValue) AsUInt32() (UInt32, error) {
	return 0, nil
}

func (nv *EmptyValue) AsUInt64() (UInt64, error) {
	return 0, nil
}

func (nv *EmptyValue) AsFloat() (Float, error) {
	return 0, nil
}

func (nv *EmptyValue) AsDouble() (Double, error) {
	return 0, nil
}

func (nv *EmptyValue) AsString() (String, error) {
	return "", nil
}

func (nv *EmptyValue) AsList() (List, error) {
	return nil, nil
}

func (nv *EmptyValue) AsRecord() (Record, error) {
	return nil, nil
}

func (nv *EmptyValue) AsDuration() (Duration, error) {
	return nil, nil
}

// AsTime() (time.Time, error)
// AsDatetime() (Datetime, error)
func (nv *EmptyValue) AsNode() (Node, error) {
	return nil, nil
}

func (nv *EmptyValue) AsEdge() (Edge, error) {
	return nil, nil
}

func (nv *EmptyValue) AsPath() (Path, error) {
	return nil, nil
}
func (nv *EmptyValue) AsLocalDatetime() (LocalDatetime, error) {
	return nil, nil
}

func (nv *EmptyValue) AsLocalTime() (LocalTime, error) {
	return nil, nil
}

func (nv *EmptyValue) AsZonedTime() (ZonedTime, error) {
	return nil, nil
}
func (nv *EmptyValue) AsZonedDatetime() (ZonedDatetime, error) {
	return nil, nil
}

func (nv *EmptyValue) AsDate() (Date, error) {
	return nil, nil
}

func (s *String) String() string {
	inputbytes := []byte(*s)
	// Create a byte slice to store the output, but don't specify the length!!!
	outputBytes := []byte{}
	// Initialize the index to track the current position
	idx := 0
	for _, ch := range inputbytes {
		if ch == '\r' {
			// Carriage return: move to the beginning of the line by setting idx to 0
			idx = 0
		} else {
			// If the index is out of range, append the character to the output
			if idx >= len(outputBytes) {
				outputBytes = append(outputBytes, ch)
			} else {
				// Otherwise, overwrite the character at the index
				outputBytes[idx] = ch
			}
			idx++
		}
	}
	return fmt.Sprintf(`%s`, string(outputBytes))
}

func (m mapValue) string() string {
	var kvStr []string = make([]string, 0, len(m))
	var keys []string = make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	for _, k := range keys {
		v := m[k]
		kvTemp := fmt.Sprintf(`%s:%s`, k, v)
		kvStr = append(kvStr, kvTemp)
	}
	return strings.Join(kvStr, ",")
}
