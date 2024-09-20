package introspect

import (
	"fmt"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/DanielLiu1123/gencoder/pkg/util"
	"github.com/spf13/cobra"
	"log"
	"sync"
)

type introspectOptions struct {
	config string
	output string
}

func NewCmdIntrospect(globalOptions *model.GlobalOptions) *cobra.Command {

	opt := &introspectOptions{}

	c := &cobra.Command{
		Use:     "introspect",
		Short:   "Print table information from database configuration",
		Aliases: []string{"intro", "i"},
		Example: `  # Print metadata of database tables from default config file (gencoder.yaml)
  $ gencoder introspect

  # Print metadata of database tables from a specific config file
  $ gencoder introspect -f myconfig.yaml
  
  # Print metadata of database tables from a specific config file in JSON/YAML format
  $ gencoder introspect -f myconfig.yaml -o [json|yaml]
`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				log.Fatalf("introspect command does not accept any arguments")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, opt, globalOptions)
		},
	}

	c.Flags().StringVarP(&opt.config, "config", "f", globalOptions.Config, "Config file to use")
	c.Flags().StringVarP(&opt.output, "output", "o", "json", "Output format, one of (json, yaml)")
	_ = c.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
	})

	return c
}

func run(_ *cobra.Command, _ []string, opt *introspectOptions, _ *model.GlobalOptions) {

	renderContextsFunc := sync.OnceValue(func() []*model.RenderContext {
		cfg := util.ReadConfig(opt.config)
		return util.CollectRenderContexts(cfg.Databases...)
	})

	switch opt.output {
	case "json":
		fmt.Println(util.ToJson(renderContextsFunc()))
	case "yaml", "yml":
		fmt.Println(util.ToYaml(renderContextsFunc()))
	default:
		log.Fatalf("unsupported output format: %s, must be one of (json, yaml)", opt.output)
	}
}
