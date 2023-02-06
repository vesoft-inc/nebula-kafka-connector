package specv5

import (
	"fmt"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
)

const (
	fmtNodeInsertStatement = "USE %s INSERT NODE %s %s" // "USE graph INSERT NODE name (props),(props)"
	fmtNodeValue           = "({%s})"                   // "(props)"
)

type (
	Node struct {
		Name  string  `yaml:"name"`
		ID    *NodeID `yaml:"id"`
		Props Props   `yaml:"props,omitempty"`

		propsWithID Props
	}

	Nodes []*Node

	NodeOption func(*Node)
)

func NewNode(name string, opts ...NodeOption) *Node {
	n := &Node{
		Name: name,
	}

	for _, opt := range opts {
		opt(n)
	}

	return n
}

func WithNodeID(id *NodeID) NodeOption {
	return func(n *Node) {
		n.ID = id
	}
}

func WithNodeProps(props ...*Prop) NodeOption {
	return func(n *Node) {
		n.Props = append(n.Props, props...)
	}
}

func (n *Node) Complete() {
	if n.ID != nil {
		n.ID.Complete()
	}
	n.Props.Complete()
}

func (n *Node) Validate() error {
	if n.Name == "" {
		return n.importError(errors.ErrNoNodeName)
	}

	if n.ID == nil {
		return n.importError(errors.ErrNoNodeID)
	}

	if err := n.ID.Validate(); err != nil {
		return n.importError(err)
	}

	if err := n.Props.Validate(); err != nil {
		return n.importError(err)
	}

	n.propsWithID = Props{&Prop{
		Name:   n.ID.Name,
		Type:   n.ID.Type,
		picker: n.ID.picker,
	}}.Append(n.Props...)

	return nil
}

func (n *Node) InsertStatement(graphName string, records ...Record) (string, error) {
	values := make([]string, 0, len(records))
	for _, record := range records {
		propsValueList, err := n.propsWithID.ValueList(record)
		if err != nil {
			return "", n.importError(err).SetGraphName(graphName)
		}
		values = append(values, fmt.Sprintf(fmtNodeValue, strings.Join(propsValueList, ", ")))
	}
	return fmt.Sprintf(fmtNodeInsertStatement, graphName, n.Name, strings.Join(values, ", ")), nil
}

func (n *Node) importError(err error, formatWithArgs ...any) *errors.ImportError { //nolint:unparam
	return errors.AsOrNewImportError(err, formatWithArgs...).SetNodeName(n.Name)
}

func (ns Nodes) Complete() {
	for i := range ns {
		ns[i].Complete()
	}
}

func (ns Nodes) Validate() error {
	for i := range ns {
		if err := ns[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}
