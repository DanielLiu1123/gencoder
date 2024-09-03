package root

import (
	"github.com/DanielLiu1123/gencoder/pkg/cmd/gen"
	"github.com/aymerick/raymond"
	"github.com/spf13/cobra"
	"reflect"
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
	// eq
	raymond.RegisterHelper("eq", func(left, right string) bool {
		return reflect.DeepEqual(left, right)
	})
}
