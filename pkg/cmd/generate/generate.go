package generate

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/DanielLiu1123/gencoder/pkg/handlebars"
	"github.com/DanielLiu1123/gencoder/pkg/jsruntime"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/DanielLiu1123/gencoder/pkg/util"
	"github.com/spf13/cobra"
)

type generateOptions struct {
	config                string
	importHelper          string
	commandLineProperties map[string]string
	commandLineTemplates  string
}

func NewCmdGenerate(globalOptions *model.GlobalOptions) *cobra.Command {
	opt := &generateOptions{}
	var props []string

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

  # Generate boilerplate code from URL with custom properties
  $ gencoder generate --templates "https://github.com/DanielLiu1123/gencoder/tree/main/templates" --properties="package=com.example,author=Freeman"
`,
		PreRun: func(cmd *cobra.Command, args []string) {
			validateArgs(args)
			opt.commandLineProperties = parseProperties(props)
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
	c.Flags().StringSliceVarP(&props, "properties", "p", []string{}, "Add properties, will override properties in config file, --properties=\"k1=v1\" --properties=\"k2=v2,k3=v3\"")
	c.Flags().StringVarP(&opt.commandLineTemplates, "templates", "t", "", "Override templates directory, can be path or URL, e.g. https://github.com/DanielLiu1123/gencoder/tree/main/templates")

	return c
}

func validateArgs(args []string) {
	if len(args) > 0 {
		log.Fatalf("generate command does not accept any arguments")
	}
}

func parseProperties(props []string) map[string]string {
	properties := make(map[string]string)
	for _, prop := range props {
		parts := strings.Split(prop, "=")
		if len(parts) != 2 {
			log.Fatalf("Invalid property: %s", prop)
		}
		properties[parts[0]] = parts[1]
	}
	return properties
}

func registerCustomHelpers(location string) {
	content := fetchHelperContent(location)
	jsruntime.RunJS(content)
}

func fetchHelperContent(location string) string {
	if strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://") {
		return fetchFromURL(location)
	}
	return fetchFromFile(location)
}

func fetchFromURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func fetchFromFile(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func run(_ *cobra.Command, _ []string, opt *generateOptions, _ *model.GlobalOptions) {
	cfg, err := loadConfig(opt)
	if err != nil {
		log.Fatal(err)
	}

	templates, err := util.LoadTemplates(cfg, opt.commandLineTemplates)
	if err != nil {
		log.Fatal(err)
	}

	registerPartialTemplates(templates)

	renderContexts := util.CollectRenderContexts(cfg, opt.commandLineProperties)

	if len(renderContexts) > 0 {
		generateForAllContexts(cfg, templates, renderContexts)
	} else {
		generateFromBoilerplate(cfg, templates, opt.commandLineProperties)
	}
}

func loadConfig(opt *generateOptions) (*model.Config, error) {
	cfg, err := util.ReadConfig(opt.config)
	if (err != nil && !errors.Is(err, os.ErrNotExist)) || (errors.Is(err, os.ErrNotExist) && opt.commandLineTemplates == "") {
		return nil, err
	}
	return cfg, nil
}

func registerPartialTemplates(templates []*model.Tpl) {
	for _, t := range templates {
		if t.GeneratedFileName == "" {
			handlebars.RegisterPartial(t.TemplateName, t.Source)
		}
	}
}

func generateForAllContexts(cfg *model.Config, templates []*model.Tpl, renderContexts []*model.RenderContext) {
	for _, ctx := range renderContexts {
		for _, t := range templates {
			generate(cfg, t, ctx)
		}
	}
}

func generateFromBoilerplate(cfg *model.Config, templates []*model.Tpl, commandLineProperties map[string]string) {
	properties := mergeProperties(cfg.Properties, commandLineProperties)
	renderContext := &model.RenderContext{Properties: properties, Config: cfg}
	for _, t := range templates {
		generate(cfg, t, renderContext)
	}
}

func mergeProperties(configProps map[string]string, cmdLineProps map[string]string) map[string]string {
	merged := make(map[string]string)
	for k, v := range configProps {
		merged[k] = v
	}
	for k, v := range cmdLineProps {
		merged[k] = v
	}
	return merged
}

func generate(cfg *model.Config, tpl *model.Tpl, ctx *model.RenderContext) {
	if tpl.GeneratedFileName == "" { // partial template
		return
	}

	context := util.ToMap(ctx)
	newContent := handlebars.Render(tpl.Template, context)
	fileName := getFileName(tpl.GeneratedFileName, context)

	if _, err := os.Stat(fileName); err == nil {
		updateExistingFile(cfg, fileName, newContent)
	} else {
		createNewFile(fileName, newContent)
	}
}

func updateExistingFile(cfg *model.Config, fileName, newContent string) {
	oldContent, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	realContent := replaceBlocks(cfg, string(oldContent), newContent)
	writeFile(fileName, realContent)
}

func createNewFile(fileName, content string) {
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
				writeBlock(&realContent, newBlocks, currentBlockID, &currentBlock)
			}
			currentBlockID = extractBlockID(trimmed, cfg.BlockMarker.GetStart())
			currentBlock.WriteString(line + "\n")
		} else if strings.Contains(trimmed, cfg.BlockMarker.GetEnd()) && currentBlockID != "" {
			currentBlock.WriteString(line + "\n")
			writeBlock(&realContent, newBlocks, currentBlockID, &currentBlock)
			currentBlockID = ""
			currentBlock.Reset()
		} else if currentBlockID != "" {
			currentBlock.WriteString(line + "\n")
		} else {
			realContent.WriteString(line + "\n")
		}
	}

	if currentBlockID != "" {
		writeBlock(&realContent, newBlocks, currentBlockID, &currentBlock)
	}

	return strings.TrimSuffix(realContent.String(), "\n")
}

func writeBlock(realContent *strings.Builder, newBlocks map[string]string, blockID string, currentBlock *strings.Builder) {
	if newBlock, exists := newBlocks[blockID]; exists {
		realContent.WriteString(newBlock + "\n")
	} else {
		realContent.WriteString(strings.TrimSuffix(currentBlock.String(), "\n") + "\n")
	}
}

func extractBlockID(line, marker string) string {
	return strings.TrimSpace(line[strings.Index(line, marker)+len(marker):])
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
				blocks[currentBlockID] = strings.TrimSuffix(currentBlock.String(), "\n")
				currentBlock.Reset()
			}
			currentBlockID = extractBlockID(trimmed, cfg.BlockMarker.GetStart())
			currentBlock.WriteString(line + "\n")
		} else if strings.Contains(trimmed, cfg.BlockMarker.GetEnd()) && currentBlockID != "" {
			currentBlock.WriteString(line + "\n")
			blocks[currentBlockID] = strings.TrimSuffix(currentBlock.String(), "\n")
			currentBlockID = ""
			currentBlock.Reset()
		} else if currentBlockID != "" {
			currentBlock.WriteString(line + "\n")
		}
	}

	if currentBlockID != "" {
		blocks[currentBlockID] = strings.TrimSuffix(currentBlock.String(), "\n")
	}

	return blocks
}
