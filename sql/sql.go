package sql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/reyhanmichiels/go-pkg/log"
)

var ErrNotFound = sql.ErrNoRows

const (
	failedConnectDBMessage = "[FATAL] cannot connect to db %s leader: %s on port %d, with error: %s"
)

type Interface interface {
	Leader() Command
	Follower() Command
	Stop()
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
	if s.cfg.Driver != "mysql" {
		s.log.Fatal(context.Background(), fmt.Sprintf("driver %s is not supported", s.cfg.Driver))
	}

	conf := s.cfg.Leader
	if !isLeader {
		conf = s.cfg.Follower
	}

	if conf.MockDB != nil {
		return sqlx.NewDb(s.cfg.Leader.MockDB, s.cfg.Driver)
	}

	ssl := `false`
	if conf.SSL {
		ssl = `true`
	}
	uri := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?tls=%s&parseTime=true", conf.User, conf.Password, conf.Host, conf.Port, conf.DB, ssl)

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

func (s *sqlDB) isFollowerEnabled() bool {
	isHostNotEmpty := s.cfg.Follower.Host != ""
	isHostDifferent := s.cfg.Follower.Host != s.cfg.Leader.Host && s.cfg.Follower.Port == s.cfg.Leader.Port
	isPortDifferent := s.cfg.Follower.Host == s.cfg.Leader.Host && s.cfg.Follower.Port != s.cfg.Leader.Port
	return isHostNotEmpty && (isHostDifferent || isPortDifferent)
}
