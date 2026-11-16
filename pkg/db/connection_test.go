package db

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/xo/dburl"
)

func TestPostgresConnectionAliases(t *testing.T) {
	for _, scheme := range []string{"postgres", "postgresql", "pg", "pgsql", "pgx", "pq"} {
		t.Run(scheme, func(t *testing.T) {
			u, err := dburl.Parse(scheme + "://user:password@localhost:5432/dbname?sslmode=disable")
			require.NoError(t, err)
			require.Equal(t, "postgres", u.Dialect)
			dsn, driver, err := connectionParameters(u)
			require.NoError(t, err)
			require.Equal(t, "postgres", driver)
			require.Contains(t, dsn, "sslmode=disable")
			// Parsing and opening are deliberately offline: never Ping or query.
			_, err = pq.NewConnector(dsn)
			require.NoError(t, err)
			conn, err := OpenConnection(u)
			require.NoError(t, err)
			t.Cleanup(func() { _ = conn.Close() })
			require.IsType(t, &pq.Driver{}, conn.Driver())
		})
	}
}

func TestPostgresConnectionPreservesCredentialsAndOptions(t *testing.T) {
	u, err := dburl.Parse("postgres://user%20name:p%27a%5Css@localhost:5433/mydb?sslmode=verify-full&application_name=gen%20test&search_path=custom&connect_timeout=7&binary_parameters=yes&options=-c%20statement_timeout%3D1000")
	require.NoError(t, err)
	dsn, driver, err := connectionParameters(u)
	require.NoError(t, err)
	require.Equal(t, "postgres", driver)
	for _, option := range []string{"user='user name'", "password='p\\'a\\\\ss'", "host=localhost", "port=5433", "dbname=mydb", "sslmode=verify-full", "application_name='gen test'", "search_path=custom", "connect_timeout=7", "binary_parameters=yes", "options='-c statement_timeout=1000'"} {
		require.Contains(t, dsn, option)
	}
	_, err = pq.NewConnector(dsn)
	require.NoError(t, err)
}

func TestOtherConnectionParametersUnchanged(t *testing.T) {
	for _, raw := range []string{"mysql://user:password@localhost:3306/dbname?parseTime=true", "sqlserver://user:password@localhost:1433/dbname?encrypt=disable"} {
		u, err := dburl.Parse(raw)
		require.NoError(t, err)
		dsn, driver, err := connectionParameters(u)
		require.NoError(t, err)
		require.Equal(t, u.Driver, driver)
		require.Equal(t, u.DSN, dsn)
	}
}
