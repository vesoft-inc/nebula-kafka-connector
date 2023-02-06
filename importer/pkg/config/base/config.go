package configbase

import (
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/manager"
)

type Configurator interface {
	Optimize(configPath string) error
	Build() error
	GetLogger() logger.Logger
	GetClientPool() client.Pool
	GetManager() manager.Manager
}
