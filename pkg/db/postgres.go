package db

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"sort"

	"github.com/DanielLiu1123/gencoder/pkg/model"
)

// GenPostgresTable generates a PostgreSQL table and fills the Table structure.
func GenPostgresTable(ctx context.Context, db *sql.DB, schema string, name string, ignoreColumns []string) (*model.Table, error) {
	t, err := getPostgresTableInfo(ctx, db, schema, name)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	columns, err := getPostgresColumnsInfo(ctx, db, schema, name, ignoreColumns)
	if err != nil {
		return nil, err
	}
	t.Columns = columns

	indexes, err := getPostgresIndexesInfo(ctx, db, schema, name)
	if err != nil {
		return nil, err
	}
	t.Indexes = indexes

	return t, nil
}

func getPostgresTableInfo(ctx context.Context, db *sql.DB, schema, name string) (*model.Table, error) {
	const tableSql = `
		SELECT table_schema,
			   table_name,
			   obj_description(pg_class.oid) AS table_comment
		FROM information_schema.tables
		JOIN pg_class ON relname = table_name
		WHERE table_schema = $1
		  AND table_name = $2;
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

func getPostgresColumnsInfo(ctx context.Context, db *sql.DB, schema, name string, ignoreColumns []string) ([]*model.Column, error) {
	const columnsSql = `
		SELECT a.attnum                             AS ordinal,
			   a.attname                            AS column_name,
			   format_type(a.atttypid, a.atttypmod) AS data_type,
			   NOT a.attnotnull                     AS is_nullable,
			   pg_get_expr(ad.adbin, ad.adrelid)    AS default_value,
			   COALESCE(ct.contype = 'p', false)    AS is_primary,
			   d.description                        AS comment
		FROM pg_attribute a
				 JOIN pg_class c ON c.oid = a.attrelid
				 JOIN pg_namespace n ON n.oid = c.relnamespace
				 LEFT JOIN pg_constraint ct ON ct.conrelid = c.oid AND a.attnum = ANY (ct.conkey) AND ct.contype = 'p'
				 LEFT JOIN pg_attrdef ad ON ad.adrelid = c.oid AND ad.adnum = a.attnum
				 LEFT JOIN pg_description d ON d.objoid = c.oid AND d.objsubid = a.attnum
		WHERE a.attisdropped = false
		  AND n.nspname = $1
		  AND c.relname = $2
		  AND a.attnum > 0
		ORDER BY a.attnum;
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

func getPostgresIndexesInfo(ctx context.Context, db *sql.DB, schema, name string) ([]*model.Index, error) {
	const indexesSql = `
		SELECT ic.relname                                                                 AS index_name,
			   i.indisunique                                                              AS is_unique,
			   i.indisprimary                                                             AS is_primary,
			   row_number() OVER (PARTITION BY ic.relname ORDER BY indkey_col.ordinality) AS ordinal,
			   a.attname                                                                  AS column_name
		FROM pg_index i
				 JOIN pg_class c ON c.oid = i.indrelid
				 JOIN pg_namespace n ON n.oid = c.relnamespace
				 JOIN pg_class ic ON ic.oid = i.indexrelid
				 LEFT JOIN LATERAL unnest(i.indkey) WITH ORDINALITY indkey_col(indkey_col, ordinality) ON TRUE
				 LEFT JOIN pg_attribute a ON i.indrelid = a.attrelid
			AND a.attnum = indkey_col.indkey_col
			AND a.attisdropped = false
		WHERE n.nspname = $1
		  AND c.relname = $2
		ORDER BY ic.relname, indkey_col.ordinality;
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
