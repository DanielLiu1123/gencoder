package util

import (
	"bufio"
	"context"
	"database/sql"
	"github.com/DanielLiu1123/gencoder/pkg/db"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/mailgun/raymond/v2"
	"github.com/xo/dburl"
	"gopkg.in/yaml.v3"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func ReadConfig(configPath string) (*model.Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg model.Config
	err = yaml.Unmarshal(file, &cfg)
	return &cfg, err
}

func LoadTemplates(cfg *model.Config) ([]*model.Tpl, error) {
	var templates []*model.Tpl

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

		t := &model.Tpl{
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

func CollectRenderContexts(dbConfigs ...*model.DatabaseConfig) []*model.RenderContext {
	renderContexts := make([]*model.RenderContext, 0)
	for _, dbCfg := range dbConfigs {
		contexts := collectRenderContextsForDBConfig(dbCfg)
		renderContexts = append(renderContexts, contexts...)
	}
	return renderContexts
}

func collectRenderContextsForDBConfig(dbCfg *model.DatabaseConfig) []*model.RenderContext {

	u, err := dburl.Parse(dbCfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	driver := u.Driver

	conn, err := sql.Open(driver, u.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	contexts := make([]*model.RenderContext, 0)
	for _, tbCfg := range dbCfg.Tables {
		schema := tbCfg.Schema
		if schema == "" {
			schema = dbCfg.Schema
		}

		var table *model.Table

		switch driver {
		case "mysql":
			if schema == "" {
				arr := strings.Split(u.Path, "/")
				if len(arr) > 1 {
					schema = arr[1]
				}
			}
			table, err = db.GenMySQLTable(context.Background(), conn, schema, tbCfg.Name, tbCfg.IgnoreColumns)
		case "postgres":
			if schema == "" {
				schema = "public"
			}
			table, err = db.GenPostgresTable(context.Background(), conn, schema, tbCfg.Name, tbCfg.IgnoreColumns)
		default:
			log.Fatalf("unsupported driver: %s", driver)
		}

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

// ===============================
// helper functions for Handlebars
// ===============================

// ToSnakeCase converts "userName" or "UserName" to "user_name"
//
// Example:
// ToSnakeCase("userName") => "user_name"
// ToSnakeCase("UserName") => "user_name"
func ToSnakeCase(input string) string {
	var result []rune
	for i, r := range input {
		if unicode.IsUpper(r) {
			// Add an underscore before uppercase letters, except at the start
			if i > 0 && !unicode.IsUpper(rune(input[i-1])) {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// ToCamelCase converts "user_name" to "userName"
func ToCamelCase(input string) string {
	words := strings.Split(input, "_")
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			// Convert the first character of each word to uppercase
			words[i] = strings.ToUpper(string(words[i][0])) + words[i][1:]
		}
	}
	return strings.Join(words, "")
}

// ToPascalCase converts "user_name" or "userName" to "UserName"
func ToPascalCase(input string) string {
	// Convert to CamelCase first if input is snake_case
	if strings.Contains(input, "_") {
		input = ToCamelCase(input)
	}
	// Capitalize the first character
	if len(input) > 0 {
		input = strings.ToUpper(string(input[0])) + input[1:]
	}
	return input
}
