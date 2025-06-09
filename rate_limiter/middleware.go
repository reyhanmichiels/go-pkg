package rate_limiter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/reyhanmichiels/go-pkg/v2/appcontext"
	"github.com/reyhanmichiels/go-pkg/v2/checker"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
)

type GinMiddleware struct {
	*mgin.Middleware
}

type HTTPResp struct {
	Message HTTPMessage `json:"message"`
	Meta    Meta        `json:"metadata"`
}

type HTTPMessage struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Meta struct {
	Path       string `json:"path"`
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
}

func (rl *rateLimiter) InitMiddleware(paths []string, limiter *limiter.Limiter) gin.HandlerFunc {
	middleware := ginMiddleware(limiter)

	return func(ctx *gin.Context) {
		if checker.ArrayContains(paths, ctx.Request.URL.Path) {
			middleware, ok := rl.middlewarePaths[ctx.Request.URL.Path]
			if !ok {
				ctx.Next()
				return
			}

			middleware.Handle(ctx)
			return
		}

		middleware.Handle(ctx)
	}
}

func ginMiddleware(limiter *limiter.Limiter) GinMiddleware {
	return GinMiddleware{
		&mgin.Middleware{
			Limiter:        limiter,
			OnError:        mgin.DefaultErrorHandler,
			OnLimitReached: limitReachedHandler,
			KeyGetter:      mgin.DefaultKeyGetter,
			ExcludedKey:    nil,
		},
	}
}

func skipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func limitReachedHandler(c *gin.Context) {
	httpStatus, displayError := errors.Compile(errors.NewWithCode(codes.CodeTooManyRequest, "Limit Exceeded"), appcontext.GetAcceptLanguage(c))

	c.JSON(http.StatusTooManyRequests, HTTPResp{
		Message: HTTPMessage{
			Title: displayError.Title,
			Body:  displayError.Body,
		},
		Meta: Meta{
			Path:       c.Request.Host + c.Request.URL.String(),
			StatusCode: httpStatus,
			Status:     http.StatusText(httpStatus),
			Message:    fmt.Sprintf("%s %s [%d] %s", c.Request.Method, c.Request.URL.RequestURI(), httpStatus, displayError.Error()),
			Timestamp:  time.Now().Format(time.RFC3339),
		},
	})
}
