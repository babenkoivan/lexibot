package locale

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const DefaultPath = "./locale"

type UserLocale interface {
	MakeLocalizer(userID int) *i18n.Localizer
}

type UserLocaleStore interface {
	GetLocale(userID int) string
}

type userLocale struct {
	bundle      *i18n.Bundle
	localeStore UserLocaleStore
}

func (l *userLocale) MakeLocalizer(userID int) *i18n.Localizer {
	lang := l.localeStore.GetLocale(userID)

	if lang == "" {
		lang = language.English.String()
	}

	return i18n.NewLocalizer(l.bundle, lang)
}

func NewUserLocale(localePath string, localeStore UserLocaleStore) (UserLocale, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	files, err := ioutil.ReadDir(localePath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		bundle.MustLoadMessageFile(fmt.Sprintf("%s/%s", localePath, f.Name()))
	}

	return &userLocale{bundle, localeStore}, nil
}
