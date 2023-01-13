package config

import (
	"fmt"
	"path/filepath"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/importer"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec"
)

type (
	Source struct {
		SourceConfig source.Config `yaml:",inline"`
		GraphName    string        `yaml:"graphName"`
		Nodes        spec.Nodes    `yaml:"nodes,omitempty"`
		Edges        spec.Edges    `yaml:"edges,omitempty"`
	}

	Sources []Source
)

func (s *Source) BuildGraph(opts ...spec.GraphOption) (*spec.Graph, error) {
	options := make([]spec.GraphOption, 0, len(s.Nodes)+len(s.Edges)+len(opts))
	for i := range s.Nodes {
		node := s.Nodes[i]
		options = append(options, spec.WithGraphNodes(node))
	}
	for i := range s.Edges {
		edge := s.Edges[i]
		options = append(options, spec.WithGraphEdges(edge))
	}
	options = append(options, opts...)
	graph := spec.NewGraph(s.GraphName, options...)
	graph.Complete()
	if err := graph.Validate(); err != nil {
		return nil, err
	}
	return graph, nil
}

func (s *Source) BuildImporters(c client.Client) ([]importer.Importer, error) {
	graph, err := s.BuildGraph()
	if err != nil {
		return nil, err
	}
	importers := make([]importer.Importer, 0, len(s.Nodes)+len(s.Edges))
	for k := range s.Nodes {
		node := s.Nodes[k]
		i := importer.NewNodeImporter(graph, node, c)
		importers = append(importers, i)
	}

	for k := range s.Edges {
		edge := s.Edges[k]
		i := importer.NewEdgeImporter(graph, edge, c)
		importers = append(importers, i)
	}
	return importers, nil
}

// OptimizeConfigPath Change source path to an absolute path in order to set it to the relative path of the config file
func (ss Sources) OptimizeConfigPath(configPath string) error {
	configAbsPath, err := filepath.Abs(filepath.Dir(configPath))
	if err != nil {
		return err
	}

	for i := range ss {
		if filepath.IsAbs(ss[i].SourceConfig.Path) {
			continue
		}
		ss[i].SourceConfig.Path = filepath.Join(configAbsPath, ss[i].SourceConfig.Path)
	}

	return nil
}

// OptimizeWildCard Convert sources with wildcards to a new sources
func (ss Sources) OptimizeWildCard() (Sources, error) {
	nss := make(Sources, 0, len(ss))
	for i := range ss {
		paths, err := filepath.Glob(ss[i].SourceConfig.Path)
		if err != nil {
			return nil, err
		}
		if len(paths) == 0 {
			return nil, fmt.Errorf("no file found for %s", ss[i].SourceConfig.Path)
		}

		for _, path := range paths {
			ns := ss[i]
			ns.SourceConfig.Path = path
			nss = append(nss, ns)
		}
	}

	return nss, nil
}
