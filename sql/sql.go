package sql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/reyhanmichiels/go-pkg/v2/log"
)

var ErrNotFound = sql.ErrNoRows

const (
	failedConnectDBMessage = "[FATAL] cannot connect to db %s leader: %s on port %d, with error: %s"
)

type txKey struct{} // txKey is a context key for the transaction.

type Interface interface {
	Leader() Command
	Follower() Command
	Stop()

	Rebind(query string) string

	QueryRow(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Row, error)
	Query(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error)
	Get(ctx context.Context, name string, query string, dest interface{}, args ...interface{}) error

	NamedExec(ctx context.Context, name string, query string, args interface{}) (sql.Result, error)
	Exec(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error)
	Transaction(ctx context.Context, name string, txOpts TxOptions, f func(context.Context) error) error
}

type sqlDB struct {
	endOnce  *sync.Once
	leader   Command
	follower Command
	cfg      Config
	log      log.Interface
}

type Config struct {
	LogQuery bool
	Driver   string
	Leader   ConnConfig
	Follower ConnConfig
}

type ConnConfig struct {
	Host     string
	Port     int
	DB       string
	User     string
	Password string
	SSL      bool
	Schema   string
	Options  ConnOptions
	MockDB   *sql.DB
}

type ConnOptions struct {
	MaxLifeTime time.Duration
	MaxIdle     int
	MaxOpen     int
}

func Init(cfg Config, log log.Interface) Interface {
	sqlDB := sqlDB{
		endOnce: &sync.Once{},
		cfg:     cfg,
		log:     log,
	}

	sqlDB.initDB()

	return &sqlDB
}

func (s *sqlDB) Leader() Command {
	return s.leader
}

func (s *sqlDB) Follower() Command {
	return s.follower
}

func (s *sqlDB) Stop() {
	s.endOnce.Do(func() {
		ctx := context.Background()
		if s.leader != nil {
			if err := s.leader.Close(); err != nil {
				s.log.Error(ctx, err)
			}
		}

		if s.follower != nil {
			if err := s.follower.Close(); err != nil {
				s.log.Error(ctx, err)
			}
		}
	})
}

func (s *sqlDB) Rebind(query string) string {
	return s.leader.Rebind(query)
}

func (s *sqlDB) QueryRow(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Row, error) {
	return s.follower.QueryRow(ctx, name, query, args...)
}

func (s *sqlDB) Query(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error) {
	return s.follower.Query(ctx, name, query, args...)
}

func (s *sqlDB) Get(ctx context.Context, name string, query string, dest interface{}, args ...interface{}) error {
	return s.follower.Get(ctx, name, query, dest, args...)
}

func (s *sqlDB) NamedExec(ctx context.Context, name string, query string, args interface{}) (sql.Result, error) {
	if tx, ok := s.getTx(ctx); ok {
		return tx.NamedExec(name, query, args)
	}
	return s.leader.NamedExec(ctx, name, query, args)
}

func (s *sqlDB) Exec(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error) {
	if tx, ok := s.getTx(ctx); ok {
		return tx.Exec(name, query, args...)
	}
	return s.leader.Exec(ctx, name, query, args...)
}

// Transaction executes a transaction. If the given function returns an error, the transaction
// is rolled back. Otherwise, it is automatically committed before `Transaction()` returns.
func (s *sqlDB) Transaction(ctx context.Context, name string, txOpts TxOptions, f func(context.Context) error) error {
	tx, err := s.leader.BeginTx(ctx, name, txOpts)
	if err != nil {
		return err
	}
	c := context.WithValue(ctx, txKey{}, tx)
	err = f(c)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// getTx retrieves the transaction from the context.
func (s *sqlDB) getTx(ctx context.Context) (CommandTx, bool) {
	tx, ok := ctx.Value(txKey{}).(CommandTx)
	return tx, ok
}

func (s *sqlDB) initDB() {
	ctx := context.Background()

	db := s.connect(true)
	s.log.Info(ctx, fmt.Sprintf("SQL: [LEADER] driver=%s db=%s @%s:%v ssl=%v", s.cfg.Driver, s.cfg.Leader.DB, s.cfg.Leader.Host, s.cfg.Leader.Port, s.cfg.Leader.SSL))

	s.leader = initCommand(db, s.log, s.cfg.LogQuery)
	s.follower = s.leader

	if s.isFollowerEnabled() {
		db = s.connect(false)
		s.log.Info(ctx, fmt.Sprintf("SQL: [FOLLOWER] driver=%s db=%s @%s:%v ssl=%v", s.cfg.Driver, s.cfg.Follower.DB, s.cfg.Follower.Host, s.cfg.Follower.Port, s.cfg.Leader.SSL))

		s.follower = initCommand(db, s.log, s.cfg.LogQuery)
	}
}

func (s *sqlDB) connect(isLeader bool) *sqlx.DB {
	conf := s.cfg.Leader
	if !isLeader {
		conf = s.cfg.Follower
	}

	if conf.MockDB != nil {
		return sqlx.NewDb(s.cfg.Leader.MockDB, s.cfg.Driver)
	}

	uri, err := s.getURI(conf)
	if err != nil {
		s.log.Fatal(context.Background(), fmt.Sprintf("[FATAL] cannot get URI for db %s leader: %s on port %d, with error: %s", conf.DB, conf.Host, conf.Port, err))
	}

	sqlxDB, err := sqlx.Open(s.cfg.Driver, uri)
	if err != nil {
		s.log.Fatal(context.Background(), fmt.Sprintf(failedConnectDBMessage, conf.DB, conf.Host, conf.Port, err))
	}

	err = sqlxDB.Ping()
	if err != nil {
		s.log.Fatal(context.Background(), fmt.Sprintf(failedConnectDBMessage, conf.DB, conf.Host, conf.Port, err))
	}

	sqlxDB.SetMaxOpenConns(conf.Options.MaxOpen)
	sqlxDB.SetMaxIdleConns(conf.Options.MaxIdle)
	sqlxDB.SetConnMaxLifetime(conf.Options.MaxLifeTime)

	return sqlxDB
}

func (s *sqlDB) getURI(conf ConnConfig) (string, error) {
	switch s.cfg.Driver {
	case "postgres":
		ssl := `disable`
		if conf.SSL {
			ssl = `require`
		}
		if conf.Schema == "" {
			conf.Schema = "public"
		}
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=%s", conf.Host, conf.Port, conf.User, conf.Password, conf.DB, conf.Schema, ssl), nil
	case "mysql":
		ssl := `false`
		if conf.SSL {
			ssl = `true`
		}
		return fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?tls=%s&parseTime=true", conf.User, conf.Password, conf.Host, conf.Port, conf.DB, ssl), nil
	default:
		return "", fmt.Errorf(`DB Driver [%s] is not supported`, s.cfg.Driver)
	}
}

func (s *sqlDB) isFollowerEnabled() bool {
	isHostNotEmpty := s.cfg.Follower.Host != ""
	isHostDifferent := s.cfg.Follower.Host != s.cfg.Leader.Host && s.cfg.Follower.Port == s.cfg.Leader.Port
	isPortDifferent := s.cfg.Follower.Host == s.cfg.Leader.Host && s.cfg.Follower.Port != s.cfg.Leader.Port
	return isHostNotEmpty && (isHostDifferent || isPortDifferent)
}
