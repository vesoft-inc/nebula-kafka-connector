package utils

import (
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/clients"
)

type pathSelector struct {
	paths []string
	index int
}

func NewPathSelector(paths []string) *pathSelector {
	return &pathSelector{
		paths: paths,
		index: 0,
	}
}

func NewPathSelectorMap(storages []*clients.ServiceInfo) map[string]*pathSelector {
	sMap := make(map[string]*pathSelector)
	for _, storage := range storages {
		sMap[storage.Host] = NewPathSelector(storage.DataPaths)
	}
	return sMap
}

func (p *pathSelector) EvenlyGet() string {
	path := string(p.paths[p.index])
	p.index = (p.index + 1) % len(p.paths)
	return path
}
