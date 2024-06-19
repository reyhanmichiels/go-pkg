package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichiels/go-pkg/errors"
	"github.com/reyhanmichiels/go-pkg/log"
)

var ErrNotObtained = redislock.ErrNotObtained

const (
	Nil = redis.Nil
)

type Locker *redislock.Lock

type Interface interface {
	Get(ctx context.Context, key string) (string, error)
	SetEX(ctx context.Context, key string, val string, expTime time.Duration) error
	Lock(ctx context.Context, key string, expTime time.Duration) (*redislock.Lock, error)
	LockRelease(ctx context.Context, lock *redislock.Lock) error
	Del(ctx context.Context, key string) error
	FlushAll(ctx context.Context) error
	FlushAllAsync(ctx context.Context) error
	FlushDB(ctx context.Context) error
	FlushDBAsync(ctx context.Context) error
	GetDefaultTTL(ctx context.Context) time.Duration
}

type TLSConfig struct {
	Enabled            bool
	InsecureSkipVerify bool
}

type Config struct {
	Protocol   string
	Host       string
	Port       string
	Username   string
	Password   string
	DefaultTTL time.Duration
	TLS        TLSConfig
}

type cache struct {
	conf  Config
	rdb   *redis.Client
	log   log.Interface
	rlock *redislock.Client
}

func Init(cfg Config, log log.Interface) Interface {
	c := &cache{
		conf: cfg,
		log:  log,
	}
	c.connect(context.Background())
	return c
}

func (c *cache) connect(ctx context.Context) {
	redisOpts := redis.Options{
		Network:  c.conf.Protocol,
		Addr:     fmt.Sprintf("%s:%s", c.conf.Host, c.conf.Port),
		Username: c.conf.Username,
		Password: c.conf.Password,
	}

	if c.conf.TLS.Enabled {
		redisOpts.TLSConfig = &tls.Config{
			InsecureSkipVerify: c.conf.TLS.InsecureSkipVerify, // #nosec G402
		}
	}

	client := redis.NewClient(&redisOpts)

	err := client.Ping(ctx).Err()
	if err != nil {
		c.log.Fatal(ctx, fmt.Sprintf("[FATAL] cannot connect to redis on address @%s:%v, with error: %s", c.conf.Host, c.conf.Port, err))
	}
	c.rdb = client
	c.log.Info(ctx, fmt.Sprintf("REDIS: Address @%s:%v", c.conf.Host, c.conf.Port))

	c.rlock = redislock.New(client)
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	s, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return s, err
	}

	return s, nil
}

func (c *cache) SetEX(ctx context.Context, key string, val string, expTime time.Duration) error {

	if expTime <= 0 {
		expTime = c.conf.DefaultTTL
	}

	err := c.rdb.SetEx(ctx, key, val, expTime).Err()
	if err != nil {
		return errors.NewWithCode(codes.CodeInternalServerError, err.Error())
	}

	return nil
}

func (c *cache) Lock(ctx context.Context, key string, expTime time.Duration) (*redislock.Lock, error) {
	// Obtain lock
	lock, err := c.rlock.Obtain(ctx, key, expTime, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		return nil, err
	} else if err != nil {
		return nil, errors.NewWithCode(codes.CodeInternalServerError, err.Error())
	}

	return lock, nil
}

func (c *cache) LockRelease(ctx context.Context, lock *redislock.Lock) error {
	if lock != nil {
		err := lock.Release(ctx)
		if errors.Is(err, redislock.ErrLockNotHeld) {
			return err
		} else if err != nil {
			return errors.NewWithCode(codes.CodeInternalServerError, err.Error())
		}
	}

	return nil
}

func (c *cache) Del(ctx context.Context, key string) error {
	var keysCount int64
	// Use SCAN with COUNT = 0 to advance the cursor
	iter := c.rdb.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		c.log.Info(ctx, fmt.Sprintf("deleted key: %s", iter.Val()))
		c.rdb.Del(ctx, iter.Val())
		keysCount++
	}
	if err := iter.Err(); err != nil {
		return err
	}
	c.log.Info(ctx, fmt.Sprintf("successfully deleted %d numbers of key", keysCount))

	return nil
}

func (c *cache) FlushAll(ctx context.Context) error {
	return c.rdb.FlushAll(ctx).Err()
}

func (c *cache) FlushAllAsync(ctx context.Context) error {
	return c.rdb.FlushAllAsync(ctx).Err()
}

func (c *cache) FlushDB(ctx context.Context) error {
	return c.rdb.FlushDB(ctx).Err()
}

func (c *cache) FlushDBAsync(ctx context.Context) error {
	return c.rdb.FlushDBAsync(ctx).Err()
}

func (c *cache) GetDefaultTTL(ctx context.Context) time.Duration {
	return c.conf.DefaultTTL
}
