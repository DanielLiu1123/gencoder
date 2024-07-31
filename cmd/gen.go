package cmd

import (
	"bufio"
	"context"
	"github.com/DanielLiu1123/gencoder/info"
	"github.com/aymerick/raymond"
	"github.com/spf13/cobra"
	"github.com/xo/dburl"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
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
		if err != nil {
			panic(err)
		}

		for _, dbCfg := range cfg.Databases {
			db, err := dburl.Open(dbCfg.Dsn)
			if err != nil {
				panic(err)
			}

			for _, table := range dbCfg.Tables {
				table, err := info.GenMySQLTable(context.Background(), db, "testdb", table.Name)
				if err != nil {
					panic(err)
				}

				for _, tpl := range templates {
					content, err := tpl.Template.Exec(table)
					if err != nil {
						panic(err)
					}

					err = os.WriteFile(tpl.Name, []byte(content), 0644)
					if err != nil {
						panic(err)
					}
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
	var templates []*tpl

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".hbs") {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content := string(b)

		template, err := raymond.Parse(content)
		if err != nil {
			return err
		}

		t := &tpl{
			Name:     getGeneratedFileName(&content),
			Template: template,
		}

		templates = append(templates, t)
		return nil
	})

	if err != nil {
		return nil, err
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

	panic(fileNamePrefix + " not found")
}

type tpl struct {
	Name     string
	Template *raymond.Template
}
