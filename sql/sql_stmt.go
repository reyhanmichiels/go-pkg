package sql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type CommandStmt interface {
	Close() error
	Get(name string, dest interface{}, args ...interface{}) error
	QueryRow(name string, args ...interface{}) (*sqlx.Row, error)
	Query(name string, args ...interface{}) (*sqlx.Rows, error)
	Exec(name string, args ...interface{}) (sql.Result, error)
}

type commandStmt struct {
	ctx  context.Context
	stmt *sqlx.Stmt
}

func initStmt(ctx context.Context, stmt *sqlx.Stmt) CommandStmt {
	c := &commandStmt{
		ctx:  ctx,
		stmt: stmt,
	}

	return c
}

func (c *commandStmt) Close() error {
	return c.stmt.Close()
}

func (c *commandStmt) Get(name string, dest interface{}, args ...interface{}) error {
	return c.stmt.GetContext(c.ctx, dest, args...)
}

func (c *commandStmt) QueryRow(name string, args ...interface{}) (*sqlx.Row, error) {
	return c.stmt.QueryRowxContext(c.ctx, args...), nil
}

func (c *commandStmt) Query(name string, args ...interface{}) (*sqlx.Rows, error) {
	return c.stmt.QueryxContext(c.ctx, args...)
}

func (c *commandStmt) Exec(name string, args ...interface{}) (sql.Result, error) {
	return c.stmt.ExecContext(c.ctx, args...)
}
