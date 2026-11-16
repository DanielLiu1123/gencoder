package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/xo/dburl"
)

// OpenConnection uses the drivers supported by gencoder. dburl's PostgreSQL
// aliases default to pgx, but gencoder uses lib/pq. Generate its DSN explicitly
// rather than assuming a dialect name identifies a driver or a compatible DSN.
// Like sql.Open, this does not establish a connection until the first query.
func OpenConnection(u *dburl.URL) (*sql.DB, error) {
	dsn, driver, err := connectionParameters(u)
	if err != nil {
		return nil, err
	}
	return sql.Open(driver, dsn)
}

func connectionParameters(u *dburl.URL) (dsn, driver string, err error) {
	if u.Dialect == "postgres" {
		return dburl.GenPq(u)
	}
	return u.DSN, u.Driver, nil
}
