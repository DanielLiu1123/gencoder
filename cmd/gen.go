package cmd

import (
	"context"
	"fmt"
	"github.com/DanielLiu1123/gencoder/info"
	"github.com/spf13/cobra"
	"github.com/xo/dburl"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	config *string
)

func init() {
	config = genCmd.Flags().StringP("config", "c", "gencoder.yaml", "config file to use")
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate code from database metadata",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := readConfig()
		if err != nil {
			panic(err)
		}

		for i, database := range cfg.Databases {
			db, err := dburl.Open(database.Dsn)
			if err != nil {
				panic(err)
			}

			for _, table := range database.Tables {
				table, err := info.GenMySQLTable(context.Background(), db, "testdb", table.Name)
				if err != nil {
					panic(err)
				}
				fmt.Printf("table %d: %v\n", i, table)
			}
		}
	},
}

func readConfig() (*info.Config, error) {
	file, err := os.ReadFile(*config)
	if err != nil {
		return nil, err
	}

	var config info.Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
