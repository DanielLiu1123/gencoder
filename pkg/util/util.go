package util

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/DanielLiu1123/gencoder/pkg/db"
	"github.com/DanielLiu1123/gencoder/pkg/handlebars"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/xo/dburl"
	"gopkg.in/yaml.v3"
)

// ReadConfig reads the configuration file from the given path
func ReadConfig(configPath string) (*model.Config, error) {
	var cfg model.Config
	file, err := os.ReadFile(configPath)
	if err != nil {
		return &cfg, err
	}

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}

// LoadFiles loads all files from the given path
func LoadFiles(cfg *model.Config) ([]*model.File, error) {
	templatePath := cfg.GetTemplates()

	if isGitHubURL(templatePath) {
		var err error
		templatePath, err = cloneGitHubRepo(templatePath)
		if err != nil {
			return nil, err
		}
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				log.Printf("failed to remove temp dir: %s", err)
			}
		}(templatePath)
	}

	return loadFilesFromPath(templatePath, cfg)
}

func isGitHubURL(url string) bool {
	matched, _ := regexp.Match(`^.*(github\.com)/.*`, []byte(url))
	return matched
}

func cloneGitHubRepo(url string) (string, error) {
	re := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)(/tree/([^/]+)/(.*))?`)
	matches := re.FindStringSubmatch(url)
	if matches == nil {
		return "", fmt.Errorf("invalid GitHub URL format")
	}

	owner, repo := matches[1], matches[2]
	branch := "main"
	if matches[4] != "" {
		branch = matches[4]
	}
	dirInRepo := matches[5]

	tmpDir, err := os.MkdirTemp("", "templates")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %v", err)
	}

	cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
	cmd := exec.Command("git", "clone", "--branch", branch, "--depth", "1", cloneURL, tmpDir)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to clone repository: %v", err)
	}

	if dirInRepo != "" {
		return filepath.Join(tmpDir, dirInRepo), nil
	}
	return tmpDir, nil
}

func loadFilesFromPath(rootPath string, cfg *model.Config) ([]*model.File, error) {

	var templates []*model.File

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var f model.File

		f.Name = d.Name()
		rel, err := filepath.Rel(rootPath, path)
		if err != nil {
			return err
		}
		f.RelativePath = rel
		f.Content = b
		f.Type = model.FileTypeNormal

		// template or partial
		if strings.HasSuffix(d.Name(), ".hbs") || strings.HasSuffix(d.Name(), ".mustache") {
			content := string(b)
			output := getFileNameTemplate(content, cfg)
			if output != "" {
				f.Type = model.FileTypeTemplate
				f.Output = output
			} else {
				f.Type = model.FileTypePartial
			}
			f.Template = handlebars.Compile(content)
		}

		templates = append(templates, &f)
		return nil
	})

	return templates, err
}

// CollectRenderContexts collects render contexts for the given database configurations
func CollectRenderContexts(cfg *model.Config, commandLineProperties map[string]string) []*model.RenderContext {
	var renderContexts []*model.RenderContext
	for _, dbCfg := range cfg.Databases {
		contexts := collectRenderContextsForDBConfig(cfg, dbCfg)
		renderContexts = append(renderContexts, contexts...)
	}

	for _, rc := range renderContexts {
		for k, v := range commandLineProperties {
			rc.Properties[k] = v
		}
	}

	return renderContexts
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

func collectRenderContextsForDBConfig(cfg *model.Config, dbCfg *model.DatabaseConfig) []*model.RenderContext {
	u, err := dburl.Parse(dbCfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(u.Driver, u.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	var mu sync.Mutex
	var contexts []*model.RenderContext
	var wg sync.WaitGroup

	for _, tbCfg := range dbCfg.Tables {
		wg.Add(1)
		go func(tbCfg *model.TableConfig) {
			defer wg.Done()

			schema := getSchema(tbCfg, dbCfg, u)
			table, err := generateTable(conn, u.Driver, schema, tbCfg)
			if err != nil {
				log.Fatal(err)
			}

			if table == nil {
				log.Printf("table %s.%s not found, skipping", schema, tbCfg.Name)
				return
			}

			ctx := createRenderContext(cfg, dbCfg, tbCfg, table)

			mu.Lock()
			contexts = append(contexts, ctx)
			mu.Unlock()
		}(tbCfg)
	}

	wg.Wait()

	return contexts
}

func getSchema(tbCfg *model.TableConfig, dbCfg *model.DatabaseConfig, u *dburl.URL) string {
	if tbCfg.Schema != "" {
		return tbCfg.Schema
	}
	if dbCfg.Schema != "" {
		return dbCfg.Schema
	}
	if u.Driver == "mysql" {
		arr := strings.Split(u.Path, "/")
		if len(arr) > 1 {
			return arr[1] // database name as schema in MySQL
		}
	}
	if u.Driver == "postgres" {
		return "public" // default schema in PostgreSQL
	}
	if u.Driver == "sqlserver" {
		return "dbo" // default schema in SQL Server
	}
	return ""
}

func generateTable(conn *sql.DB, driver, schema string, tbCfg *model.TableConfig) (*model.Table, error) {
	switch driver {
	case "mysql":
		return db.GenMySQLTable(context.Background(), conn, schema, tbCfg.Name, tbCfg.IgnoreColumns)
	case "postgres":
		return db.GenPostgresTable(context.Background(), conn, schema, tbCfg.Name, tbCfg.IgnoreColumns)
	case "sqlserver":
		return db.GenMssqlTable(context.Background(), conn, schema, tbCfg.Name, tbCfg.IgnoreColumns)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}
}

func createRenderContext(cfg *model.Config, dbCfg *model.DatabaseConfig, tbCfg *model.TableConfig, table *model.Table) *model.RenderContext {
	properties := make(map[string]string)
	for k, v := range cfg.Properties {
		properties[k] = v
	}
	for k, v := range dbCfg.Properties {
		properties[k] = v
	}
	for k, v := range tbCfg.Properties {
		properties[k] = v
	}

	return &model.RenderContext{
		Table:          table,
		Properties:     properties,
		Config:         cfg,
		DatabaseConfig: dbCfg,
		TableConfig:    tbCfg,
	}
}

// WriteFile writes the content to the given file, creating directories if necessary
func WriteFile(filename string, content []byte) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	err := os.WriteFile(filename, content, 0644)
	if err != nil {
		return err
	}
	return nil
}
