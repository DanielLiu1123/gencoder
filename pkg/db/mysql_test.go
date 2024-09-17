package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/xo/dburl"
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
	cmd := exec.Command("docker", "run", "--name", containerID, "-e", "MYSQL_ROOT_PASSWORD=root", "-e", "MYSQL_DATABASE=testdb", "-p", "3306:3306", "-d", "mysql:latest")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return containerID, nil
}

// Helper function to stop and remove the MySQL container
func stopContainer(containerID string) {
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

	defer stopContainer(containerID)

	// Wait for MySQL to initialize
	time.Sleep(10 * time.Second)

	db, err := dburl.Open("mysql://root:root@localhost:3306/testdb")
	if err != nil {
		t.Fatalf("Failed to open database connection: %s", err)
	}

	_, err = db.Exec(`CREATE TABLE testdb.user (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(64) NOT NULL COMMENT 'Username, required',
        password VARCHAR(128) NOT NULL,
        email VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'User email, required',
        first_name VARCHAR(64) COMMENT 'First name of the user',
        last_name VARCHAR(64) COMMENT 'Last name of the user',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation timestamp',
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record update timestamp',
        status ENUM('active', 'inactive', 'suspended') DEFAULT 'active' COMMENT 'Account status',
		deleted_at TIMESTAMP COMMENT 'Record deletion timestamp',
        INDEX idx_name (username),
        UNIQUE INDEX idx_email (email),
        INDEX idx_status_created (status, created_at),
        INDEX idx_full_name (first_name, last_name)
    ) COMMENT='User account information';`)
	if err != nil {
		t.Fatalf("Failed to create test table: %s", err)
	}

	schema := "testdb"
	table := "user"

	tb, err := GenMySQLTable(context.Background(), db, schema, table, []string{"deleted_at"})
	if err != nil {
		t.Fatalf("Failed to generate MySQL table: %s", err)
	}

	assert.NotNil(t, tb)
	assert.Equal(t, "testdb", tb.Schema)
	assert.Equal(t, "user", tb.Name)
	assert.Equal(t, "User account information", tb.Comment)

	assert.Equal(t, 9, len(tb.Columns))
	assert.Equal(t, true, tb.Columns[0].IsPrimaryKey)
	assert.Equal(t, false, tb.Columns[1].IsPrimaryKey)

	assert.Equal(t, 5, len(tb.Indexes))
}
