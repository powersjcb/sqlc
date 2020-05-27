// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package querytest

import (
	"context"
)

const cTEFilter = `-- name: CTEFilter :many
WITH filter_count AS (
	SELECT count(*) FROM bar WHERE ready = $1
)
SELECT filter_count.count
FROM filter_count
`

func (q *Queries) CTEFilter(ctx context.Context, ready bool) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, cTEFilter, ready)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]int64, 0)
	for rows.Next() {
		var count int64
		if err := rows.Scan(&count); err != nil {
			return nil, err
		}
		items = append(items, count)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
