package importer

import (
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec"
)

type (
	Importer interface {
		Graph() *spec.Graph
		Node() *spec.Node
		Edge() *spec.Edge
		Import(records ...spec.Record) (*ImportResp, error)
	}

	ImportResp struct {
		Latency time.Duration
		ReqTime time.Duration
	}

	ImportResult struct {
		Resp *ImportResp
		Err  error
	}

	defaultImporter struct {
		graph  *spec.Graph
		node   *spec.Node
		edge   *spec.Edge
		client client.Client
	}
)

func NewNodeImporter(graph *spec.Graph, node *spec.Node, c client.Client) Importer {
	return &defaultImporter{
		graph:  graph,
		node:   node,
		client: c,
	}
}

func NewEdgeImporter(graph *spec.Graph, edge *spec.Edge, c client.Client) Importer {
	return &defaultImporter{
		graph:  graph,
		edge:   edge,
		client: c,
	}
}

func (i *defaultImporter) Graph() *spec.Graph {
	return i.graph
}

func (i *defaultImporter) Node() *spec.Node {
	return i.node
}

func (i *defaultImporter) Edge() *spec.Edge {
	return i.edge
}

func (i *defaultImporter) Import(records ...spec.Record) (*ImportResp, error) {
	statement, err := i.importStatement(records...)
	if err != nil {
		return nil, err
	}
	start := time.Now()
	rs, err := i.client.Execute(statement)
	if err != nil {
		return nil, errors.NewImportError(err).SetGraphName(i.graph.Name).SetStatement(statement)
	}
	if !rs.IsSucceed() {
		return nil, errors.NewImportError(err, "the status is %s ", rs.GetStatus()).
			SetGraphName(i.graph.Name).
			SetStatement(statement)
	}

	return &ImportResp{
		ReqTime: time.Since(start),
		Latency: time.Duration(rs.GetLatency()) * time.Microsecond,
	}, nil
}

func (i *defaultImporter) importStatement(records ...spec.Record) (string, error) {
	var statement string
	var err error
	if i.node != nil {
		statement, err = i.graph.NodeStatement(i.node, records...)
		if err != nil {
			return "", err
		}
	} else {
		statement, err = i.graph.EdgeStatement(i.edge, records...)
		if err != nil {
			return "", err
		}
	}
	return statement, nil
}
