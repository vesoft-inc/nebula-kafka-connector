package configv3

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/importer"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
	specv3 "github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec/v3"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/utils"
)

type (
	Source struct {
		SourceConfig source.Config `yaml:",inline"`
		Nodes        specv3.Nodes  `yaml:"tags,omitempty"`
		Edges        specv3.Edges  `yaml:"edges,omitempty"`
	}

	Sources []Source
)

func (s *Source) BuildGraph(graphName string, opts ...specv3.GraphOption) (*specv3.Graph, error) {
	options := make([]specv3.GraphOption, 0, len(s.Nodes)+len(s.Edges)+len(opts))
	for i := range s.Nodes {
		node := s.Nodes[i]
		options = append(options, specv3.WithGraphNodes(node))
	}
	for i := range s.Edges {
		edge := s.Edges[i]
		options = append(options, specv3.WithGraphEdges(edge))
	}
	options = append(options, opts...)
	graph := specv3.NewGraph(graphName, options...)
	graph.Complete()
	if err := graph.Validate(); err != nil {
		return nil, err
	}
	return graph, nil
}

func (s *Source) BuildImporters(graphName string, pool client.Pool) ([]importer.Importer, error) {
	graph, err := s.BuildGraph(graphName)
	if err != nil {
		return nil, err
	}
	importers := make([]importer.Importer, 0, len(s.Nodes)+len(s.Edges))
	for k := range s.Nodes {
		node := s.Nodes[k]
		builder := graph.InsertNodeBuilder(node)
		i := importer.New(builder, pool)
		importers = append(importers, i)
	}

	for k := range s.Edges {
		edge := s.Edges[k]
		builder := graph.InsertEdgeBuilder(edge)
		i := importer.New(builder, pool)
		importers = append(importers, i)
	}
	return importers, nil
}

// OptimizePath optimizes relative paths base to the configuration file path
func (ss Sources) OptimizePath(configPath string) error {
	configPathDir := filepath.Dir(configPath)
	for i := range ss {
		ss[i].SourceConfig.Path = utils.RelativePathBaseOn(configPathDir, ss[i].SourceConfig.Path)
	}
	return nil
}

// OptimizePathWildCard optimizes the wildcards in the paths
func (ss *Sources) OptimizePathWildCard() error {
	nss := make(Sources, 0, len(*ss))
	for i := range *ss {
		paths, err := filepath.Glob((*ss)[i].SourceConfig.Path)
		if err != nil {
			return err
		}
		if len(paths) == 0 {
			return &os.PathError{Op: "open", Path: (*ss)[i].SourceConfig.Path, Err: fs.ErrNotExist}
		}

		for _, path := range paths {
			cpy := (*ss)[i]
			cpy.SourceConfig.Path = path
			nss = append(nss, cpy)
		}
	}
	*ss = nss
	return nil
}
