package cmd

import (
	"github.com/DanielLiu1123/gencoder/pkg/cmd/generate"
	initCmd "github.com/DanielLiu1123/gencoder/pkg/cmd/init"
	"github.com/DanielLiu1123/gencoder/pkg/cmd/introspect"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/spf13/cobra"
)

func NewCmdRoot(buildInfo *model.BuildInfo) *cobra.Command {

	opt := &model.GlobalOptions{}

	c := &cobra.Command{
		Use:     "gencoder",
		Version: buildInfo.Version,
		Short:   "The ultimate code generator",
		Long:    "gencoder is a code generator that generates code from templates/databases, for any languages/frameworks.",
		Example: `  # Generate code from config file (default: gencoder.yaml)
  $ gencoder generate -f gencoder.yaml

  # Init basic config for quick start
  $ gencoder init

  # Print metadata of database tables
  $ gencoder introspect -f gencoder.yaml -o yaml`,
	}

	c.Flags().StringVarP(&opt.Config, "config", "f", "gencoder.yaml", "Config file to use")

	c.AddCommand(generate.NewCmdGenerate(opt))
	c.AddCommand(introspect.NewCmdIntrospect(opt))
	c.AddCommand(initCmd.NewCmdInit(opt))

	return c
}
