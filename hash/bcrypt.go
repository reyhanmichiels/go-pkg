package hash

import (
	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichiels/go-pkg/errors"
	bcryptLib "golang.org/x/crypto/bcrypt"
)

type bcryptInterface interface {
	GenerateFromText(text string) (string, error)
	CompareHashWithText(hashedText string, plainText string) bool
}

type bcrypt struct {
}

func initBcrypt() bcryptInterface {
	return &bcrypt{}
}

// GenerateFromText hash text with bcrypt algorithm
// using default cost (10)
func (b *bcrypt) GenerateFromText(text string) (string, error) {
	hashedTextByte, err := bcryptLib.GenerateFromPassword([]byte(text), bcryptLib.DefaultCost)
	if err != nil {
		return "", errors.NewWithCode(codes.CodeInternalServerError, err.Error())
	}

	return string(hashedTextByte), nil
}

// CompareHashWithText compare hashed text with plain text
// and return true if both is equal
func (b *bcrypt) CompareHashWithText(hashedText string, plainText string) bool {
	err := bcryptLib.CompareHashAndPassword([]byte(hashedText), []byte(plainText))
	return err == nil
}
