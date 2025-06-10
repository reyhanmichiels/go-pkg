package sql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type NamedCommandStmt interface {
	Close() error
	Get(dest interface{}, arg interface{}) error
	QueryRow(arg interface{}) *sqlx.Row
	Query(arg interface{}) (*sqlx.Rows, error)
	Exec(arg interface{}) (sql.Result, error)
}

type namedCommandStmt struct {
	ctx  context.Context
	stmt *sqlx.NamedStmt
}

func initNamedStmt(ctx context.Context, stmt *sqlx.NamedStmt) NamedCommandStmt {
	return &namedCommandStmt{
		ctx:  ctx,
		stmt: stmt,
	}
}

func (n *namedCommandStmt) Close() error {
	return n.stmt.Close()
}

func (n *namedCommandStmt) Get(dest interface{}, arg interface{}) error {
	return n.stmt.GetContext(n.ctx, dest, arg)
}

func (n *namedCommandStmt) QueryRow(arg interface{}) *sqlx.Row {
	return n.stmt.QueryRowxContext(n.ctx, arg)
}

func (n *namedCommandStmt) Query(arg interface{}) (*sqlx.Rows, error) {
	return n.stmt.QueryxContext(n.ctx, arg)
}

func (n *namedCommandStmt) Exec(arg interface{}) (sql.Result, error) {
	return n.stmt.ExecContext(n.ctx, arg)
}
