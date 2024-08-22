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

		// register partials
		for _, t := range templates {
			if t.GeneratedFileName == "" {
				raymond.RegisterPartial(t.TemplateName, t.Source)
			}
		}

		for _, dbCfg := range cfg.Databases {
			db, err := dburl.Open(dbCfg.Dsn)
			if err != nil {
				panic(err)
			}
			defer db.Close()

			for _, tbCfg := range dbCfg.Tables {
				var scm string
				if tbCfg.Schema == "" {
					scm = dbCfg.Schema
				} else {
					scm = tbCfg.Schema
				}
				table, err := info.GenMySQLTable(context.Background(), db, scm, tbCfg.Name)
				if err != nil {
					panic(err)
				}

				var ctx renderCtx
				ctx.Table = table
				ctx.Properties = make(map[string]string)
				for _, prop := range dbCfg.Properties {
					ctx.Properties[prop.Name] = prop.Value
				}
				for _, prop := range tbCfg.Properties {
					ctx.Properties[prop.Name] = prop.Value
				}

				for _, tpl := range templates {

					// Maybe a partial template
					if tpl.GeneratedFileName == "" {
						continue
					}

					content, err := tpl.Template.Exec(ctx)
					if err != nil {
						panic(err)
					}

					fileName := getFileName(tpl.GeneratedFileName, &ctx)

					dir := filepath.Dir(fileName)
					if err = os.MkdirAll(dir, 0755); err != nil {
						panic(err)
					}

					err = os.WriteFile(fileName, []byte(content), 0644)
					if err != nil {
						panic(err)
					}
				}

			}
		}
	},
}

func getFileName(filenameTpl string, ctx *renderCtx) string {
	t, err := raymond.Parse(filenameTpl)
	if err != nil {
		panic(err)
	}

	fileName, err := t.Exec(ctx)
	if err != nil {
		panic(err)
	}

	return fileName
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
			TemplateName:      d.Name(),
			GeneratedFileName: getFileNameTemplate(&content),
			Source:            content,
			Template:          template,
		}

		templates = append(templates, t)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return templates, nil
}

func getFileNameTemplate(content *string) string {
	scanner := bufio.NewScanner(strings.NewReader(*content))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, fileNamePrefix) {
			return strings.TrimSpace(line[strings.LastIndex(line, fileNamePrefix)+len(fileNamePrefix):])
		}
	}

	// Maybe a partial template
	return ""
}

type tpl struct {
	TemplateName      string
	GeneratedFileName string
	Source            string
	Template          *raymond.Template
}

type renderCtx struct {
	Table      *info.Table
	Properties map[string]string
}
