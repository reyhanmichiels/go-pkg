package codes

import (
	"math"

	"github.com/reyhanmichiels/go-pkg/v2/language"
	"github.com/reyhanmichiels/go-pkg/v2/operator"
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
	CodeAuth = Code(iota + 1000)
	CodeAuthRefreshTokenExpired
	CodeAuthAccessTokenExpired
	CodeAuthFailure
	CodeAuthInvalidToken
	CodeAuthRevokeRefreshTokenFailed
)

// json parser error
const (
	CodeJSONMarshalError = Code(iota + 1100)
	CodeJSONUnmarshalError
)

// SQL error
const (
	CodeSQL = Code(iota + 1200)
	CodeSQLInit
	CodeSQLBuilder
	CodeSQLTxBegin
	CodeSQLTxCommit
	CodeSQLTxRollback
	CodeSQLTxExec
	CodeSQLPrepareStmt
	CodeSQLRead
	CodeSQLRowScan
	CodeSQLRecordDoesNotExist
	CodeSQLUniqueConstraint
	CodeSQLConflict
	CodeSQLNoRowsAffected
)

// Cache Error
const (
	CodeRedisGet = Code(iota + 3900)
	CodeRedisSetEx
	CodeFailedLock
	CodeFailedReleaseLock
	CodeLockExist
	CodeCacheMarshal
	CodeCacheUnmarshal
	CodeCacheGetSimpleKey
	CodeCacheSetSimpleKey
	CodeCacheDeleteSimpleKey
	CodeCacheGetHashKey
	CodeCacheSetHashKey
	CodeCacheDeleteHashKey
	CodeCacheSetExpiration
	CodeCacheDecode
	CodeCacheLockNotAcquired
	CodeCacheInvalidCastType
	CodeCacheNotFound
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

	// Code JSON Parser
	CodeJSONMarshalError:   ErrMsgBadRequest,
	CodeJSONUnmarshalError: ErrMsgBadRequest,

	// Code SQL
	CodeSQL:                   ErrMsgInternalServerError,
	CodeSQLInit:               ErrMsgInternalServerError,
	CodeSQLBuilder:            ErrMsgInternalServerError,
	CodeSQLTxBegin:            ErrMsgInternalServerError,
	CodeSQLTxCommit:           ErrMsgInternalServerError,
	CodeSQLTxRollback:         ErrMsgInternalServerError,
	CodeSQLTxExec:             ErrMsgInternalServerError,
	CodeSQLPrepareStmt:        ErrMsgInternalServerError,
	CodeSQLRead:               ErrMsgInternalServerError,
	CodeSQLRowScan:            ErrMsgInternalServerError,
	CodeSQLRecordDoesNotExist: ErrMsgNotFound,
	CodeSQLUniqueConstraint:   ErrMsgConflict,
	CodeSQLConflict:           ErrMsgConflict,
	CodeSQLNoRowsAffected:     ErrMsgInternalServerError,

	// Code Cache Error
	CodeLockExist:            ErrMsgLockExist,
	CodeRedisGet:             ErrMsgInternalServerError,
	CodeRedisSetEx:           ErrMsgInternalServerError,
	CodeFailedLock:           ErrMsgInternalServerError,
	CodeFailedReleaseLock:    ErrMsgInternalServerError,
	CodeCacheMarshal:         ErrMsgInternalServerError,
	CodeCacheUnmarshal:       ErrMsgInternalServerError,
	CodeCacheGetSimpleKey:    ErrMsgInternalServerError,
	CodeCacheSetSimpleKey:    ErrMsgInternalServerError,
	CodeCacheDeleteSimpleKey: ErrMsgInternalServerError,
	CodeCacheGetHashKey:      ErrMsgInternalServerError,
	CodeCacheSetHashKey:      ErrMsgInternalServerError,
	CodeCacheDeleteHashKey:   ErrMsgInternalServerError,
	CodeCacheSetExpiration:   ErrMsgInternalServerError,
	CodeCacheDecode:          ErrMsgInternalServerError,
	CodeCacheLockNotAcquired: ErrMsgInternalServerError,
	CodeCacheInvalidCastType: ErrMsgInternalServerError,
	CodeCacheNotFound:        ErrMsgInternalServerError,
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
