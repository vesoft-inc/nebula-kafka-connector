package config

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"
	pkgconfig "github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/config"
)

func (s *configService) GetComponentConfig(req *types.GetComponentConfigReq) (resp *types.GetComponentConfigResp, err error) {
	config, err := pkgconfig.LoadConfig()
	if err != nil {
		return nil, err
	}

	component := config.GetComponent(req.Component)
	if component == nil {
		return nil, fmt.Errorf("invalid component: %s", req.Component)
	}

	return &types.GetComponentConfigResp{Config: component.Config}, nil
}
