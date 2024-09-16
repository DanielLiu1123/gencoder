package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/xo/dburl"
	"os/exec"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestGenPostgresTable(t *testing.T) {
	if !isDockerAvailable() {
		t.Log("Docker is not available. Skipping TestGenPostgresTable test.")
		return
	}

	containerID, err := startPostgresContainer()
	if err != nil {
		t.Fatalf("Failed to start Postgres container: %s", err)
	}
	defer stopContainer(containerID)

	time.Sleep(5 * time.Second)

	db, err := dburl.Open("postgres://root:root@localhost:5432/testdb?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to open database connection: %s", err)
	}

	_, err = db.Exec(`CREATE TABLE testdb.public."user" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL,
    password VARCHAR(128) NOT NULL,
    email VARCHAR(128) NOT NULL DEFAULT '',
    first_name VARCHAR(64),
    last_name VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(9) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    deleted_at TIMESTAMP,
    CONSTRAINT unique_email UNIQUE (email)
);
CREATE INDEX idx_name ON testdb.public."user" (username);
CREATE INDEX idx_status_created ON testdb.public."user" (status, created_at);
CREATE INDEX idx_full_name ON testdb.public."user" (first_name, last_name);
COMMENT ON COLUMN testdb.public."user".username IS 'Username, required';
COMMENT ON COLUMN testdb.public."user".email IS 'User email, required';
COMMENT ON COLUMN testdb.public."user".first_name IS 'First name of the user';
COMMENT ON COLUMN testdb.public."user".last_name IS 'Last name of the user';
COMMENT ON COLUMN testdb.public."user".created_at IS 'Record creation timestamp';
COMMENT ON COLUMN testdb.public."user".updated_at IS 'Record update timestamp';
COMMENT ON COLUMN testdb.public."user".status IS 'Account status';
COMMENT ON COLUMN testdb.public."user".deleted_at IS 'Record deletion timestamp';
COMMENT ON TABLE testdb.public."user" IS 'User account information';`)
	if err != nil {
		t.Fatalf("Failed to create test table: %s", err)
	}

	schema := "public"
	table := "user"

	tb, err := GenPostgresTable(context.Background(), db, schema, table, []string{"deleted_at"})
	if err != nil {
		t.Fatalf("Failed to generate Postgres table: %s", err)
	}

	assert.NotNil(t, tb)
	assert.Equal(t, "public", tb.Schema)
	assert.Equal(t, "user", tb.Name)
	assert.Equal(t, "User account information", tb.Comment)

	assert.Equal(t, 9, len(tb.Columns))
	assert.Equal(t, true, tb.Columns[0].IsPrimaryKey)
	assert.Equal(t, false, tb.Columns[1].IsPrimaryKey)

	assert.Equal(t, 5, len(tb.Indexes))
}

func startPostgresContainer() (string, error) {
	const containerID = "gencoder_test_postgres"
	cmd := exec.Command("docker", "run", "--name", containerID, "-e", "POSTGRES_USER=root", "-e", "POSTGRES_PASSWORD=root", "-e", "POSTGRES_DB=testdb", "-p", "5432:5432", "-d", "postgres:latest")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return containerID, nil
}
