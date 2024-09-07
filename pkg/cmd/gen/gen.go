package gen

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"github.com/DanielLiu1123/gencoder/pkg/db"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/DanielLiu1123/gencoder/pkg/util"
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

type GenOptions struct {
	config       string
	contextsOnly bool // Only show render contexts
	output       string
}

func NewCmdGen() *cobra.Command {

	opt := &GenOptions{}

	c := &cobra.Command{
		Use:   "gen",
		Short: "Generate code from database configuration",
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, opt)
		},
	}

	c.Flags().StringVarP(&opt.config, "config", "f", "gencoder.yaml", "Config file to use")
	c.Flags().BoolVarP(&opt.contextsOnly, "contexts-only", "c", false, "Only show render contexts")
	c.Flags().StringVarP(&opt.output, "output", "o", "json", "Output format, only used with --contexts-only. One of: json, yaml")

	err := c.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "yaml"}, cobra.ShellCompDirectiveDefault
	})
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func run(_ *cobra.Command, _ []string, opt *GenOptions) {
	cfg, err := readConfig(opt.config)
	if err != nil {
		log.Fatal(err)
	}

	templates, err := loadTemplates(cfg)
	if err != nil {
		log.Fatal(err)
	}

	registerPartialTemplates(templates)

	var renderContexts []*model.RenderContext
	for _, dbCfg := range cfg.Databases {
		contexts := collectRenderContexts(dbCfg)
		renderContexts = append(renderContexts, contexts...)
	}

	if opt.contextsOnly {
		switch opt.output {
		case "json":
			jsonV, err := util.ToJson(renderContexts)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(jsonV)
		case "yaml":
			yamlV, err := util.ToYaml(renderContexts)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(yamlV)
		default:
			log.Fatalf("Invalid output format: %s", opt.output)
		}
		return
	}

	for _, ctx := range renderContexts {
		for _, tpl := range templates {
			handleTemplate(cfg, tpl, ctx)
		}
	}
}

func registerPartialTemplates(templates []*tpl) {
	for _, t := range templates {
		if t.GeneratedFileName == "" {
			raymond.RegisterPartial(t.TemplateName, t.Source)
		}
	}
}

func collectRenderContexts(dbCfg *model.DatabaseConfig) []*model.RenderContext {

	dbconn, err := dburl.Open(dbCfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(dbconn)

	var contexts []*model.RenderContext
	for _, tbCfg := range dbCfg.Tables {
		schema := tbCfg.Schema
		if schema == "" {
			schema = dbCfg.Schema
		}

		table, err := db.GenMySQLTable(context.Background(), dbconn, schema, tbCfg.Name, tbCfg.IgnoreColumns)
		if err != nil {
			log.Fatal(err)
		}

		ctx := createRenderContext(dbCfg, tbCfg, table)

		contexts = append(contexts, ctx)
	}

	return contexts
}

func createRenderContext(dbCfg *model.DatabaseConfig, tbCfg *model.TableConfig, table *model.Table) *model.RenderContext {
	properties := make(map[string]string)
	for k, v := range dbCfg.Properties {
		properties[k] = v
	}
	for k, v := range tbCfg.Properties {
		properties[k] = v
	}

	return &model.RenderContext{
		Table:      table,
		Properties: properties,
	}
}

func handleTemplate(cfg *model.Config, tpl *tpl, ctx *model.RenderContext) {

	// Skip partial templates
	if tpl.GeneratedFileName == "" {
		return
	}

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

func getFileName(filenameTpl string, ctx *model.RenderContext) string {
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

func readConfig(configPath string) (*model.Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg model.Config
	err = yaml.Unmarshal(file, &cfg)
	return &cfg, err
}

func loadTemplates(cfg *model.Config) ([]*tpl, error) {
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

func getFileNameTemplate(content string, cfg *model.Config) string {
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

func parseBlocks(cfg *model.Config, content string) map[string]string {
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

func replaceBlocks(cfg *model.Config, oldContent, newContent string) string {
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
