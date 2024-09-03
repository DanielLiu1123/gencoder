package root

import (
	"github.com/DanielLiu1123/gencoder/pkg/cmd/gen"
	"github.com/DanielLiu1123/gencoder/pkg/util"
	"github.com/aymerick/raymond"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

func NewCmdRoot() *cobra.Command {

	registerHelperFunctions()

	c := &cobra.Command{
		Use:   "gencoder <command> [flags]",
		Short: "gencoder short",
		Long:  `gencoder longlonglong`,
	}

	c.AddCommand(gen.NewCmdGen())

	return c
}

func registerHelperFunctions() {

	raymond.RegisterHelper("replaceAll", func(target, old, new string) string {
		return strings.ReplaceAll(target, old, new)
	})

	raymond.RegisterHelper("match", func(target, pattern string) bool {
		match, err := regexp.MatchString(pattern, target)
		if err != nil {
			panic(err)
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

}
