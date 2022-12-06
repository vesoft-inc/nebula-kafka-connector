package config

import (
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/manager"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec"
)

type (
	Graph struct {
		Name                string        `yaml:"name"`
		Nodes               Nodes         `yaml:"nodes"`
		Edges               Edges         `yaml:"edges"`
		Batch               int           `yaml:"batch,omitempty"`
		ReaderConcurrency   int           `yaml:"readerConcurrency,omitempty"`
		ImporterConcurrency int           `yaml:"importerConcurrency,omitempty"`
		StatsInterval       time.Duration `yaml:"statsInterval,omitempty"`
		Hooks               manager.Hooks `yaml:"hooks,omitempty"`
	}

	Node struct {
		spec.Node     `yaml:",inline"`
		SourceConfigs []*source.Config `yaml:"sources"`
	}
	Nodes []Node

	Edge struct {
		spec.Edge     `yaml:",inline"`
		SourceConfigs []*source.Config `yaml:"sources"`
	}

	Edges []Edge
)

func (g *Graph) BuildGraph(opts ...spec.GraphOption) (*spec.Graph, error) {
	options := make([]spec.GraphOption, 0, len(g.Nodes)+len(g.Edges)+len(opts))
	for i := range g.Nodes {
		node := g.Nodes[i]
		options = append(options, spec.WithGraphNodes(&node.Node))
	}
	for i := range g.Edges {
		edge := g.Edges[i]
		options = append(options, spec.WithGraphEdges(&edge.Edge))
	}
	options = append(options, opts...)
	graph := spec.NewGraph(g.Name, options...)
	graph.Complete()
	if err := graph.Validate(); err != nil {
		return nil, err
	}
	return graph, nil
}

func (g *Graph) BuildManager(opts ...manager.Option) (manager.Manager, error) {
	options := make([]manager.Option, 0, 5+len(opts))
	options = append(options,
		manager.WithBatch(g.Batch),
		manager.WithReaderConcurrency(g.ReaderConcurrency),
		manager.WithImporterConcurrency(g.ImporterConcurrency),
		manager.WithStatsInterval(g.StatsInterval),
		manager.WithBeforeHooks(g.Hooks.Before...),
		manager.WithAfterHooks(g.Hooks.After...),
	)
	options = append(options, opts...)

	mgr := manager.NewWithOpts(options...)

	for i := range g.Nodes {
		node := g.Nodes[i]
		if err := mgr.ImportNode(node.Name, node.SourceConfigs...); err != nil {
			return nil, err
		}
	}
	for i := range g.Edges {
		edge := g.Edges[i]
		if err := mgr.ImportEdge(edge.Name, edge.SourceConfigs...); err != nil {
			return nil, err
		}
	}

	return mgr, nil
}
