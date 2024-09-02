package cmd

import (
	"bufio"
	"context"
	"database/sql"
	"github.com/DanielLiu1123/gencoder/info"
	"github.com/aymerick/raymond"
	"github.com/spf13/cobra"
	"github.com/xo/dburl"
	"gopkg.in/yaml.v3"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	config *string
)

func NewGenCmd() *cobra.Command {

	c := &cobra.Command{
		Use:   "gen",
		Short: "Generate code from database configuration",
		Run:   run,
	}

	c.Flags().StringVarP(config, "config", "f", "gencoder.yaml", "config file to use")

	return c
}

func run(_ *cobra.Command, _ []string) {
	cfg, err := readConfig(*config)
	if err != nil {
		log.Fatal(err)
	}

	templates, err := loadTemplates(cfg)
	if err != nil {
		log.Fatal(err)
	}

	registerPartialTemplates(templates)

	for _, dbCfg := range cfg.Databases {
		processDatabase(cfg, dbCfg, templates)
	}
}

func registerPartialTemplates(templates []*tpl) {
	for _, t := range templates {
		if t.GeneratedFileName == "" {
			raymond.RegisterPartial(t.TemplateName, t.Source)
		}
	}
}

func processDatabase(cfg *info.Config, dbCfg *info.DatabaseConfig, templates []*tpl) {
	db, err := dburl.Open(dbCfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	for _, tbCfg := range dbCfg.Tables {
		schema := tbCfg.Schema
		if schema == "" {
			schema = dbCfg.Schema
		}

		table, err := info.GenMySQLTable(context.Background(), db, schema, tbCfg.Name, tbCfg.IgnoreColumns)
		if err != nil {
			log.Fatal(err)
		}

		ctx := createRenderContext(dbCfg, tbCfg, table)

		for _, tpl := range templates {
			if tpl.GeneratedFileName == "" {
				continue
			}
			handleTemplate(cfg, tpl, ctx)
		}
	}
}

func createRenderContext(dbCfg *info.DatabaseConfig, tbCfg *info.TableConfig, table *info.Table) *renderCtx {
	properties := make(map[string]string)
	for k, v := range dbCfg.Properties {
		properties[k] = v
	}
	for k, v := range tbCfg.Properties {
		properties[k] = v
	}

	return &renderCtx{
		Table:      table,
		Properties: properties,
	}
}

func handleTemplate(cfg *info.Config, tpl *tpl, ctx *renderCtx) {
	newContent, err := tpl.Template.Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fileName := getFileName(tpl.GeneratedFileName, ctx)

	if fileExists(fileName) {
		oldContent, err := readFile(fileName)
		if err != nil {
			log.Fatal(err)
		}

		realContent := replaceBlocks(cfg, oldContent, newContent)
		writeFile(fileName, realContent)
	} else {
		createFile(fileName, newContent)
	}
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}

func createFile(fileName, content string) {
	dir := filepath.Dir(fileName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal(err)
	}
	writeFile(fileName, content)
}

func writeFile(fileName, content string) {
	if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
}

func getFileName(filenameTpl string, ctx *renderCtx) string {
	t, err := raymond.Parse(filenameTpl)
	if err != nil {
		log.Fatal(err)
	}

	fileName, err := t.Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return fileName
}

func readConfig(configPath string) (*info.Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg info.Config
	err = yaml.Unmarshal(file, &cfg)
	return &cfg, err
}

func loadTemplates(cfg *info.Config) ([]*tpl, error) {
	var templates []*tpl

	err := filepath.WalkDir(cfg.GetTemplatesDir(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".hbs") {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		template, err := raymond.Parse(string(b))
		if err != nil {
			return err
		}

		t := &tpl{
			TemplateName:      d.Name(),
			GeneratedFileName: getFileNameTemplate(string(b), cfg),
			Source:            string(b),
			Template:          template,
		}

		templates = append(templates, t)
		return nil
	})

	return templates, err
}

func getFileNameTemplate(content string, cfg *info.Config) string {
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, cfg.GetOutputMarker()) {
			return strings.TrimSpace(line[strings.LastIndex(line, cfg.GetOutputMarker())+len(cfg.GetOutputMarker()):])
		}
	}

	return ""
}

func readFile(filename string) (string, error) {
	fileData, err := os.ReadFile(filename)
	return string(fileData), err
}

func parseBlocks(cfg *info.Config, content string) map[string]string {
	blocks := make(map[string]string)
	lines := strings.Split(content, "\n")
	var currentBlockID string
	var currentBlock strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, cfg.BlockMarker.GetStart()) {
			if currentBlockID != "" {
				blocks[currentBlockID] = strings.TrimRight(currentBlock.String(), "\n")
				currentBlock.Reset()
			}
			currentBlockID = strings.SplitN(trimmed, cfg.BlockMarker.GetStart(), 2)[1]
			currentBlock.WriteString(line + "\n")
		} else if strings.Contains(trimmed, cfg.BlockMarker.GetEnd()) && currentBlockID != "" {
			currentBlock.WriteString(line + "\n")
			blocks[currentBlockID] = strings.TrimRight(currentBlock.String(), "\n")
			currentBlockID = ""
			currentBlock.Reset()
		} else if currentBlockID != "" {
			currentBlock.WriteString(line + "\n")
		}
	}

	if currentBlockID != "" {
		blocks[currentBlockID] = strings.TrimRight(currentBlock.String(), "\n")
	}

	return blocks
}

func replaceBlocks(cfg *info.Config, oldContent, newContent string) string {
	newBlocks := parseBlocks(cfg, newContent)
	var realContent strings.Builder

	lines := strings.Split(oldContent, "\n")
	var currentBlockID string
	var currentBlock strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, cfg.BlockMarker.GetStart()) {
			if currentBlockID != "" {
				if newBlock, exists := newBlocks[currentBlockID]; exists {
					realContent.WriteString(newBlock + "\n")
				} else {
					realContent.WriteString(strings.TrimRight(currentBlock.String(), "\n") + "\n")
				}
				currentBlock.Reset()
			}
			currentBlockID = strings.SplitN(trimmed, cfg.BlockMarker.GetStart(), 2)[1]
			currentBlock.WriteString(line + "\n")
		} else if strings.Contains(trimmed, cfg.BlockMarker.GetEnd()) && currentBlockID != "" {
			currentBlock.WriteString(line + "\n")
			if newBlock, exists := newBlocks[currentBlockID]; exists {
				realContent.WriteString(newBlock + "\n")
			} else {
				realContent.WriteString(strings.TrimRight(currentBlock.String(), "\n") + "\n")
			}
			currentBlockID = ""
			currentBlock.Reset()
		} else if currentBlockID != "" {
			currentBlock.WriteString(line + "\n")
		} else {
			realContent.WriteString(strings.TrimRight(line, "\n") + "\n")
		}
	}

	if currentBlockID != "" {
		if newBlock, exists := newBlocks[currentBlockID]; exists {
			realContent.WriteString(newBlock + "\n")
		} else {
			realContent.WriteString(strings.TrimRight(currentBlock.String(), "\n") + "\n")
		}
	}

	return strings.TrimRight(realContent.String(), "\n")
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
