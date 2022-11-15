package spec

import "strings"

const (
	ValueTypeInt    ValueType = "INT"
	ValueTypeString ValueType = "STRING"
	ValueTypeDouble ValueType = "DOUBLE"

	ValueTypeDefault = ValueTypeString
)

var supportedValueTypes = map[ValueType]struct{}{
	ValueTypeInt:    {},
	ValueTypeString: {},
	ValueTypeDouble: {},
}

type (
	ValueType string
)

func IsSupportedValueType(t ValueType) bool {
	_, ok := supportedValueTypes[ValueType(strings.ToUpper(t.String()))]
	return ok
}

func (t ValueType) Equal(vt ValueType) bool {
	if !IsSupportedValueType(t) || !IsSupportedValueType(vt) {
		return false
	}
	return strings.EqualFold(t.String(), vt.String())
}

func (t ValueType) String() string {
	return string(t)
}
