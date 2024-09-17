package cmd

import (
	"github.com/DanielLiu1123/gencoder/pkg/cmd/generate"
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
		Example: `  # Generate code from default config file (gencoder.yaml)
  $ gencoder generate

  # Generate code from a specific config file
  $ gencoder generate -f myconfig.yaml`,
	}

	c.Flags().StringVarP(&opt.Config, "config", "f", "gencoder.yaml", "Config file to use")

	c.AddCommand(generate.NewCmdGenerate(opt))
	c.AddCommand(introspect.NewCmdIntrospect(opt))

	return c
}
