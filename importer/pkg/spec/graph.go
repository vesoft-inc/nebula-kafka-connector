package spec

import "github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"

type (
	Graph struct {
		Name  string `yaml:"name"`
		Nodes Nodes  `yaml:"nodes"`
		Edges Edges  `yaml:"edges"`
	}

	GraphOption func(*Graph)
)

func NewGraph(name string, opts ...GraphOption) *Graph {
	g := &Graph{
		Name: name,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

func WithGraphNodes(nodes ...*Node) GraphOption {
	return func(g *Graph) {
		g.AddNodes(nodes...)
	}
}

func WithGraphEdges(edges ...*Edge) GraphOption {
	return func(g *Graph) {
		g.AddEdges(edges...)
	}
}

func (g *Graph) AddNodes(nodes ...*Node) {
	g.Nodes = append(g.Nodes, nodes...)
}

func (g *Graph) AddEdges(edges ...*Edge) {
	g.Edges = append(g.Edges, edges...)
}

func (g *Graph) Complete() {
	if g.Nodes != nil {
		g.Nodes.Complete()
	}
	if g.Edges != nil {
		g.Edges.Complete()
	}
}

func (g *Graph) Validate() error {
	if g.Name == "" {
		return errors.ErrNoGraphName
	}
	if err := g.Nodes.Validate(); err != nil {
		return err
	}
	//revive:disable-next-line:if-return
	if err := g.Edges.Validate(); err != nil {
		return err
	}

	return nil
}

func (g *Graph) NodeStatement(n *Node, records ...Record) (string, error) {
	statement, err := n.ValueStatement(g.Name, records...)
	if err != nil {
		return "", g.importError(err).SetNodeName(n.Name)
	}
	return statement, nil
}

func (g *Graph) EdgeStatement(e *Edge, records ...Record) (string, error) {
	statement, err := e.ValueStatement(g.Name, records...)
	if err != nil {
		return "", g.importError(err).SetEdgeName(e.Name)
	}
	return statement, nil
}

func (g *Graph) GetNodeByName(name string) (*Node, bool) {
	for _, n := range g.Nodes {
		if n.Name == name {
			return n, true
		}
	}
	return nil, false
}

func (g *Graph) GetEdgeByName(name string) (*Edge, bool) {
	for _, e := range g.Edges {
		if e.Name == name {
			return e, true
		}
	}
	return nil, false
}

func (g *Graph) importError(err error, formatWithArgs ...any) *errors.ImportError {
	return errors.AsOrNewImportError(err, formatWithArgs...).SetGraphName(g.Name)
}
