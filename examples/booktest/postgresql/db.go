package booktest

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)



type PGTX interface {
	Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func New(db PGTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db PGTX
}

func (q *Queries) WithTx(db pgx.Tx) *Queries {
	return &Queries{
		db: db,
	}
}
