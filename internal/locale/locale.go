package locale

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const DefaultPath = "./locale"

type Locale interface {
	MakeLocalizer(userID int) *i18n.Localizer
}

type locale struct {
	bundle      *i18n.Bundle
	localeStore LocaleStore
}

func (l *locale) MakeLocalizer(userID int) *i18n.Localizer {
	lang := l.localeStore.GetLocale(userID)

	if lang == "" {
		lang = language.English.String()
	}

	return i18n.NewLocalizer(l.bundle, lang)
}

func NewLocale(localePath string, localeStore LocaleStore) (Locale, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	files, err := ioutil.ReadDir(localePath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		bundle.MustLoadMessageFile(fmt.Sprintf("%s/%s", localePath, f.Name()))
	}

	return &locale{bundle, localeStore}, nil
}
