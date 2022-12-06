package spec

import (
	"fmt"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
)

const (
	ValueTypeInt      ValueType = "INT"
	ValueTypeString   ValueType = "STRING"
	ValueTypeDouble   ValueType = "DOUBLE"
	ValueTypeDateTime ValueType = "DATETIME"

	ValueTypeDefault = ValueTypeString
)

var supportedValueTypes = map[ValueType]struct{}{
	ValueTypeInt:      {},
	ValueTypeString:   {},
	ValueTypeDouble:   {},
	ValueTypeDateTime: {},
}

type (
	ValueType string
)

func IsSupportedValueType(t ValueType) bool {
	_, ok := supportedValueTypes[ValueType(strings.ToUpper(t.String()))]
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
	if !IsSupportedValueType(t) || !IsSupportedValueType(vt) {
		return false
	}
	return strings.EqualFold(t.String(), vt.String())
}

func (t ValueType) String() string {
	return string(t)
}
