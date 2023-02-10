package specv5

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/bytebufferpool"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
)

type (
	Node struct {
		Name  string  `yaml:"name"`
		ID    *NodeID `yaml:"id"`
		Props Props   `yaml:"props,omitempty"`

		graphName    string
		insertPrefix string // "USE %s INSERT NODE %s "
		propsWithID  Props
	}

	Nodes []*Node

	NodeOption func(*Node)
)

func NewNode(name string, opts ...NodeOption) *Node {
	n := &Node{
		Name: name,
	}
	n.Options(opts...)

	return n
}

func WithNodeGraphName(name string) NodeOption {
	return func(n *Node) {
		n.graphName = name
	}
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

func (n *Node) Options(opts ...NodeOption) *Node {
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *Node) Complete() {
	if n.ID != nil {
		n.ID.Complete()
	}
	n.Props.Complete()

	n.insertPrefix = fmt.Sprintf(
		"USE %s INSERT NODE %s ",
		n.graphName,
		n.Name,
	)
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

func (n *Node) InsertStatement(records ...Record) (string, error) {
	buff := bytebufferpool.Get()
	defer bytebufferpool.Put(buff)

	buff.SetString(n.insertPrefix)

	for i, record := range records {
		propsValueList, err := n.propsWithID.ValueList(record)
		if err != nil {
			return "", n.importError(err)
		}

		if i > 0 {
			_, _ = buff.WriteString(", ")
		}

		// "({%s})"
		_, _ = buff.WriteString("({")
		_, _ = buff.WriteStringSlice(propsValueList, ", ")
		_, _ = buff.WriteString("})")
	}
	return buff.String(), nil
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
