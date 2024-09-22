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
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return &cfg, err
	}
	return &cfg, nil
}

// LoadTemplates loads all templates from the given configuration and command line templates
func LoadTemplates(cfg *model.Config, commandLineTemplates string) ([]*model.Tpl, error) {

	template := cfg.GetTemplates()
	if commandLineTemplates != "" {
		template = commandLineTemplates
	}

	var templates []*model.Tpl

	// If url is provided, download templates
	isGitUrl, err := regexp.Match(`^.*(github\.com)/.*`, []byte(template))
	if err != nil {
		return nil, err
	}

	if isGitUrl {
		// Parse GitHub URL to extract the repo and branch (if provided)
		re := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)(/tree/([^/]+)/(.*))?`)
		matches := re.FindStringSubmatch(template)
		if matches == nil {
			return nil, fmt.Errorf("invalid GitHub URL format")
		}

		owner := matches[1]
		repo := matches[2]
		branch := "main"
		if matches[4] != "" {
			branch = matches[4]
		}
		dirInRepo := matches[5]

		tmpDir, err := os.MkdirTemp("", "templates")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		cloneUrl := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
		cmd := exec.Command("git", "clone", "--branch", branch, "--depth", "1", cloneUrl, tmpDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to clone repository: %v", err)
		}

		// If a directory within the repo is specified, update the path
		if dirInRepo != "" {
			template = filepath.Join(tmpDir, dirInRepo)
		} else {
			template = tmpDir
		}
	}

	err = filepath.WalkDir(template, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		template := handlebars.Compile(string(b))

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

// CollectRenderContexts collects render contexts for the given database configurations
func CollectRenderContexts(cfg *model.Config, commandLineProperties map[string]string) []*model.RenderContext {
	renderContexts := make([]*model.RenderContext, 0)
	for _, dbCfg := range cfg.Databases {
		contexts := collectRenderContextsForDBConfig(cfg, dbCfg)
		renderContexts = append(renderContexts, contexts...)
	}

	// Use command line properties to override properties in config file
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

	var mu sync.Mutex
	contexts := make([]*model.RenderContext, 0)
	var wg sync.WaitGroup

	for _, tbCfg := range dbCfg.Tables {
		wg.Add(1)
		go func(tbCfg *model.TableConfig) {
			defer wg.Done()

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

			// table not found
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
