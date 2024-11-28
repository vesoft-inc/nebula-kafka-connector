package config

import (
	"fmt"
	"strings"
)

type validateRule func(c *Config) error

var validateRules map[string]validateRule

func validatePath(c *Config) error {
	var l = make([]string, 0)
	if c.InstallPath == "" {
		l = append(l, "installPath is empty")
	}
	if c.PackagePath == "" {
		l = append(l, "packagePath is empty")
	}
	if len(l) > 0 {
		return fmt.Errorf("%s", strings.Join(l, ", "))
	}
	return nil
}

func validateTLS(c *Config) error {
	var l = make([]string, 0)
	if c.CertFile == "" {
		l = append(l, "certFile is empty")
	}
	if c.KeyFile == "" {
		l = append(l, "keyFile is empty")
	}
	if c.CaFile == "" {
		l = append(l, "caFile is empty")
	}
	if len(l) > 0 {
		return fmt.Errorf("%s", strings.Join(l, ", "))
	}
	return nil
}

func validateServiceGroup(c *Config) error {
	if c.Spec == nil || len(c.Spec.ServiceGroups) == 0 {
		return fmt.Errorf("serviceGroups is empty")
	}
	return nil
}

func validateBlankConfigItem(c *Config) error {
	if c.Spec == nil {
		return nil
	}
	if c.Spec.Metad.Config != nil {
		for k, v := range c.Spec.Metad.Config {
			if v == nil || v == "" {
				return fmt.Errorf("metad config item %s is empty", k)
			}
		}
	}
	for _, sg := range c.Spec.ServiceGroups {
		if sg.Graphd.Config != nil {
			for k, v := range sg.Graphd.Config {
				if v == nil || v == "" {
					return fmt.Errorf("graphd config item %s is empty", k)
				}
			}
		}
		if sg.Storaged.Config != nil {
			for k, v := range sg.Storaged.Config {
				if v == nil || v == "" {
					return fmt.Errorf("storaged config item %s is empty", k)
				}
			}
		}
	}
	return nil
}

func init() {
	validateRules = map[string]validateRule{
		"path":  validatePath,
		"tls":   validateTLS,
		"sg":    validateServiceGroup,
		"blank": validateBlankConfigItem,
	}
}
