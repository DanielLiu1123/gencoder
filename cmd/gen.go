package cmd

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
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
	defaultOutputMarker = "gencoder generated file:"
)

var (
	config *string
)

func init() {
	config = genCmd.Flags().StringP("config", "f", "gencoder.yaml", "config file to use")
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate code from database metadata",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := readConfig()
		if err != nil {
			panic(err)
		}

		templates, err := loadTemplates(cfg)
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
				var schema string
				if tbCfg.Schema == "" {
					schema = dbCfg.Schema
				} else {
					schema = tbCfg.Schema
				}
				table, err := info.GenMySQLTable(context.Background(), db, schema, tbCfg.Name)
				if err != nil {
					panic(err)
				}

				var ctx renderCtx
				ctx.Table = table
				ctx.Properties = make(map[string]string)
				for k, v := range dbCfg.Properties {
					ctx.Properties[k] = v
				}
				for k, v := range tbCfg.Properties {
					ctx.Properties[k] = v
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

					if _, err := os.Stat(fileName); err == nil {
						// File exists, replace specific block
						err = replaceBlockInFile(fileName, content, tpl.BlockID)
						if err != nil {
							panic(err)
						}
					} else {
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

func loadTemplates(cfg *info.Config) ([]*tpl, error) {
	var templates []*tpl

	err := filepath.WalkDir(cfg.TemplatesDir, func(path string, d fs.DirEntry, err error) error {
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
			GeneratedFileName: getFileNameTemplate(&content, cfg),
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

func getFileNameTemplate(content *string, cfg *info.Config) string {
	scanner := bufio.NewScanner(strings.NewReader(*content))

	var generatedFileName string
	if cfg.OutputMarker != "" {
		generatedFileName = cfg.OutputMarker
	} else {
		generatedFileName = defaultOutputMarker
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, generatedFileName) {
			return strings.TrimSpace(line[strings.LastIndex(line, generatedFileName)+len(generatedFileName):])
		}
	}

	// Maybe a partial template
	return ""
}

func replaceBlockInFile(fileName, newContent, blockID string) error {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	startTag := fmt.Sprintf("gencoder block start: %s", blockID)
	endTag := fmt.Sprintf("gencoder block end: %s", blockID)

	startIndex := bytes.Index(fileData, []byte(startTag))
	endIndex := bytes.Index(fileData, []byte(endTag))
	if startIndex == -1 || endIndex == -1 || startIndex > endIndex {
		return errors.New("block not found or malformed")
	}

	// Replace the content between the start and end tags
	newFileData := append(fileData[:startIndex+len(startTag)], append([]byte(newContent), fileData[endIndex:]...)...)

	return os.WriteFile(fileName, newFileData, 0644)
}

type tpl struct {
	TemplateName      string            // template file name
	GeneratedFileName string            // generated file name, if empty, it's a partial template
	Source            string            // template source code
	Template          *raymond.Template // compiled template
}

type renderCtx struct {
	Table      *info.Table
	Properties map[string]string
}
