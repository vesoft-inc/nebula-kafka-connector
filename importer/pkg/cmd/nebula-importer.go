package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/manager"
)

type (
	ImporterOptions struct {
		common.IOStreams
		Arguments    []string
		ConfigFile   string
		cfg          config.Config
		logger       logger.Logger
		useNopLogger bool // for test
		cli          client.Client
		mgr          manager.Manager
	}
)

func NewImporterOptions(streams common.IOStreams) *ImporterOptions {
	return &ImporterOptions{
		IOStreams: streams,
	}
}

func NewDefaultImporterCommand() *cobra.Command {
	o := NewImporterOptions(common.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	})
	return NewImporterCommand(o)
}

func NewImporterCommand(o *ImporterOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nebula-importer",
		Short: `The NebulaGraph Importer Tool.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				if err != nil {
					l := o.logger

					if l == nil || o.useNopLogger {
						l = logger.NopLogger
					}

					e := errors.NewImportError(err)
					fields := logger.MapToFields(e.Fields())
					l.SkipCaller(1).WithError(e.Cause()).Error("failed to execute", fields...)
				}
				if o.cli != nil {
					_ = o.cli.Close()
				}
				if o.logger != nil {
					_ = o.logger.Sync()
					_ = o.logger.Close()
				}
			}()
			err = o.Complete(cmd, args)
			if err != nil {
				return err
			}
			err = o.Validate()
			if err != nil {
				return err
			}
			return o.Run(cmd, args)
		},
	}
	o.AddFlags(cmd)
	return cmd
}

func (*ImporterOptions) Complete(_ *cobra.Command, _ []string) error {
	return nil
}

func (o *ImporterOptions) Validate() error {
	if err := o.cfg.FromFile(o.ConfigFile); err != nil {
		return err
	}
	cfg := o.cfg

	l, err := cfg.BuildLogger()
	if err != nil {
		return err
	}
	o.logger = l

	cli, err := cfg.BuildClient(client.WithLogger(l))
	if err != nil {
		return err
	}
	o.cli = cli

	mgr, err := cfg.BuildManager(cli, cfg.Sources, manager.WithLogger(l))
	if err != nil {
		return err
	}
	o.mgr = mgr
	return nil
}

func (o *ImporterOptions) Run(_ *cobra.Command, _ []string) error {
	if err := o.mgr.Start(); err != nil {
		return err
	}
	//revive:disable-next-line:if-return
	if err := o.mgr.Wait(); err != nil {
		return err
	}
	return nil
}

func (o *ImporterOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.ConfigFile, "config", "c", o.ConfigFile,
		"specify nebula-importer configure file")
}
