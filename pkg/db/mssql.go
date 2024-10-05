package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"slices"
	"sort"
)

// GenMssqlTable generates an MSSQL table and fills the Table structure.
func GenMssqlTable(ctx context.Context, db *sql.DB, schema string, name string, ignoreColumns []string) (*model.Table, error) {
	t, err := getMssqlTableInfo(ctx, db, schema, name)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	columns, err := getMssqlColumnsInfo(ctx, db, schema, name, ignoreColumns)
	if err != nil {
		return nil, err
	}
	t.Columns = columns

	indexes, err := getMssqlIndexesInfo(ctx, db, schema, name)
	if err != nil {
		return nil, err
	}
	t.Indexes = indexes

	return t, nil
}

func getMssqlTableInfo(ctx context.Context, db *sql.DB, schema, name string) (*model.Table, error) {
	const tableSql = `
		SELECT s.name AS table_schema,
		       t.name AS table_name,
		       p.value AS table_comment
		FROM sys.tables t
		JOIN sys.schemas s ON t.schema_id = s.schema_id
		LEFT JOIN sys.extended_properties p ON p.major_id = t.object_id AND p.minor_id = 0 AND p.name = 'MS_Description'
		WHERE s.name = @p1
		  AND t.name = @p2;
	`
	var t model.Table
	err := db.QueryRowContext(ctx, tableSql, schema, name).Scan(&t.Schema, &t.Name, &t.Comment)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func getMssqlColumnsInfo(ctx context.Context, db *sql.DB, schema, name string, ignoreColumns []string) ([]*model.Column, error) {
	const columnsSql = `
SELECT 
    c.column_id AS ordinal,
    c.name AS column_name,
    tp.name AS data_type,
    c.is_nullable AS is_nullable,
    dc.definition AS default_value,
    CASE 
        WHEN EXISTS (
            SELECT 1
            FROM sys.index_columns ic
            JOIN sys.indexes i ON ic.index_id = i.index_id
            WHERE ic.object_id = c.object_id AND ic.column_id = c.column_id AND i.is_primary_key = 1
        ) 
        THEN 1 ELSE 0 
    END AS is_primary,
    ep.value AS comment
FROM sys.columns c
JOIN sys.types tp ON c.user_type_id = tp.user_type_id
JOIN sys.tables t ON c.object_id = t.object_id
JOIN sys.schemas s ON t.schema_id = s.schema_id
LEFT JOIN sys.default_constraints dc ON c.default_object_id = dc.object_id
LEFT JOIN sys.extended_properties ep ON ep.major_id = c.object_id AND ep.minor_id = c.column_id AND ep.name = 'MS_Description'
WHERE s.name = @p1
  AND t.name = @p2
ORDER BY c.column_id;
	`
	rows, err := db.QueryContext(ctx, columnsSql, schema, name)
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

func getMssqlIndexesInfo(ctx context.Context, db *sql.DB, schema, name string) ([]*model.Index, error) {
	const indexesSql = `
		SELECT i.name                                AS index_name,
		       i.is_unique                           AS is_unique,
		       i.is_primary_key                      AS is_primary,
		       ic.key_ordinal                        AS ordinal,
		       c.name                                AS column_name
		FROM sys.indexes i
		JOIN sys.tables t ON i.object_id = t.object_id
		JOIN sys.schemas s ON t.schema_id = s.schema_id
		JOIN sys.index_columns ic ON i.object_id = ic.object_id AND i.index_id = ic.index_id
		JOIN sys.columns c ON ic.object_id = c.object_id AND ic.column_id = c.column_id
		WHERE s.name = @p1
		  AND t.name = @p2
		ORDER BY i.name, ic.key_ordinal;
	`
	rows, err := db.QueryContext(ctx, indexesSql, schema, name)
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
