package spec

import (
	"fmt"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/picker"
)

const (
	fmtNodeIDValueStatement = "{%s: %s}" // "{name: value}"
)

// TODO: 3.x support hash
var supportedNodeIDFunctions = map[string]struct{}{}

type (
	NodeID struct {
		Name        string        `yaml:"name"`
		Type        ValueType     `yaml:"type"`
		Index       int           `yaml:"index"`
		ConcatItems []interface{} `yaml:"concatItems,omitempty"` // only string and int is support, int is for Index
		Function    *string       `yaml:"function"`

		picker picker.Picker
	}
)

func IsSupportedNodeIDFunction(function string) bool {
	_, ok := supportedNodeIDFunctions[strings.ToUpper(function)]
	return ok
}

func (id *NodeID) Complete() {
	if id.Type == "" {
		id.Type = ValueTypeDefault
	}
}

func (id *NodeID) Validate() error {
	if id.Name == "" {
		return id.importError(errors.ErrNoNodeIDName)
	}
	if !IsSupportedNodeIDValueType(id.Type) {
		return id.importError(errors.ErrUnsupportedValueType, "unsupported type %s", id.Type)
	}
	if id.Function != nil && !IsSupportedNodeIDFunction(*id.Function) {
		return id.importError(errors.ErrUnsupportedFunction, "unsupported function %s", *id.Function)
	}
	if err := id.initPicker(); err != nil {
		return id.importError(err, "init picker failed")
	}

	return nil
}

func (id *NodeID) ValueStatement(record Record) (string, error) {
	val, err := id.picker.Pick(record)
	if err != nil {
		if len(id.ConcatItems) > 0 {
			return "", id.importError(err, "record concat items %v pick failed", id.ConcatItems).SetRecord(record)
		}
		return "", id.importError(err, "record index %d pick failed", id.Index).SetRecord(record)
	}
	return fmt.Sprintf(fmtNodeIDValueStatement, id.Name, val.Val), nil
}

func (id *NodeID) initPicker() error {
	pickerConfig := picker.Config{
		Type:     string(id.Type),
		Function: id.Function,
	}

	if len(id.ConcatItems) > 0 {
		for i, item := range id.ConcatItems {
			switch val := item.(type) {
			case int:
				pickerConfig.ConcatItems.AddIndex(val)
			case string:
				pickerConfig.ConcatItems.AddConstant(val)
			default:
				return id.importError(
					errors.ErrUnsupportedConcatItemType,
					"ConcatItems only support int or string, but the %d is %T", i, val,
				)
			}
		}
	} else {
		pickerConfig.Indices = []int{id.Index}
	}

	var err error
	id.picker, err = pickerConfig.Build()
	return err
}

func (id *NodeID) importError(err error, formatWithArgs ...any) *errors.ImportError {
	return errors.AsOrNewImportError(err, formatWithArgs...).SetPropName(id.Name)
}
