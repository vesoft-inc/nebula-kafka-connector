package specv3

import (
	"fmt"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/utils"
)

const (
	// INSERT VERTEX <tag_type> ( <prop_name_list> VALUES
	//		<vid> : ( <prop_value_list> )
	//		[, <vid> : ( <prop_value_list> ), ...]
	fmtNodeInsertStatement  = "INSERT VERTEX %s VALUES %s"
	fmtNodeNamePropNameList = "%s(%s)"
	fmtNodeValue            = "%s:(%s)"
)

type (
	// Node is VERTEX in 3.x
	Node struct {
		Name  string  `yaml:"name"`
		ID    *NodeID `yaml:"id"`
		Props Props   `yaml:"props,omitempty"`

		namePropNameList string // name(prop_name, ..., prop_name)
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
		n.ID.Name = strVID
	}
	n.Props.Complete()

	n.namePropNameList = fmt.Sprintf(
		fmtNodeNamePropNameList,
		utils.ConvertIdentifier(n.Name),
		strings.Join(n.Props.NameList(), ", "),
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

	return nil
}

func (n *Node) InsertStatement(graphName string, records ...Record) (string, error) {
	values := make([]string, 0, len(records))
	for _, record := range records {
		idValue, err := n.ID.Value(record)
		if err != nil {
			return "", n.importError(err).SetGraphName(graphName)
		}
		propValueList, err := n.Props.ValueList(record)
		if err != nil {
			return "", n.importError(err).SetGraphName(graphName)
		}
		values = append(values, fmt.Sprintf(fmtNodeValue, idValue, strings.Join(propValueList, ", ")))
	}
	return fmt.Sprintf(fmtNodeInsertStatement, n.namePropNameList, strings.Join(values, ", ")), nil
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
