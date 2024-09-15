package cmd

import (
	"github.com/DanielLiu1123/gencoder/pkg/cmd/generate"
	"github.com/DanielLiu1123/gencoder/pkg/cmd/introspect"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/DanielLiu1123/gencoder/pkg/util"
	"github.com/mailgun/raymond/v2"
	"github.com/spf13/cobra"
	"log"
	"regexp"
	"strings"
)

func NewCmdRoot() *cobra.Command {

	registerHelperFunctions()

	opt := &model.GlobalOptions{}

	c := &cobra.Command{
		Use:   "gencoder <command> [flags]",
		Short: "gencoder short",
		Long:  `gencoder longlonglong`,
	}

	c.Flags().StringVarP(&opt.Config, "config", "f", "gencoder.yaml", "Config file to use")

	c.AddCommand(generate.NewCmdGenerate(opt))
	c.AddCommand(introspect.NewCmdIntrospect(opt))

	return c
}

func registerHelperFunctions() {

	raymond.RegisterHelper("replaceAll", func(target, old, new string) string {
		return strings.ReplaceAll(target, old, new)
	})

	raymond.RegisterHelper("match", func(pattern, target string) bool {
		match, err := regexp.MatchString(pattern, target)
		if err != nil {
			log.Fatal(err)
		}
		return match
	})

	raymond.RegisterHelper("eq", func(left, right string) bool {
		return left == right
	})

	raymond.RegisterHelper("ne", func(left, right string) bool {
		return left != right
	})

	raymond.RegisterHelper("snakeCase", func(s string) string {
		return util.ToSnakeCase(s)
	})

	raymond.RegisterHelper("camelCase", func(s string) string {
		return util.ToCamelCase(s)
	})

	raymond.RegisterHelper("pascalCase", func(s string) string {
		return util.ToPascalCase(s)
	})

	raymond.RegisterHelper("upperFirst", func(s string) string {
		if len(s) == 0 {
			return ""
		}
		return strings.ToUpper(string(s[0])) + s[1:]
	})

	raymond.RegisterHelper("lowerFirst", func(s string) string {
		if len(s) == 0 {
			return ""
		}
		return strings.ToLower(string(s[0])) + s[1:]
	})

	raymond.RegisterHelper("uppercase", func(s string) string {
		return strings.ToUpper(s)
	})

	raymond.RegisterHelper("lowercase", func(s string) string {
		return strings.ToLower(s)
	})

	raymond.RegisterHelper("trim", func(s string) string {
		return strings.TrimSpace(s)
	})

	raymond.RegisterHelper("removePrefix", func(s, prefix string) string {
		return strings.TrimPrefix(s, prefix)
	})

	raymond.RegisterHelper("removeSuffix", func(s, suffix string) string {
		return strings.TrimSuffix(s, suffix)
	})

}
