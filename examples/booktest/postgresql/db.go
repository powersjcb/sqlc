package booktest

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)



type PGTX interface {
	//ExecEx(context.Context, string, ...interface{}) (sql.Result, error)
	//PrepareEx(context.Context, string) (*sql.Stmt, error)
	//QueryEx(context.Context, string, ...interface{}) (*sql.Rows, error)
	//QueryRowEx(context.Context, string, ...interface{}) *sql.Row

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
