package config

import "github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"

type Log struct {
	Level   *string       `yaml:"level,omitempty"`
	Console *bool         `yaml:"console,omitempty"`
	Files   []string      `yaml:"files,omitempty"`
	Fields  logger.Fields `yaml:"fields,omitempty"`
}

func (l *Log) BuildLogger(opts ...logger.Option) (logger.Logger, error) {
	options := make([]logger.Option, 0, 4+len(opts))
	if l != nil {
		if l.Level != nil && *l.Level != "" {
			options = append(options, logger.WithLevelText(*l.Level))
		}
		if l.Console != nil {
			options = append(options, logger.WithConsole(*l.Console))
		}
		if len(l.Files) > 0 {
			options = append(options, logger.WithFiles(l.Files...))
		}
		if len(l.Fields) > 0 {
			options = append(options, logger.WithFields(l.Fields...))
		}
	}
	options = append(options, opts...)
	return logger.New(options...)
}
