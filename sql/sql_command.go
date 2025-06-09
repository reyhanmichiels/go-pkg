package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichiels/go-pkg/errors"
	"github.com/reyhanmichiels/go-pkg/log"
)

type Command interface {
	Close() error
	Ping(ctx context.Context) error

	Rebind(query string) string

	QueryRow(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Row, error)
	Query(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error)
	Get(ctx context.Context, name string, query string, dest interface{}, args ...interface{}) error

	NamedExec(ctx context.Context, name string, query string, args interface{}) (sql.Result, error)
	Exec(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error)
	BeginTx(ctx context.Context, name string, opts TxOptions) (CommandTx, error)
	ExecuteInTransaction(ctx context.Context, name string, opt TxOptions, fn func(tx CommandTx) error) error
}

type command struct {
	db       *sqlx.DB
	log      log.Interface
	logQuery bool
}

type TxOptions struct {
	Isolation sql.IsolationLevel
	ReadOnly  bool
}

func initCommand(db *sqlx.DB, log log.Interface, logQuery bool) Command {
	c := &command{
		db:       db,
		log:      log,
		logQuery: logQuery,
	}

	return c
}

func (c *command) Close() error {
	return c.db.Close()
}

func (c *command) Rebind(query string) string {
	return c.db.Rebind(query)
}

func (c *command) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c *command) QueryRow(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Row, error) {
	if c.logQuery {
		c.log.Info(ctx, fmt.Sprintf(queryLogMessage, name, replaceBindVarsWithArgs(query, args...)))
	}
	row := c.db.QueryRowxContext(ctx, query, args...)
	return row, row.Err()
}

func (c *command) Query(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error) {
	if c.logQuery {
		c.log.Info(ctx, fmt.Sprintf(queryLogMessage, name, replaceBindVarsWithArgs(query, args...)))
	}
	return c.db.QueryxContext(ctx, query, args...)
}

func (c *command) NamedExec(ctx context.Context, name string, query string, args interface{}) (sql.Result, error) {
	if c.logQuery {
		c.log.Info(ctx, fmt.Sprintf(queryLogMessage, name, replaceBindVarsWithArgs(query)))
	}
	return c.db.NamedExecContext(ctx, query, args)
}

func (c *command) Exec(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error) {
	if c.logQuery {
		c.log.Info(ctx, fmt.Sprintf(queryLogMessage, name, replaceBindVarsWithArgs(query, args...)))
	}
	return c.db.ExecContext(ctx, query, args...)
}

func (c *command) BeginTx(ctx context.Context, name string, opt TxOptions) (CommandTx, error) {
	opts := &sql.TxOptions{
		Isolation: opt.Isolation,
		ReadOnly:  opt.ReadOnly,
	}
	tx, err := c.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return initTx(ctx, name, tx, c.log, c.logQuery), nil
}

func (c *command) ExecuteInTransaction(ctx context.Context, name string, opt TxOptions, fn func(tx CommandTx) error) error {
	tx, err := c.BeginTx(ctx, name, opt)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = fn(tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.NewWithCode(codes.CodeInternalServerError, err.Error())
	}

	return nil
}

func (c *command) Get(ctx context.Context, name string, query string, dest interface{}, args ...interface{}) error {
	if c.logQuery {
		c.log.Info(ctx, fmt.Sprintf(queryLogMessage, name, replaceBindVarsWithArgs(query, args...)))
	}
	return c.db.GetContext(ctx, dest, query, args...)
}
