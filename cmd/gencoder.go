package cmd

import (
	"github.com/aymerick/raymond"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

var gencoderCmd = &cobra.Command{
	Use:   "gencoder",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {

	registerHelperFunctions()

	gencoderCmd.AddCommand(NewGenCmd())
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
}

func Execute() {
	err := gencoderCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
