package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/xo/dburl"
	"os/exec"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestGenMySQLTable(t *testing.T) {

	err := exec.Command("docker", "info").Run()
	if err != nil {
		t.Skip("Docker not available, skipping MySQL tests")
	}

	ctx := context.Background()
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

	host, err := mysqlContainer.Host(ctx)
	require.NoError(t, err)
	port, err := mysqlContainer.MappedPort(ctx, "3306")
	require.NoError(t, err)

	dsn := fmt.Sprintf("mysql://root:root@%s:%s/testdb", host, port.Port())
	db, err := dburl.Open(dsn)
	require.NoError(t, err)

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
	require.NoError(t, err)

	schema := "testdb"
	table := "user"

	tb, err := GenMySQLTable(context.Background(), db, schema, table, []string{"deleted_at"})
	require.NoError(t, err)

	assert.NotNil(t, tb)
	assert.Equal(t, "testdb", tb.Schema)
	assert.Equal(t, "user", tb.Name)
	assert.Equal(t, "User account information", tb.Comment)

	assert.Equal(t, 9, len(tb.Columns))
	assert.Equal(t, true, tb.Columns[0].IsPrimaryKey)
	assert.Equal(t, false, tb.Columns[1].IsPrimaryKey)

	assert.Equal(t, 5, len(tb.Indexes))
}
