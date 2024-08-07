package auth

import "github.com/golang-jwt/jwt/v5"

type Claim struct {
	UserID    int64
	TokenType string
	jwt.RegisteredClaims
}

type User struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	RoleID int64  `json:"roleID"`
}
