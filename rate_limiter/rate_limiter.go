package rate_limiter

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/reyhanmichiels/go-pkg/log"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type Interface interface {
	Limiter() gin.HandlerFunc
}

type ConfigPath struct {
	Enabled bool
	Period  string
	Limit   int64
	Path    string
}

type Config struct {
	Enabled bool
	Period  string
	Limit   int64
	Paths   []ConfigPath
}

type rateLimiter struct {
	cfg             Config
	log             log.Interface
	middleware      gin.HandlerFunc
	middlewarePaths map[string]GinMiddleware
}

func Init(cfg Config, log log.Interface) Interface {
	rl := &rateLimiter{
		cfg: cfg,
		log: log,
	}

	rl.InitConfiguration()

	return rl
}

func (rl *rateLimiter) InitConfiguration() {
	if rl.cfg.Enabled {
		paths := []string{}

		middlewares := make(map[string]GinMiddleware)
		for _, conf := range rl.cfg.Paths {
			paths = append(paths, conf.Path)

			if !conf.Enabled {
				continue
			}

			limiterPath := getLimiter(Config{Period: conf.Period, Limit: conf.Limit}, rl.log)
			middlewares[conf.Path] = ginMiddleware(limiterPath)
		}

		rl.middlewarePaths = middlewares
		rl.middleware = rl.InitMiddleware(paths, getLimiter(rl.cfg, rl.log))
		return
	}

	rl.middleware = skipMiddleware()
}

func getLimiter(conf Config, log log.Interface) *limiter.Limiter {
	ctx := context.Background()
	time, err := time.ParseDuration(conf.Period)
	if err != nil {
		log.Fatal(ctx, err)
	}

	rate := limiter.Rate{
		Period: time,
		Limit:  conf.Limit,
	}

	store := memory.NewStore()
	return limiter.New(store, rate)
}

func (rl *rateLimiter) Limiter() gin.HandlerFunc {
	return rl.middleware
}
