package job

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/config"
)

func getServiceConfig(c *config.Config, serviceGroup string, serviceType string, host string) (map[string]string, error) {
	configs := make(map[string]string)
	var customizeConfig map[string]any
	if serviceType == "metad" {
		customizeConfig = c.Spec.Metad.Config
	} else {
		var svcgrp *config.ServiceGroupSpec
		for _, sg := range c.Spec.ServiceGroups {
			if sg.Name == serviceGroup {
				svcgrp = sg
				break
			}
		}
		if svcgrp == nil {
			return nil, fmt.Errorf("service group %s not found", serviceGroup)
		}
		switch serviceType {
		case "graphd":
			customizeConfig = svcgrp.Graphd.Config
		case "storaged":
			customizeConfig = svcgrp.Storaged.Config
		default:
			return nil, fmt.Errorf("invalid service type %s", serviceType)
		}
	}
	for k, v := range customizeConfig {
		configs[k] = fmt.Sprintf("%v", v)
	}
	configs["local_ip"] = host
	configs["meta_server_addrs"] = c.GetMetadAddress()
	return configs, nil
}

func isValidServiceType(serviceType string) bool {
	return serviceType == "graphd" || serviceType == "storaged" || serviceType == "metad"
}
