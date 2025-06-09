package translator

import (
	"context"
	"fmt"
	"os"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/reyhanmichiels/go-pkg/v2/appcontext"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/language"
	"github.com/reyhanmichiels/go-pkg/v2/log"
)

type Interface interface {
	Translate(ctx context.Context, key interface{}, params ...string) (string, error)
}

type Config struct {
	FallbackLanguageID   string
	SupportedLanguageIDs []string
	TranslationDir       string
}

type translator struct {
	translator *ut.UniversalTranslator
	log        log.Interface
}

func Init(conf Config, log log.Interface) Interface {
	fallback, supported, err := parseLanguageId(conf)
	if err != nil {
		log.Fatal(context.Background(), err)
	}

	t := ut.New(fallback, supported...)

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(context.Background(), err)
	}

	err = t.Import(ut.FormatJSON, fmt.Sprintf("%s/%s", pwd, conf.TranslationDir))
	if err != nil {
		log.Fatal(context.Background(), err)
	}

	if err := t.VerifyTranslations(); err != nil {
		log.Fatal(context.Background(), err)
	}

	return &translator{
		translator: t,
		log:        log,
	}
}

func parseLanguageId(conf Config) (locales.Translator, []locales.Translator, error) {
	var (
		locales []locales.Translator
	)

	localeIds := []string{conf.FallbackLanguageID}
	localeIds = append(localeIds, conf.SupportedLanguageIDs...)

	for _, v := range localeIds {
		switch v {
		case language.English:
			locales = append(locales, en.New())
		case language.Indonesian:
			locales = append(locales, id.New())
		default:
			return nil, nil, errors.NewWithCode(codes.CodeTranslatorError, fmt.Sprintf("unsupported languages ID %s", v))
		}
	}

	if len(locales) < 1 {
		return nil, nil, errors.NewWithCode(codes.CodeTranslatorError, "unsupported fallback language")
	}

	return locales[0], locales, nil
}

func (ut *translator) Translate(ctx context.Context, key interface{}, params ...string) (string, error) {
	if key == nil || key.(string) == "" {
		return "", nil
	}

	language := appcontext.GetAcceptLanguage(ctx)
	trans, found := ut.translator.GetTranslator(language)
	if !found {
		trans = ut.translator.GetFallback()
	}

	return trans.T(key, params...)

}
