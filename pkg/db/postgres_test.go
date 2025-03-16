package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/xo/dburl"
	"os/exec"
	"testing"

	_ "github.com/lib/pq"
)

func TestGenPostgresTable(t *testing.T) {
	err := exec.Command("docker", "info").Run()
	if err != nil {
		t.Skip("Docker not available, skipping MySQL tests")
	}

	ctx := context.Background()
	postgresContainer, err := postgres.Run(ctx,
		"postgres:latest",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("root"),
		postgres.WithPassword("root"),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)
	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://root:root@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err := dburl.Open(dsn)
	require.NoError(t, err)

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
	require.NoError(t, err)

	schema := "public"
	table := "user"

	tb, err := GenPostgresTable(context.Background(), db, schema, table, []string{"deleted_at"})
	require.NoError(t, err)

	assert.NotNil(t, tb)
	assert.Equal(t, "public", tb.Schema)
	assert.Equal(t, "user", tb.Name)
	assert.Equal(t, "User account information", tb.Comment)

	assert.Equal(t, 9, len(tb.Columns))
	assert.Equal(t, true, tb.Columns[0].IsPrimaryKey)
	assert.Equal(t, false, tb.Columns[0].IsNullable)
	assert.Equal(t, false, tb.Columns[1].IsPrimaryKey)
	assert.Equal(t, "first_name", tb.Columns[4].Name)
	assert.Equal(t, true, tb.Columns[4].IsNullable)

	assert.Equal(t, 5, len(tb.Indexes))
	assert.Equal(t, "user_pkey", tb.Indexes[0].Name)    // Primary key
	assert.Equal(t, "unique_email", tb.Indexes[1].Name) // Unique index
	assert.Equal(t, "idx_full_name", tb.Indexes[2].Name)
	assert.Equal(t, "first_name", tb.Indexes[2].Columns[0].Name)
	assert.Equal(t, "last_name", tb.Indexes[2].Columns[1].Name)
	assert.Equal(t, "idx_name", tb.Indexes[3].Name)
	assert.Equal(t, "idx_status_created", tb.Indexes[4].Name)
	assert.Equal(t, "status", tb.Indexes[4].Columns[0].Name)
	assert.Equal(t, "created_at", tb.Indexes[4].Columns[1].Name)
}
