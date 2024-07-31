package cmd

import (
	"bufio"
	"context"
	"github.com/DanielLiu1123/gencoder/info"
	"github.com/aymerick/raymond"
	"github.com/spf13/cobra"
	"github.com/xo/dburl"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

const (
	fileNamePrefix = "gencoder generated file:"
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

		templates, err := loadTemplates(cfg.TemplatesDir)

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

				for _, tpl := range templates {
					result, err := tpl.Exec(table)
					if err != nil {
						panic(err)
					}

					filename := table.Name + ".go"
					err = os
				}

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

func loadTemplates(dir string) ([]*tpl, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var templates = make([]*tpl, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasPrefix(entry.Name(), ".hbs") {
			continue
		}

		b, err := os.ReadFile(entry.Name())
		if err != nil {
			return nil, err
		}

		content := string(b)

		template, err := raymond.Parse(content)
		if err != nil {
			return nil, err
		}

		t := &tpl{
			Name:     getGeneratedFileName(&content),
			Template: template,
		}

		templates = append(templates, t)
	}

	return templates, nil
}

func getGeneratedFileName(content *string) string {
	scanner := bufio.NewScanner(strings.NewReader(*content))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, fileNamePrefix) {
			return strings.TrimSpace(line[strings.LastIndex(line, fileNamePrefix)+len(fileNamePrefix):])
		}
	}

	return ""
}

type tpl struct {
	Name     string
	Template *raymond.Template
}
