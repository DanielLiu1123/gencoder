package info

import (
	"context"
	"database/sql"
	"strings"
)

// GenMySQLTable generates a MySQL table.
func GenMySQLTable(ctx context.Context, db *sql.DB, schema, table string) (*Table, error) {

	// table info
	const tableSql = `
select table_schema  as table_schema,
       table_name    as name,
       table_comment as comment
from information_schema.tables
where table_schema = ?
  and table_name = ?;
`
	tableRow := db.QueryRowContext(ctx, tableSql, schema, table)
	var t Table
	if err := tableRow.Scan(&t.Schema, &t.TableName, &t.Comment); err != nil {
		return nil, err
	}

	// columns info
	const columnsSql = `
select ordinal_position                     as ordinal,
       column_name                          as name,
       column_type                          as type,
       if(is_nullable = 'YES', true, false) as is_nullable,
       column_default                       as default_value,
       if(column_key = 'PRI', true, false)  as is_primary_key,
       column_comment                       as comment
from information_schema.columns t
where table_schema = ?
  and table_name = ?
order by ordinal;
`
	columnRow, err := db.QueryContext(ctx, columnsSql, schema, table)
	if err != nil {
		return nil, err
	}
	var columns []*Column
	for columnRow.Next() {
		var c Column
		if err := columnRow.Scan(&c.Ordinal, &c.ColumnName, &c.ColumnType, &c.IsNullable, &c.DefaultValue, &c.IsPrimaryKey, &c.Comment); err != nil {
			return nil, err
		}
		columns = append(columns, &c)
	}

	// indexes info
	const indexesSql = `
select table_schema                                    as 'scheme',
       table_name                                      as 'table',
       index_name                                      as name,
       if(non_unique = 0, true, false)                 as is_unique,
       group_concat(column_name order by seq_in_index) as columns
from information_schema.statistics
where table_schema = ?
  and table_name = ?
group by table_schema, table_name, index_name, non_unique;
`
	indexRow, err := db.QueryContext(ctx, indexesSql, schema, table)
	if err != nil {
		return nil, err
	}
	var indexes []*Index
	for indexRow.Next() {
		var idx Index
		var columns string
		if err := indexRow.Scan(&idx.Schema, &idx.TableName, &idx.IndexName, &idx.IsUnique, &columns); err != nil {
			return nil, err
		}
		idx.Columns = strings.Split(columns, ",")
		indexes = append(indexes, &idx)
	}

	t.Columns = columns
	t.Indexes = indexes
	return &t, nil
}
