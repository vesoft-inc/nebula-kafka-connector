package spec

import (
	"fmt"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
)

const (
	fmtPropValueStatement  = "%s: %s" // "name: value"
	fmtPropsValueStatement = "{%s}"   // "{prop, ..., prop}"
)

type (
	Prop struct {
		Name   string    `yaml:"name"`
		Type   ValueType `yaml:"type"`
		Index  int       `yaml:"index"`
		Ignore bool      `yaml:"ignore,omitempty"`
	}

	Props []*Prop
)

func (p *Prop) Complete() {
	if p.Ignore {
		return
	}
	if p.Type == "" {
		p.Type = ValueTypeDefault
	}
}

func (p *Prop) Validate() error {
	if p.Ignore {
		return nil
	}
	if p.Name == "" {
		return p.importError(errors.ErrNoPropName)
	}
	if !IsSupportedValueType(p.Type) {
		return p.importError(errors.ErrUnsupportedValueType, "unsupported type %s", p.Type)
	}
	return nil
}

func (p *Prop) ValueStatement(record Record) (string, error) {
	if p.Ignore {
		return "", nil
	}
	if p.Index < 0 || p.Index >= len(record) {
		return "", p.importError(errors.ErrNoRecord, "record index %d not exists", p.Index).SetRecord(record)
	}
	val := record[p.Index]
	if p.Type.Equal(ValueTypeString) {
		val = fmt.Sprintf("%q", val)
	}
	return fmt.Sprintf(fmtPropValueStatement, p.Name, val), nil
}

func (p *Prop) importError(err error, formatWithArgs ...any) *errors.ImportError {
	return errors.AsOrNewImportError(err, formatWithArgs...).SetPropName(p.Name)
}

func (ps Props) Complete() {
	for i := range ps {
		ps[i].Complete()
	}
}

func (ps Props) Validate() error {
	for i := range ps {
		if err := ps[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (ps Props) ValueStatement(record Record) (string, error) {
	statements := make([]string, 0, len(ps))
	for _, prop := range ps {
		if prop.Ignore {
			continue
		}
		statement, err := prop.ValueStatement(record)
		if err != nil {
			return "", err
		}
		statements = append(statements, statement)
	}
	return fmt.Sprintf(fmtPropsValueStatement, strings.Join(statements, ", ")), nil
}

func (ps Props) Append(props ...*Prop) Props {
	cpy := make(Props, len(ps)+len(props))
	copy(cpy, ps)
	copy(cpy[len(ps):], props)
	return cpy
}
