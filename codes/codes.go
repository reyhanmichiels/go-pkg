package codes

import (
	"math"

	"github.com/reyhanmichiels/go-pkg/language"
	"github.com/reyhanmichiels/go-pkg/operator"
)

type Code uint32

type AppMessage map[Code]Message

type DisplayMessage struct {
	StatusCode int    `json:"statusCode"`
	Title      string `json:"title"`
	Body       string `json:"body"`
}

const NoCode Code = math.MaxUint32

// success code
const (
	CodeSuccess = Code(iota + 10)
	CodeCreated
	CodeAccepted
)

// common errors
const (
	CodeInvalidValue = Code(iota + 100)
	CodeContextDeadlineExceeded
	CodeContextCanceled
	CodeInternalServerError
	CodeServerUnavailable
	CodeNotImplemented
	CodeBadRequest
	CodeNotFound
	CodeConflict
	CodeUnauthorized
	CodeForbidden
	CodeTooManyRequest
	CodeMarshal
	CodeUnmarshal
)

// translator errors
const (
	CodeTranslatorError = Code(iota + 500)
)

// auth error
const (
	CodeAuth = Code(iota + 1700)
	CodeAuthRefreshTokenExpired
	CodeAuthAccessTokenExpired
	CodeAuthFailure
	CodeAuthInvalidToken
	CodeAuthRevokeRefreshTokenFailed
)

var ErrorMessage = AppMessage{
	CodeInvalidValue:            ErrMsgBadRequest,
	CodeContextDeadlineExceeded: ErrMsgContextTimeout,
	CodeContextCanceled:         ErrMsgContextTimeout,
	CodeInternalServerError:     ErrMsgInternalServerError,
	CodeServerUnavailable:       ErrMsgServiceUnavailable,
	CodeNotImplemented:          ErrMsgNotImplemented,
	CodeBadRequest:              ErrMsgBadRequest,
	CodeNotFound:                ErrMsgNotFound,
	CodeConflict:                ErrMsgConflict,
	CodeUnauthorized:            ErrMsgUnauthorized,
	CodeForbidden:               ErrMsgForbidden,
	CodeTooManyRequest:          ErrMsgTooManyRequest,
	CodeMarshal:                 ErrMsgBadRequest,
	CodeUnmarshal:               ErrMsgBadRequest,

	// Code Translator
	CodeTranslatorError: ErrMsgTranslatorlib,

	// Code Auth
	CodeAuth:                         ErrMsgUnauthorized,
	CodeAuthRefreshTokenExpired:      ErrMsgRefreshTokenExpired,
	CodeAuthAccessTokenExpired:       ErrMsgAccessTokenExpired,
	CodeAuthFailure:                  ErrMsgUnauthorized,
	CodeAuthInvalidToken:             ErrMsgInvalidToken,
	CodeAuthRevokeRefreshTokenFailed: ErrMsgRevokeRefreshTokenFailed,
}

var SuccessMessage = AppMessage{
	CodeSuccess:  SuccessDefault,
	CodeCreated:  SuccessCreated,
	CodeAccepted: SuccessAccepted,
}

func Compile(code Code, lang string) DisplayMessage {
	if appMsg, ok := SuccessMessage[code]; ok {
		return DisplayMessage{
			StatusCode: appMsg.StatusCode,
			Title:      operator.Ternary(lang == language.Indonesian, appMsg.TitleID, appMsg.TitleEN),
			Body:       operator.Ternary(lang == language.Indonesian, appMsg.BodyID, appMsg.BodyEN),
		}
	}

	return DisplayMessage{
		StatusCode: SuccessDefault.StatusCode,
		Title:      operator.Ternary(lang == language.Indonesian, SuccessDefault.TitleID, SuccessDefault.TitleEN),
		Body:       operator.Ternary(lang == language.Indonesian, SuccessDefault.BodyID, SuccessDefault.BodyEN),
	}
}
