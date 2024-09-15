package cmd

import (
	"github.com/DanielLiu1123/gencoder/pkg/cmd/generate"
	"github.com/DanielLiu1123/gencoder/pkg/cmd/introspect"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/spf13/cobra"
)

func NewCmdRoot() *cobra.Command {

	opt := &model.GlobalOptions{}

	c := &cobra.Command{
		Use:   "gencoder <command> [flags]",
		Short: "The ultimate code generator",
		Long: `
gencoder is a code generator that generates code from templates/databases, for any language or framework.

$ gencoder generate -f gencoder.yaml
$ gencoder introspect -f gencoder.yaml`,
	}

	c.Flags().StringVarP(&opt.Config, "config", "f", "gencoder.yaml", "Config file to use")

	c.AddCommand(generate.NewCmdGenerate(opt))
	c.AddCommand(introspect.NewCmdIntrospect(opt))

	return c
}
