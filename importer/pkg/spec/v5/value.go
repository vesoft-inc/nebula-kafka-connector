package specv5

import (
	"strings"
)

const (
	dbNULL = "NULL"

	ValueTypeInt           ValueType = "INT"
	ValueTypeString        ValueType = "STRING"
	ValueTypeDouble        ValueType = "DOUBLE"
	ValueTypeDate          ValueType = "DATE"
	ValueTypeTime          ValueType = "TIME"
	ValueTypeLocalTime     ValueType = "LOCAL TIME"
	ValueTypeDateTime      ValueType = "DATETIME"
	ValueTypeLocalDateTime ValueType = "LOCAL DATETIME"

	ValueTypeDefault = ValueTypeString
)

var (
	supportedPropValueTypes = map[ValueType]struct{}{
		ValueTypeInt:           {},
		ValueTypeString:        {},
		ValueTypeDouble:        {},
		ValueTypeDate:          {},
		ValueTypeTime:          {},
		ValueTypeLocalTime:     {},
		ValueTypeDateTime:      {},
		ValueTypeLocalDateTime: {},
	}

	supportedNodeIDValueTypes = map[ValueType]struct{}{
		ValueTypeInt:    {},
		ValueTypeString: {},
	}
)

type ValueType string

func IsSupportedPropValueType(t ValueType) bool {
	_, ok := supportedPropValueTypes[ValueType(strings.ToUpper(t.String()))]
	return ok
}

func IsSupportedNodeIDValueType(t ValueType) bool {
	_, ok := supportedNodeIDValueTypes[ValueType(strings.ToUpper(t.String()))]
	return ok
}

func (t ValueType) Equal(vt ValueType) bool {
	if !IsSupportedPropValueType(t) || !IsSupportedPropValueType(vt) {
		return false
	}
	return strings.EqualFold(t.String(), vt.String())
}

func (t ValueType) String() string {
	return string(t)
}
