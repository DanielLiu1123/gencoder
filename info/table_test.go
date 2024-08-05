package info

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Helper function to check if Docker is available
func isDockerAvailable() bool {
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	return err == nil
}

// Helper function to start a MySQL container
func startMySQLContainer() (string, error) {
	const containerID = "gencoder_test_mysql"
	cmd := exec.Command("docker", "run", "--name", containerID, "-e", "MYSQL_ROOT_PASSWORD=123456", "-e", "MYSQL_DATABASE=testdb", "-p", "3306:3306", "-d", "mysql:latest")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return containerID, nil
}

// Helper function to stop and remove the MySQL container
func stopAndRemoveMySQLContainer(containerID string) {
	err := exec.Command("docker", "stop", containerID).Run()
	if err != nil {
		panic(err)
	}
	err = exec.Command("docker", "rm", containerID).Run()
	if err != nil {
		panic(err)
	}
}

func TestGenMySQLTable(t *testing.T) {
	if !isDockerAvailable() {
		t.Log("Docker is not available. Skipping test.")
		return
	}

	containerID, err := startMySQLContainer()
	if err != nil {
		t.Fatalf("Failed to start MySQL container: %s", err)
	}

	defer stopAndRemoveMySQLContainer(containerID)

	// Wait for MySQL to initialize
	time.Sleep(10 * time.Second)

	dsn := "root:123456@tcp(127.0.0.1:3306)/testdb"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %s", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			t.Fatalf("Failed to close the database: %s", err)
		}
	}(db)

	_, err = db.Exec(`CREATE TABLE user (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(64) NOT NULL,
		email VARCHAR(128),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_name (name),
		UNIQUE INDEX idx_email (email)
	)`)
	if err != nil {
		t.Fatalf("Failed to create test table: %s", err)
	}

	schema := "testdb"
	table := "user"

	tb, err := GenMySQLTable(context.Background(), db, schema, table)
	if err != nil {
		t.Fatalf("Failed to generate MySQL table: %s", err)
	}

	assert.NotNil(t, tb)
	assert.Equal(t, "testdb", tb.Schema)
	assert.Equal(t, "user", tb.TableName)
	assert.Equal(t, "", tb.Comment)

	assert.Equal(t, 4, len(tb.Columns))
	assert.Equal(t, true, tb.Columns[0].IsPrimaryKey)
	assert.Equal(t, false, tb.Columns[1].IsPrimaryKey)

	assert.Equal(t, 3, len(tb.Indexes))
	assert.Contains(t, tb.Indexes, &Index{Schema: schema,
		TableName: table, IndexName: "PRIMARY", IsUnique: true,
		Columns: []string{"id"}})
	assert.Contains(t, tb.Indexes, &Index{Schema: schema,
		TableName: table, IndexName: "idx_name", IsUnique: false,
		Columns: []string{"name"}})
	assert.Contains(t, tb.Indexes, &Index{Schema: schema,
		TableName: table, IndexName: "idx_email", IsUnique: true,
		Columns: []string{"email"}})
}
