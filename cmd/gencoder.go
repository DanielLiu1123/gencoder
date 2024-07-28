package cmd

import (
	"github.com/spf13/cobra"
	"os"
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
	gencoderCmd.AddCommand(genCmd)
}

func Execute() {
	err := gencoderCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
