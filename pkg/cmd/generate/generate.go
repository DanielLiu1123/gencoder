package generate

import (
	"github.com/DanielLiu1123/gencoder/pkg/handlebars"
	"github.com/DanielLiu1123/gencoder/pkg/jsruntime"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/DanielLiu1123/gencoder/pkg/util"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type generateOptions struct {
	config       string
	importHelper string
}

func NewCmdGenerate(globalOptions *model.GlobalOptions) *cobra.Command {

	opt := &generateOptions{}

	c := &cobra.Command{
		Use:     "generate",
		Short:   "Generate code from database configuration",
		Aliases: []string{"gen", "g"},
		Example: `  # Generate code from default config file (gencoder.yaml)
  $ gencoder generate

  # Generate code from a specific config file
  $ gencoder generate -f myconfig.yaml

  # Generate code with custom import helper JavaScript file
  $ gencoder generate -f myconfig.yaml --import-helper helpers.js
`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				log.Fatalf("generate command does not accept any arguments")
			}
			if opt.importHelper != "" {
				registerCustomHelpers(opt.importHelper)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, opt, globalOptions)
		},
	}

	c.Flags().StringVarP(&opt.config, "config", "f", globalOptions.Config, "Config file to use")
	c.Flags().StringVarP(&opt.importHelper, "import-helper", "i", "", "Import helper JavaScript file, can be URL ([http|https]://...) or file path")

	return c
}

func registerCustomHelpers(location string) {

	// URL
	if strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://") {
		resp, err := http.Get(location)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		jsruntime.RunJS(string(bytes))

		return
	}

	// file path
	bytes, err := os.ReadFile(location)
	if err != nil {
		log.Fatal(err)
	}
	jsruntime.RunJS(string(bytes))
}

func run(_ *cobra.Command, _ []string, opt *generateOptions, _ *model.GlobalOptions) {

	cfg := util.ReadConfig(opt.config)

	templates, err := util.LoadTemplates(cfg)
	if err != nil {
		log.Fatal(err)
	}

	registerPartialTemplates(templates)

	renderContexts := util.CollectRenderContexts(cfg.Databases...)
	for _, ctx := range renderContexts {
		for _, t := range templates {
			generate(cfg, t, ctx)
		}
	}
}

func registerPartialTemplates(templates []*model.Tpl) {
	for _, t := range templates {
		if t.GeneratedFileName == "" {
			handlebars.RegisterPartial(t.TemplateName, t.Source)
		}
	}
}

func generate(cfg *model.Config, tpl *model.Tpl, ctx *model.RenderContext) {

	// Skip partial templates
	if tpl.GeneratedFileName == "" {
		return
	}

	context := util.ToMap(ctx)

	newContent := handlebars.Render(tpl.Template, context)

	fileName := getFileName(tpl.GeneratedFileName, context)

	if _, err := os.Stat(fileName); err == nil {

		oldContent, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatal(err)
		}

		realContent := replaceBlocks(cfg, string(oldContent), newContent)
		writeFile(fileName, realContent)

	} else {

		createFile(fileName, newContent)
	}
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

func getFileName(filenameTpl string, ctx map[string]interface{}) string {
	tpl := handlebars.Compile(filenameTpl)
	return handlebars.Render(tpl, ctx)
}

// Thanks to ChatGPT :)
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
			currentBlockID = strings.TrimSpace(trimmed[strings.Index(trimmed, cfg.BlockMarker.GetStart())+len(cfg.BlockMarker.GetStart()):])
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
			currentBlockID = strings.TrimSpace(strings.SplitN(trimmed, cfg.BlockMarker.GetStart(), 2)[1])
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
