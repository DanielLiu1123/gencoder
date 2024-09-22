package introspect

import (
	"fmt"
	"log"
	"sync"

	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/DanielLiu1123/gencoder/pkg/util"
	"github.com/spf13/cobra"
)

type introspectOptions struct {
	config string
	output string
}

func NewCmdIntrospect(globalOptions *model.GlobalOptions) *cobra.Command {
	opt := &introspectOptions{}

	cmd := &cobra.Command{
		Use:     "introspect",
		Short:   "Print table information from database configuration",
		Aliases: []string{"intro", "i"},
		Example: `  # Print metadata of database tables from default config file (gencoder.yaml)
  $ gencoder introspect

  # Print metadata of database tables from a specific config file
  $ gencoder introspect -f myconfig.yaml
  
  # Print metadata of database tables from a specific config file in JSON/YAML format
  $ gencoder introspect -f myconfig.yaml -o [json|yaml]`,
		PreRun: validateArgs,
		Run:    func(cmd *cobra.Command, args []string) { run(cmd, args, opt, globalOptions) },
	}

	cmd.Flags().StringVarP(&opt.config, "config", "f", globalOptions.Config, "Config file to use")
	cmd.Flags().StringVarP(&opt.output, "output", "o", "json", "Output format, one of (json, yaml)")
	_ = cmd.RegisterFlagCompletionFunc("output", completeOutputFormat)

	return cmd
}

func validateArgs(_ *cobra.Command, args []string) {
	if len(args) > 0 {
		log.Fatalf("introspect command does not accept any arguments")
	}
}

func completeOutputFormat(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
}

func run(_ *cobra.Command, _ []string, opt *introspectOptions, _ *model.GlobalOptions) {
	renderContexts := sync.OnceValue(func() []*model.RenderContext {
		cfg, err := util.ReadConfig(opt.config)
		if err != nil {
			log.Fatalf("failed to read config: %v", err)
		}
		return util.CollectRenderContexts(cfg, nil)
	})

	switch opt.output {
	case "json":
		fmt.Println(util.ToJson(renderContexts()))
	case "yaml", "yml":
		fmt.Println(util.ToYaml(renderContexts()))
	default:
		log.Fatalf("unsupported output format: %s, must be one of (json, yaml)", opt.output)
	}
}
