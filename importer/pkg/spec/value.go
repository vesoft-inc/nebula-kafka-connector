package spec

import (
	"fmt"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
)

const (
	dbNULL = "NULL"

	ValueTypeInt      ValueType = "INT"
	ValueTypeString   ValueType = "STRING"
	ValueTypeDouble   ValueType = "DOUBLE"
	ValueTypeDateTime ValueType = "DATETIME"

	ValueTypeDefault = ValueTypeString
)

var (
	supportedPropValueTypes = map[ValueType]struct{}{
		ValueTypeInt:      {},
		ValueTypeString:   {},
		ValueTypeDouble:   {},
		ValueTypeDateTime: {},
	}

	supportedNodeIDValueTypes = map[ValueType]struct{}{
		ValueTypeInt:    {},
		ValueTypeString: {},
	}
)

type (
	ValueType string
)

func IsSupportedPropValueType(t ValueType) bool {
	_, ok := supportedPropValueTypes[ValueType(strings.ToUpper(t.String()))]
	return ok
}

func IsSupportedNodeIDValueType(t ValueType) bool {
	_, ok := supportedNodeIDValueTypes[ValueType(strings.ToUpper(t.String()))]
	return ok
}

func (t ValueType) ValueStatement(val string) (string, error) {
	switch t {
	case ValueTypeInt:
		return val, nil
	case ValueTypeString:
		return fmt.Sprintf("%q", val), nil
	case ValueTypeDouble:
		return val, nil
	case ValueTypeDateTime:
		return fmt.Sprintf("DATETIME(%q)", val), nil
	}
	return "", errors.ErrUnsupportedValueType
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
