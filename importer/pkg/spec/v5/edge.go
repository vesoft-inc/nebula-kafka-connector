package specv5

import (
	"fmt"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
)

const (
	fmtEdgeInsertStatement = "USE %s INSERT EDGE %s %s" // "USE graph INSERT EDGE name (props),(props)"
	fmtEdgeValue           = "({%s})-[{%s}]->({%s})"    // "(srcID)-[props]->(dstID)"
)

type (
	Edge struct {
		Name  string       `yaml:"name"`
		Src   *EdgeNodeRef `yaml:"src"`
		Dst   *EdgeNodeRef `yaml:"dst"`
		Props Props        `yaml:"props,omitempty"`
	}

	EdgeNodeRef struct {
		Name string  `yaml:"name"`
		ID   *NodeID `yaml:"id"`
	}

	Edges []*Edge

	EdgeOption func(*Edge)
)

func NewEdge(name string, opts ...EdgeOption) *Edge {
	e := &Edge{
		Name: name,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func WithEdgeSrc(src *EdgeNodeRef) EdgeOption {
	return func(e *Edge) {
		e.Src = src
	}
}

func WithEdgeDst(dst *EdgeNodeRef) EdgeOption {
	return func(e *Edge) {
		e.Dst = dst
	}
}

func WithEdgeProps(props ...*Prop) EdgeOption {
	return func(e *Edge) {
		e.Props = append(e.Props, props...)
	}
}

func (e *Edge) Complete() {
	if e.Src != nil {
		e.Src.Complete()
	}
	if e.Dst != nil {
		e.Dst.Complete()
	}
	e.Props.Complete()
}

func (e *Edge) Validate() error {
	if e.Name == "" {
		return e.importError(errors.ErrNoEdgeName)
	}

	if e.Src == nil {
		return e.importError(errors.ErrNoEdgeSrc)
	}

	if err := e.Src.Validate(); err != nil {
		return e.importError(err)
	}

	if e.Dst == nil {
		return e.importError(errors.ErrNoEdgeDst)
	}

	if err := e.Dst.Validate(); err != nil {
		return e.importError(err)
	}

	if err := e.Props.Validate(); err != nil {
		return e.importError(err)
	}

	return nil
}

func (e *Edge) InsertStatement(graphName string, records ...Record) (string, error) {
	values := make([]string, 0, len(records))
	for _, record := range records {
		srcIDValue, err := e.Src.IDValue(record)
		if err != nil {
			return "", e.importError(err).SetGraphName(graphName)
		}
		dstIDValue, err := e.Dst.IDValue(record)
		if err != nil {
			return "", e.importError(err).SetGraphName(graphName)
		}
		propsValueList, err := e.Props.ValueList(record)
		if err != nil {
			return "", e.importError(err).SetGraphName(graphName)
		}
		values = append(values, fmt.Sprintf(
			fmtEdgeValue,
			srcIDValue,
			strings.Join(propsValueList, ", "),
			dstIDValue,
		))
	}
	return fmt.Sprintf(fmtEdgeInsertStatement, graphName, e.Name, strings.Join(values, ", ")), nil
}

func (e *Edge) importError(err error, formatWithArgs ...any) *errors.ImportError { //nolint:unparam
	return errors.AsOrNewImportError(err, formatWithArgs...).SetEdgeName(e.Name)
}

func (n *EdgeNodeRef) Complete() {
	if n.ID != nil {
		n.ID.Complete()
	}
}

func (n *EdgeNodeRef) Validate() error {
	if n.Name == "" {
		return n.importError(errors.ErrNoNodeName)
	}
	if n.ID == nil {
		return n.importError(errors.ErrNoNodeID)
	}
	//revive:disable-next-line:if-return
	if err := n.ID.Validate(); err != nil {
		return err
	}
	return nil
}

func (n *EdgeNodeRef) IDValue(record Record) (string, error) {
	return n.ID.Value(record)
}

func (n *EdgeNodeRef) importError(err error, formatWithArgs ...any) *errors.ImportError {
	return errors.AsOrNewImportError(err, formatWithArgs...).SetNodeName(n.Name)
}

func (es Edges) Complete() {
	for i := range es {
		es[i].Complete()
	}
}

func (es Edges) Validate() error {
	for i := range es {
		if err := es[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}
