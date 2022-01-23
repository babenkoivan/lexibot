package app

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"lexibot/internal/config"
	"lexibot/internal/locale"
)

const DefaultLocalePath = "./locale"

type userLocale struct {
	bundle      *i18n.Bundle
	configStore config.ConfigStore
}

func (l *userLocale) MakeLocalizer(userID int) *i18n.Localizer {
	config := l.configStore.Get(userID)
	lang := language.English.String()

	if config != nil {
		lang = config.LangUI
	}

	return i18n.NewLocalizer(l.bundle, lang)
}

func NewUserLocale(localePath string, configStore config.ConfigStore) (locale.UserLocale, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	files, err := ioutil.ReadDir(localePath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		bundle.MustLoadMessageFile(fmt.Sprintf("%s/%s", localePath, f.Name()))
	}

	return &userLocale{bundle, configStore}, nil
}
