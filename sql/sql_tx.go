package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/log"
)

type CommandTx interface {
	Commit() error
	Rollback()

	QueryRow(name string, query string, args ...interface{}) (*sqlx.Row, error)
	Query(name string, query string, args ...interface{}) (*sqlx.Rows, error)
	Get(name string, query string, dest interface{}, args ...interface{}) error

	Prepare(name string, query string) (CommandStmt, error)
	PrepareNamed(name string, query string) (NamedCommandStmt, error)
	NamedExec(name string, query string, args interface{}) (sql.Result, error)
	Exec(name string, query string, args ...interface{}) (sql.Result, error)
}

type commandTx struct {
	ctx      context.Context
	name     string
	tx       *sqlx.Tx
	log      log.Interface
	logQuery bool
}

func initTx(ctx context.Context, name string, tx *sqlx.Tx, log log.Interface, logQuery bool) CommandTx {
	c := &commandTx{
		ctx:      ctx,
		name:     name,
		tx:       tx,
		log:      log,
		logQuery: logQuery,
	}

	return c
}

func (x *commandTx) Commit() error {
	return x.tx.Commit()
}

// Rollback needs to be called with defer right after calling BeginTx.
// Read here: https://go.dev/doc/database/execute-transactions.
func (x *commandTx) Rollback() {
	if err := x.tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		x.log.Error(x.ctx, err)
	}
}

func (x *commandTx) QueryRow(name string, query string, args ...interface{}) (*sqlx.Row, error) {
	if x.logQuery {
		x.log.Info(x.ctx, fmt.Sprintf(queryLogMessage, name, replaceBindVarsWithArgs(query, args...)))
	}
	row := x.tx.QueryRowxContext(x.ctx, query, args...)
	return row, row.Err()
}

func (x *commandTx) Query(name string, query string, args ...interface{}) (*sqlx.Rows, error) {
	if x.logQuery {
		x.log.Info(x.ctx, fmt.Sprintf(queryLogMessage, name, replaceBindVarsWithArgs(query, args...)))
	}
	return x.tx.QueryxContext(x.ctx, query, args...)
}

func (x *commandTx) NamedExec(name string, query string, args interface{}) (sql.Result, error) {
	if x.logQuery {
		x.log.Info(x.ctx, fmt.Sprintf(queryLogMessage, name, replaceBindVarsWithArgs(query)))
	}
	return x.tx.NamedExecContext(x.ctx, query, args)
}

func (x *commandTx) Prepare(name string, query string) (CommandStmt, error) {
	if x.logQuery {
		x.log.Info(x.ctx, fmt.Sprintf(queryLogMessage, name, replaceBindVarsWithArgs(query)))
	}
	stmt, err := x.tx.PreparexContext(x.ctx, query)
	if err != nil {
		return nil, err
	}
	return initStmt(x.ctx, stmt), nil
}

func (x *commandTx) PrepareNamed(name string, query string) (NamedCommandStmt, error) {
	if x.logQuery {
		x.log.Info(x.ctx, fmt.Sprintf(queryLogMessage, name, replaceBindVarsWithArgs(query)))
	}
	stmt, err := x.tx.PrepareNamedContext(x.ctx, query)
	if err != nil {
		return nil, err
	}
	return initNamedStmt(x.ctx, stmt), nil
}

func (x *commandTx) Exec(name string, query string, args ...interface{}) (sql.Result, error) {
	if x.logQuery {
		x.log.Info(x.ctx, fmt.Sprintf(queryLogMessage, name, replaceBindVarsWithArgs(query, args...)))
	}
	return x.tx.ExecContext(x.ctx, query, args...)
}

func (x *commandTx) Get(name string, query string, dest interface{}, args ...interface{}) error {
	if x.logQuery {
		x.log.Info(x.ctx, fmt.Sprintf(queryLogMessage, name, replaceBindVarsWithArgs(query, args...)))
	}
	return x.tx.GetContext(x.ctx, dest, query, args...)
}
