package errors

import (
	goerr "errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichiels/go-pkg/language"
	"github.com/reyhanmichiels/go-pkg/operator"
)

type App struct {
	Code  codes.Code `json:"code"`
	Title string     `json:"title"`
	Body  string     `json:"body"`
	sys   error
}

func (e *App) Error() string {
	return e.sys.Error()
}

func Compile(err error, lang string) (int, App) {
	code := GetCode(err)

	if appErr, ok := codes.ErrorMessage[code]; ok {
		return appErr.StatusCode, App{
			Code:  code,
			Title: operator.Ternary(lang == language.Indonesian, appErr.TitleID, appErr.TitleEN),
			Body:  operator.Ternary(lang == language.Indonesian, appErr.BodyID, appErr.BodyEN),
			sys:   err,
		}
	}

	// Default Error
	return http.StatusInternalServerError, App{
		Code:  code,
		Title: "Service Error Not Defined",
		Body:  "Unknown error. Please contact admin",
		sys:   err,
	}
}

func NewWithCode(code codes.Code, msg string, val ...interface{}) error {
	return create(nil, code, msg, val...)
}

func create(cause error, code codes.Code, msg string, val ...interface{}) error {
	if code == codes.NoCode {
		code = GetCode(cause)
	}

	err := &stacktrace{
		message: fmt.Sprintf(msg, val...),
		cause:   cause,
		code:    code,
	}

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return err
	}
	err.file, err.line = file, line

	f := runtime.FuncForPC(pc)
	if f == nil {
		return err
	}
	err.function = shortFuncName(f)

	return err
}

func shortFuncName(f *runtime.Func) string {
	longName := f.Name()

	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}

// Implement golang errors.Is, reports whether any error in err's chain matches target.
func Is(err error, target error) bool {
	return goerr.Is(err, target)
}

// Implement golang errors.As, finds the first error in err's chain that matches target,
// and if one is found, sets target to that error value and returns true
func As(err error, target any) bool {
	return goerr.As(err, target)
}

func GetCode(err error) codes.Code {
	if err, ok := err.(*stacktrace); ok { // nolint:errorlint
		return err.code
	}
	return codes.NoCode
}

func GetCaller(err error) (string, int, string, error) {
	st, ok := err.(*stacktrace) // nolint:errorlint
	if !ok {
		return "", 0, "", create(nil, codes.NoCode, "failed to cast to stacktrace")
	}

	return st.file, st.line, st.message, nil
}

func GetCallerString(err error) (string, error) {
	file, line, msg, err := GetCaller(err)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%#v --- %s", file, line, msg), nil
}
