package db

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"sort"

	"github.com/DanielLiu1123/gencoder/pkg/model"
)

// GenMySQLTable generates a MySQL table and fills the Table structure.
func GenMySQLTable(ctx context.Context, db *sql.DB, schema, table string, ignoreColumns []string) (*model.Table, error) {
	t, err := getMySQLTableInfo(ctx, db, schema, table)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	columns, err := getMySQLColumnsInfo(ctx, db, schema, table, ignoreColumns)
	if err != nil {
		return nil, err
	}
	t.Columns = columns

	indexes, err := getMySQLIndexesInfo(ctx, db, schema, table)
	if err != nil {
		return nil, err
	}
	t.Indexes = indexes

	return t, nil
}

func getMySQLTableInfo(ctx context.Context, db *sql.DB, schema, table string) (*model.Table, error) {
	const tableSql = `
		select table_schema, table_name, table_comment
		from information_schema.tables
		where table_schema = ? and table_name = ?;
	`
	var t model.Table
	err := db.QueryRowContext(ctx, tableSql, schema, table).Scan(&t.Schema, &t.Name, &t.Comment)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func getMySQLColumnsInfo(ctx context.Context, db *sql.DB, schema, table string, ignoreColumns []string) ([]*model.Column, error) {
	const columnsSql = `
		select ordinal_position, column_name, column_type,
			   is_nullable = 'YES' as is_nullable,
			   column_default, column_key = 'PRI' as is_primary_key,
			   column_comment
		from information_schema.columns
		where table_schema = ? and table_name = ?
		order by ordinal_position;
	`
	rows, err := db.QueryContext(ctx, columnsSql, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []*model.Column
	for rows.Next() {
		var col model.Column
		if err := rows.Scan(&col.Ordinal, &col.Name, &col.Type, &col.IsNullable, &col.DefaultValue, &col.IsPrimaryKey, &col.Comment); err != nil {
			return nil, err
		}
		if !slices.Contains(ignoreColumns, col.Name) {
			columns = append(columns, &col)
		}
	}
	return columns, nil
}

func getMySQLIndexesInfo(ctx context.Context, db *sql.DB, schema, table string) ([]*model.Index, error) {
	const indexesSql = `
		select index_name, non_unique = 0 as is_unique,
			   index_name = 'PRIMARY' as is_primary,
			   seq_in_index as ordinal, column_name
		from information_schema.statistics
		where table_schema = ? and table_name = ?
		order by index_name, seq_in_index;
	`
	rows, err := db.QueryContext(ctx, indexesSql, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indexMap := make(map[string]*model.Index)
	for rows.Next() {
		var indexName, columnName string
		var isUnique, isPrimary bool
		var ordinal int

		if err := rows.Scan(&indexName, &isUnique, &isPrimary, &ordinal, &columnName); err != nil {
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

	indexes := make([]*model.Index, 0, len(indexMap))
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

	return indexes, nil
}
