package util

import (
	"context"
	"fmt"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/xo/dburl"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestReadConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.yaml")
	configContent := `
templates: "./templates"
outputMarker: "@output:"
databases:
  - name: testdb
    dsn: "mysql://user:pass@localhost:3306/testdb"
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := ReadConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, "./templates", cfg.Templates)
	assert.Equal(t, "@output:", cfg.OutputMarker)
	assert.Len(t, cfg.Databases, 1)
	assert.Equal(t, "testdb", cfg.Databases[0].Name)
}

func TestLoadTemplates(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	templateContent := `
// @gencoder.generated:{{table.name}}.go
package main

// This is a test template
`
	err = os.WriteFile(filepath.Join(tempDir, "test.go.hbs"), []byte(templateContent), 0644)
	require.NoError(t, err)

	cfg := &model.Config{
		Templates:    tempDir,
		OutputMarker: "@gencoder.generated:",
	}

	templates, err := LoadFiles(cfg, "")
	assert.NoError(t, err)
	assert.Len(t, templates, 1)
	assert.Equal(t, "test.go.hbs", templates[0].Name)
	assert.Equal(t, "{{table.name}}.go", templates[0].Output)
}

func TestCollectRenderContexts(t *testing.T) {
	// Check if Docker is available
	ctx := context.Background()

	err := exec.Command("docker", "info").Run()
	if err != nil {
		t.Skip("Docker is not available. Skipping test.")
	}

	// Start MySQL container
	mysqlContainer, err := mysql.Run(ctx,
		"mysql:latest",
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("root"),
		mysql.WithPassword("root"),
	)
	require.NoError(t, err)
	defer func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Get connection details
	host, err := mysqlContainer.Host(ctx)
	require.NoError(t, err)
	port, err := mysqlContainer.MappedPort(ctx, "3306")
	require.NoError(t, err)

	// Construct DSN
	dsn := fmt.Sprintf("mysql://root:root@%s:%s/testdb", host, port.Port())

	// Create test table
	db, err := dburl.Open(dsn)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE testdb.user ( 
        id INT AUTO_INCREMENT PRIMARY KEY, 
        username VARCHAR(64) NOT NULL COMMENT 'Username, required', 
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation timestamp' 
    ) COMMENT='User account information';
	`)
	require.NoError(t, err)

	// Prepare config
	cfg := &model.Config{
		Databases: []*model.DatabaseConfig{
			{
				Name:   "testdb",
				Dsn:    dsn,
				Schema: "testdb",
				Tables: []*model.TableConfig{
					{
						Name: "user",
					},
				},
			},
		},
		Properties: map[string]string{
			"k1": "v1",
		},
	}

	// Run CollectRenderContexts
	contexts := CollectRenderContexts(cfg, map[string]string{"k1": "v01", "k2": "v2"})

	// Assertions
	assert.Len(t, contexts, 1)
	assert.Equal(t, "user", contexts[0].Table.Name)
	assert.Equal(t, "testdb", contexts[0].Table.Schema)
	assert.Len(t, contexts[0].Table.Columns, 3)

	// Check column details
	columnNames := []string{"id", "username", "created_at"}
	for i, col := range contexts[0].Table.Columns {
		assert.Equal(t, columnNames[i], col.Name)
	}

	// Check properties
	assert.Len(t, contexts[0].Properties, 2)
	assert.Equal(t, "v01", contexts[0].Properties["k1"])
	assert.Equal(t, "v2", contexts[0].Properties["k2"])
}
