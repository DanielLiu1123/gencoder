package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"slices"
	"sort"
)

// GenMySQLTable generates a MySQL table and fills the Table structure.
func GenMySQLTable(ctx context.Context, db *sql.DB, schema, table string, ignoreColumns []string) (*model.Table, error) {
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
	var t model.Table
	if err := tableRow.Scan(&t.Schema, &t.Name, &t.Comment); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
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

	columns := make([]*model.Column, 0)
	for columnRows.Next() {
		var col model.Column
		if err := columnRows.Scan(&col.Ordinal, &col.Name, &col.Type, &col.IsNullable, &col.DefaultValue, &col.IsPrimaryKey, &col.Comment); err != nil {
			return nil, err
		}

		if len(ignoreColumns) > 0 && slices.Contains(ignoreColumns, col.Name) {
			continue
		}

		columns = append(columns, &col)
	}
	t.Columns = columns

	// Indexes and Columns info
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

	indexMap := make(map[string]*model.Index)
	for indexRows.Next() {
		var indexName, columnName string
		var isUnique, isPrimary bool
		var ordinal int

		if err := indexRows.Scan(&indexName, &isUnique, &isPrimary, &ordinal, &columnName); err != nil {
			return nil, err
		}

		if _, exists := indexMap[indexName]; !exists {
			indexMap[indexName] = &model.Index{
				Name:      indexName,
				IsUnique:  isUnique,
				IsPrimary: isPrimary,
				Columns:   []*model.IndexColumn{},
			}
		}

		indexMap[indexName].Columns = append(indexMap[indexName].Columns, &model.IndexColumn{
			Ordinal: ordinal,
			Name:    columnName,
		})
	}

	indexes := make([]*model.Index, 0)
	for _, index := range indexMap {
		indexes = append(indexes, index)
	}

	sort.Slice(indexes, func(i, j int) bool {
		if indexes[i].IsPrimary != indexes[j].IsPrimary {
			return indexes[i].IsPrimary
		}
		if indexes[i].IsUnique != indexes[j].IsUnique {
			return indexes[i].IsUnique
		}
		return indexes[i].Name < indexes[j].Name
	})

	t.Indexes = indexes

	return &t, nil
}
