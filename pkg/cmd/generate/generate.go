package generate

import (
	"errors"
	"fmt"
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
	config        string
	importHelpers []string
	includeNonTpl bool

	// Override config file gencoder.yaml
	Templates  string
	Properties map[string]string // Add properties, will override properties in config file
	output     string
}

func NewCmdGenerate(globalOptions *model.GlobalOptions) *cobra.Command {
	opt := &generateOptions{}
	var props []string

	c := &cobra.Command{
		Use:     "generate",
		Short:   "Generate code from templates and database table metadata",
		Aliases: []string{"gen", "g"},
		Example: `
  # Generate code from config file (default: gencoder.yaml), config json schema: https://raw.githubusercontent.com/DanielLiu1123/gencoder/refs/heads/main/schema.json
  $ gencoder generate -f gencoder.yaml

  # Generate code from a template project with custom properties
  $ gencoder generate --templates "https://github.com/user/template-project" --properties "package=com.example,author=Freeman" --include-non-tpl

  # Generate code using custom helpers, build-in helpers: https://github.com/DanielLiu1123/gencoder/blob/main/pkg/jsruntime/helper.js
  $ gencoder generate --helpers helpers.js`,
		PreRun: func(cmd *cobra.Command, args []string) {
			validateArgs(args)
			opt.Properties = parseProperties(props)

			// Show deprecation warning if old flag is used
			if cmd.Flags().Changed("import-helpers") {
				fmt.Fprintf(os.Stderr, "Warning: --import-helpers flag is deprecated, please use --helpers instead\n")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, opt, globalOptions)
		},
	}

	c.Flags().StringVarP(&opt.config, "config", "f", globalOptions.Config, "Config file to use")
	c.Flags().StringSliceVar(&opt.importHelpers, "helpers", []string{}, "Import helper JavaScript file, can be URL ([http|https]://...) or file path")
	c.Flags().StringSliceVarP(&opt.importHelpers, "import-helpers", "i", []string{}, "Import helper JavaScript file, can be URL ([http|https]://...) or file path (deprecated, use --helpers instead)")
	c.Flags().StringSliceVarP(&props, "properties", "p", []string{}, "Add properties, will override properties in config file, --properties=\"k1=v1\" --properties=\"k2=v2,k3=v3\"")
	c.Flags().StringVarP(&opt.Templates, "templates", "t", "", "Override templates directory, can be path or URL, e.g. https://github.com/DanielLiu1123/gencoder/tree/main/templates")
	c.Flags().BoolVarP(&opt.includeNonTpl, "include-non-tpl", "a", false, "Include non-template files in the 'templates' option")
	c.Flags().StringVarP(&opt.output, "output", "o", "", "Output directory for generated files, default is the current directory")

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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

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

	mergeCmdOptionsToConfig(cfg, opt)

	// Register custom helpers
	for _, helper := range opt.importHelpers {
		registerCustomHelpers(helper)
	}
	for _, helper := range cfg.ImportHelpers {
		registerCustomHelpers(helper)
	}

	files, err := util.LoadFiles(cfg)
	if err != nil {
		log.Fatal(err)
	}

	registerPartials(files)

	renderContexts := util.CollectRenderContexts(cfg, opt.Properties)

	if opt.includeNonTpl {
		for _, f := range files {
			generateForNormalFiles(cfg, f)
		}
	}

	if len(renderContexts) > 0 {
		generateForAllContexts(cfg, files, renderContexts)
	} else {
		properties := mergeProperties(cfg.Properties, opt.Properties)
		renderContext := &model.RenderContext{Properties: properties, Config: cfg}
		for _, t := range files {
			generateForTemplateFiles(cfg, t, renderContext)
		}
	}
}

func mergeCmdOptionsToConfig(cfg *model.Config, opt *generateOptions) {
	if opt.output != "" {
		cfg.Output = opt.output
	}
	if opt.Templates != "" {
		cfg.Templates = opt.Templates
	}
}

func loadConfig(opt *generateOptions) (*model.Config, error) {
	cfg, err := util.ReadConfig(opt.config)
	if (err != nil && !errors.Is(err, os.ErrNotExist)) || (errors.Is(err, os.ErrNotExist) && opt.Templates == "") {
		return nil, err
	}
	return cfg, nil
}

func registerPartials(files []*model.File) {
	for _, f := range files {
		if f.Type == model.FileTypePartial {
			handlebars.RegisterPartial(f.Name, string(f.Content))
		}
	}
}

func generateForAllContexts(cfg *model.Config, files []*model.File, renderContexts []*model.RenderContext) {
	for _, ctx := range renderContexts {
		for _, f := range files {
			generateForTemplateFiles(cfg, f, ctx)
		}
	}
}

func generateForNormalFiles(cfg *model.Config, f *model.File) {
	if f.Type != model.FileTypeNormal {
		return
	}

	out := filepath.Join(cfg.Output, f.RelativePath)

	if _, err := os.Stat(out); err != nil {
		createNewFile(out, f.Content)
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

func generateForTemplateFiles(cfg *model.Config, tpl *model.File, ctx *model.RenderContext) {
	if tpl.Type != model.FileTypeTemplate {
		return
	}

	context := util.ToMap(ctx)
	newContent := handlebars.Render(tpl.Template, context)
	fileName := getFileName(tpl.Output, context)
	out := filepath.Join(cfg.Output, fileName)

	if _, err := os.Stat(out); err == nil {
		updateExistingFile(cfg, out, newContent)
	} else {
		createNewFile(out, []byte(newContent))
	}
}

func updateExistingFile(cfg *model.Config, fileName, newContent string) {
	oldContent, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	realContent := replaceBlocks(cfg, string(oldContent), newContent)
	err = util.WriteFile(fileName, []byte(realContent))
	if err != nil {
		log.Fatal(err)
	}
}

func createNewFile(fileName string, content []byte) {
	err := util.WriteFile(fileName, content)
	if err != nil {
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
