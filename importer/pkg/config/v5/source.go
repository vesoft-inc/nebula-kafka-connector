package configv5

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	configbase "github.com/vesoft-inc/nebula-ng-tools/importer/pkg/config/base"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/importer"
	specv5 "github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec/v5"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/utils"
)

type (
	Source struct {
		configbase.Source `yaml:",inline"`
		Nodes             specv5.Nodes `yaml:"nodes,omitempty"`
		Edges             specv5.Edges `yaml:"edges,omitempty"`
	}

	Sources []Source
)

func (s *Source) BuildGraph(graphName string, opts ...specv5.GraphOption) (*specv5.Graph, error) {
	options := make([]specv5.GraphOption, 0, len(s.Nodes)+len(s.Edges)+len(opts))
	for i := range s.Nodes {
		node := s.Nodes[i].Options(specv5.WithNodeGraphName(graphName))
		options = append(options, specv5.WithGraphNodes(node))
	}
	for i := range s.Edges {
		edge := s.Edges[i].Options(specv5.WithEdgeGraphName(graphName))
		options = append(options, specv5.WithGraphEdges(edge))
	}
	options = append(options, opts...)
	graph := specv5.NewGraph(graphName, options...)
	graph.Complete()
	if err := graph.Validate(); err != nil {
		return nil, err
	}
	return graph, nil
}

func (s *Source) BuildImporters(wgMap *utils.WaitGroupMap, graphName string, pool client.Pool) ([]importer.Importer, error) {
	graph, err := s.BuildGraph(graphName)
	if err != nil {
		return nil, err
	}
	importers := make([]importer.Importer, 0, len(s.Nodes)+len(s.Edges))
	for k := range s.Nodes {
		node := s.Nodes[k]
		builder := graph.InsertNodeBuilder(node)
		wgMap.Add(1, node.Name)
		i := importer.New(builder, pool, importer.WithDoneFunc(func() {
			wgMap.Done(node.Name)
		}))
		importers = append(importers, i)
	}

	for k := range s.Edges {
		edge := s.Edges[k]
		builder := graph.InsertEdgeBuilder(edge)
		i := importer.New(builder, pool, importer.WithWaitFunc(func() {
			wgMap.WaitMany(edge.Src.Name, edge.Dst.Name)
		}))
		importers = append(importers, i)
	}
	return importers, nil
}

// OptimizePath optimizes relative paths base to the configuration file path
func (ss Sources) OptimizePath(configPath string) error {
	configPathDir := filepath.Dir(configPath)
	for i := range ss {
		if ss[i].SourceConfig.Local != nil {
			ss[i].SourceConfig.Local.Path = utils.RelativePathBaseOn(configPathDir, ss[i].SourceConfig.Local.Path)
		}
	}
	return nil
}

// OptimizePathWildCard optimizes the wildcards in the paths
func (ss *Sources) OptimizePathWildCard() error {
	nss := make(Sources, 0, len(*ss))
	for i := range *ss {
		if (*ss)[i].SourceConfig.Local != nil {
			paths, err := filepath.Glob((*ss)[i].SourceConfig.Local.Path)
			if err != nil {
				return err
			}
			if len(paths) == 0 {
				return &os.PathError{Op: "open", Path: (*ss)[i].SourceConfig.Local.Path, Err: fs.ErrNotExist}
			}

			for _, path := range paths {
				cpy := (*ss)[i]
				cpySourceConfig := cpy.SourceConfig.Clone()
				cpy.SourceConfig = *cpySourceConfig
				cpy.SourceConfig.Local.Path = path
				nss = append(nss, cpy)
			}
		}
	}
	*ss = nss
	return nil
}
