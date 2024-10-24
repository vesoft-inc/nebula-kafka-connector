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

func init() {
	validateRules = map[string]validateRule{
		"path": validatePath,
		"tls":  validateTLS,
	}
}
