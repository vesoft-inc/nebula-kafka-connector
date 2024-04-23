package config

import (
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"
	pkgconfig "github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/config"
)

func (s *configService) SetComponentConfig(req *types.SetComponentConfigReq) (resp *types.SetComponentConfigResp, err error) {
	config, err := pkgconfig.LoadConfig()
	if err != nil {
		return nil, err
	}

	config.SetConfig(req.Component, req.Config)

	err = pkgconfig.SaveConfig(config)
	if err != nil {
		return nil, err
	}

	return &types.SetComponentConfigResp{}, nil
}
