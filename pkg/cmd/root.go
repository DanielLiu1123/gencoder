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
		Short:   "A code generator that keeps your changes during regeneration.",
		Long:    "A code generator that keeps your changes during regeneration, powered by Handlebars.",
		Example: `
  # Generate code from config file (default: gencoder.yaml), config json schema: https://raw.githubusercontent.com/DanielLiu1123/gencoder/refs/heads/main/schema.json
  $ gencoder generate -f gencoder.yaml

  # Generate code from a template project with custom properties
  $ gencoder generate --templates "https://github.com/user/template-project" --properties "package=com.example,author=Freeman" --include-non-tpl

  # Generate code using custom helpers, build-in helpers: https://github.com/DanielLiu1123/gencoder/blob/main/pkg/jsruntime/helper.js
  $ gencoder generate --helpers helpers.js

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
