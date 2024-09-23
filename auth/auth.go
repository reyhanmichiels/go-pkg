package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichiels/go-pkg/errors"
	"github.com/reyhanmichiels/go-pkg/log"
)

type contextKey string

const (
	userAuthInfo contextKey = "UserAuthInfo"
)

type Config struct {
	AccessTokenExpireTime  time.Duration
	RefreshTokenExpireTime time.Duration
	AccessTokenType        string
	RefreshTokenType       string
	SigningKey             string
}

type Interface interface {
	CreateAccessToken(userID int64) (string, error)
	CreateRefreshToken(userID int64) (string, error)
	ValidateAccessToken(token string) (int64, error)
	ValidateRefreshToken(token string) error
	SetUserAuthInfo(ctx context.Context, user User) context.Context
	GetUserAuthInfo(ctx context.Context) (User, error)
}

type auth struct {
	cfg        Config
	log        log.Interface
	signingKey []byte
}

func Init(cfg Config, log log.Interface) Interface {
	if cfg.AccessTokenType == cfg.RefreshTokenType {
		log.Fatal(context.Background(), errors.NewWithCode(codes.CodeInternalServerError, "type value shouldn't be same"))
	}

	if cfg.AccessTokenType == "" || cfg.RefreshTokenType == "" {
		log.Fatal(context.Background(), errors.NewWithCode(codes.CodeInternalServerError, "type value shouldn't be empty"))
	}

	return &auth{
		cfg:        cfg,
		log:        log,
		signingKey: []byte(cfg.SigningKey),
	}
}

func (a *auth) CreateAccessToken(userID int64) (string, error) {
	expireTime := time.Now().Add(a.cfg.AccessTokenExpireTime)

	claim := &Claim{
		UserID:    userID,
		TokenType: a.cfg.AccessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString(a.signingKey)
	if err != nil {
		return signedToken, errors.NewWithCode(codes.CodeInternalServerError, err.Error())
	}

	return signedToken, nil
}

func (a *auth) CreateRefreshToken(userID int64) (string, error) {
	expireTime := time.Now().Add(a.cfg.RefreshTokenExpireTime)

	claim := &Claim{
		TokenType: a.cfg.RefreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString(a.signingKey)
	if err != nil {
		return signedToken, errors.NewWithCode(codes.CodeInternalServerError, err.Error())
	}

	return signedToken, nil
}

func (a *auth) ValidateAccessToken(token string) (int64, error) {
	claim := Claim{}

	_, err := jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (interface{}, error) {
		return a.signingKey, nil
	})
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return 0, errors.NewWithCode(codes.CodeAuthInvalidToken, err.Error())
	} else if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		return 0, errors.NewWithCode(codes.CodeAuthAccessTokenExpired, err.Error())
	}

	if claim.TokenType != a.cfg.AccessTokenType {
		return 0, errors.NewWithCode(codes.CodeAuthInvalidToken, "invalid token")
	}

	return claim.UserID, nil
}

func (a *auth) ValidateRefreshToken(token string) error {
	claim := Claim{}

	_, err := jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (interface{}, error) {
		return a.signingKey, nil
	})
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return errors.NewWithCode(codes.CodeAuthInvalidToken, err.Error())
	} else if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		return errors.NewWithCode(codes.CodeAuthAccessTokenExpired, err.Error())
	}

	if claim.TokenType != a.cfg.RefreshTokenType {
		return errors.NewWithCode(codes.CodeAuthInvalidToken, "invalid token")
	}

	return nil
}

func (a *auth) SetUserAuthInfo(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userAuthInfo, user)
}

func (a *auth) GetUserAuthInfo(ctx context.Context) (User, error) {
	user, ok := ctx.Value(userAuthInfo).(User)
	if !ok {
		return user, errors.NewWithCode(codes.CodeAuthFailure, "failed getting user auth info")
	}

	return user, nil
}
