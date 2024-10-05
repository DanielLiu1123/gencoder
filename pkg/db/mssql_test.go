package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mssql"
	"github.com/xo/dburl"
	"os/exec"
	"testing"

	_ "github.com/microsoft/go-mssqldb"
)

func TestGenMssqlTable(t *testing.T) {
	err := exec.Command("docker", "info").Run()
	if err != nil {
		t.Skip("Docker not available, skipping MSSQL tests")
	}

	ctx := context.Background()
	mssqlContainer, err := mssql.Run(ctx,
		"mcr.microsoft.com/mssql/server:2022-latest",
		mssql.WithAcceptEULA(),
		mssql.WithPassword("Sa123456.."),
	)
	require.NoError(t, err)
	defer func() {
		if err := mssqlContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	host, err := mssqlContainer.Host(ctx)
	require.NoError(t, err)
	port, err := mssqlContainer.MappedPort(ctx, "1433")
	require.NoError(t, err)

	dsn := fmt.Sprintf("mssql://sa:Sa123456..@%s:%s/master", host, port.Port())
	db, err := dburl.Open(dsn)
	require.NoError(t, err)

	// Create table and indexes in MSSQL
	_, err = db.Exec(`CREATE TABLE master.dbo.[user] (
		id         INT IDENTITY (1,1) PRIMARY KEY,
		username   NVARCHAR(64)  NOT NULL,
		password   NVARCHAR(128) NOT NULL,
		email      NVARCHAR(128) NOT NULL DEFAULT '',
		first_name NVARCHAR(64),
		last_name  NVARCHAR(64),
		created_at DATETIME               DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME               DEFAULT CURRENT_TIMESTAMP,
		status     NVARCHAR(9)            DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
		deleted_at DATETIME      NULL,
		CONSTRAINT unique_email UNIQUE (email)
	);
	CREATE INDEX idx_name ON master.dbo.[user] (username);
	CREATE INDEX idx_status_created ON master.dbo.[user] (status, created_at);
	CREATE INDEX idx_full_name ON master.dbo.[user] (first_name, last_name);
	
	-- Add comments
	EXEC sp_addextendedproperty
		 @name = N'MS_Description',
		 @value = N'User account information',
		 @level0type = N'SCHEMA', @level0name = 'dbo',
		 @level1type = N'TABLE', @level1name = 'user';
	EXEC sp_addextendedproperty
		 @name = N'MS_Description',
		 @value = N'Username, required',
		 @level0type = N'SCHEMA', @level0name = 'dbo',
		 @level1type = N'TABLE', @level1name = 'user',
		 @level2type = N'COLUMN', @level2name = 'username';
	EXEC sp_addextendedproperty
		 @name = N'MS_Description',
		 @value = N'User email, required',
		 @level0type = N'SCHEMA', @level0name = 'dbo',
		 @level1type = N'TABLE', @level1name = 'user',
		 @level2type = N'COLUMN', @level2name = 'email';`)
	require.NoError(t, err)

	// Call the GenMssqlTable function to generate the table model
	schema := "dbo"
	table := "user"

	tb, err := GenMssqlTable(ctx, db, schema, table, []string{"deleted_at"})
	require.NoError(t, err)

	// Assertions on the table structure
	assert.NotNil(t, tb)
	assert.Equal(t, "dbo", tb.Schema)
	assert.Equal(t, "user", tb.Name)
	assert.Equal(t, "User account information", tb.Comment)

	// Assertions on the columns
	assert.Equal(t, 9, len(tb.Columns)) // Columns without deleted_at
	assert.Equal(t, true, tb.Columns[0].IsPrimaryKey)
	assert.Equal(t, false, tb.Columns[1].IsPrimaryKey)

	// Assertions on the indexes
	assert.Equal(t, 5, len(tb.Indexes))
	assert.Contains(t, tb.Indexes[0].Name, "PK__user")  // Primary key
	assert.Equal(t, "unique_email", tb.Indexes[1].Name) // Unique index
	assert.Equal(t, "idx_full_name", tb.Indexes[2].Name)
	assert.Equal(t, "first_name", tb.Indexes[2].Columns[0].Name)
	assert.Equal(t, "last_name", tb.Indexes[2].Columns[1].Name)
	assert.Equal(t, "idx_name", tb.Indexes[3].Name)
	assert.Equal(t, "idx_status_created", tb.Indexes[4].Name)
	assert.Equal(t, "status", tb.Indexes[4].Columns[0].Name)
	assert.Equal(t, "created_at", tb.Indexes[4].Columns[1].Name)
}
