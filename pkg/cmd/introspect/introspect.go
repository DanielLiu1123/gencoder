package introspect

import (
	"fmt"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/DanielLiu1123/gencoder/pkg/util"
	"github.com/spf13/cobra"
	"log"
	"sync"
)

type IntrospectOptions struct {
	config string
	output string
}

func NewCmdIntrospect(globalOptions *model.GlobalOptions) *cobra.Command {

	opt := &IntrospectOptions{}

	c := &cobra.Command{
		Use:     "introspect",
		Short:   "Print table information from database configuration",
		Aliases: []string{"i", "intro"},
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

func run(_ *cobra.Command, _ []string, opt *IntrospectOptions, _ *model.GlobalOptions) {

	cfg, err := util.ReadConfig(opt.config)
	if err != nil {
		log.Fatal(err)
	}

	renderContextsFunc := sync.OnceValue(func() []*model.RenderContext {
		return util.CollectRenderContexts(cfg.Databases...)
	})

	switch opt.output {
	case "json":
		jsonValue, err := util.ToJson(renderContextsFunc())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(jsonValue)
	case "yaml", "yml":
		yamlValue, err := util.ToYaml(renderContextsFunc())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(yamlValue)
	default:
		log.Fatalf("unsupported output format: %s, must be one of (json, yaml)", opt.output)
	}
}
