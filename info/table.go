package info

import (
	"context"
	"database/sql"
)

// GenMySQLTable generates a MySQL table and fills the Table structure.
func GenMySQLTable(ctx context.Context, db *sql.DB, schema, table string) (*Table, error) {
	// Table info
	const tableSql = `
select table_schema,
       table_name,
       table_comment
from information_schema.tables
where table_schema = ?
  and table_name = ?;
`
	tableRow := db.QueryRowContext(ctx, tableSql, schema, table)
	var t Table
	if err := tableRow.Scan(&t.Schema, &t.Name, &t.Comment); err != nil {
		return nil, err
	}

	// Columns info
	const columnsSql = `
select ordinal_position,
       column_name,
       column_type,
       is_nullable = 'YES' as is_nullable,
       column_default,
       column_key = 'PRI' as is_primary_key,
       column_comment
from information_schema.columns
where table_schema = ?
  and table_name = ?
order by ordinal_position;
`
	columnRows, err := db.QueryContext(ctx, columnsSql, schema, table)
	if err != nil {
		return nil, err
	}
	defer columnRows.Close()

	var columns []*Column
	for columnRows.Next() {
		var col Column
		if err := columnRows.Scan(&col.Ordinal, &col.Name, &col.Type, &col.IsNullable, &col.DefaultValue, &col.IsPrimaryKey, &col.Comment); err != nil {
			return nil, err
		}
		columns = append(columns, &col)
	}
	t.Columns = columns

	// Indexes and IndexColumns info
	const indexesSql = `
select index_name,
       non_unique = 0 as is_unique,
       index_name = 'PRIMARY' as is_primary,
       seq_in_index as ordinal,
       column_name
from information_schema.statistics
where table_schema = ?
  and table_name = ?
order by index_name, seq_in_index;
`
	indexRows, err := db.QueryContext(ctx, indexesSql, schema, table)
	if err != nil {
		return nil, err
	}
	defer indexRows.Close()

	indexMap := make(map[string]*Index)
	for indexRows.Next() {
		var indexName, columnName string
		var isUnique, isPrimary bool
		var ordinal int

		if err := indexRows.Scan(&indexName, &isUnique, &isPrimary, &ordinal, &columnName); err != nil {
			return nil, err
		}

		// 如果此 index 还未在 map 中记录，则添加它
		if _, exists := indexMap[indexName]; !exists {
			indexMap[indexName] = &Index{
				Name:         indexName,
				IsUnique:     isUnique,
				IsPrimary:    isPrimary,
				IndexColumns: []*IndexColumn{},
			}
		}

		// 添加 IndexColumn 到对应的 Index 中
		indexMap[indexName].IndexColumns = append(indexMap[indexName].IndexColumns, &IndexColumn{
			Ordinal:    ordinal,
			ColumnName: columnName,
		})
	}

	// 从 map 中提取所有索引
	var indexes []*Index
	for _, index := range indexMap {
		indexes = append(indexes, index)
	}
	t.Indexes = indexes

	return &t, nil
}
