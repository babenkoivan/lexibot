package localization

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

const DefaultPath = "./locales"

type Localizer interface {
	MustLocalize(config *i18n.LocalizeConfig) string
	MatchMessage(localizedText string, messageIDs []string) (messageID string)
}

type LocalizerFactory interface {
	New(userID int) Localizer
}

type localizer struct {
	*i18n.Localizer
}

func (l *localizer) MatchMessage(localizedText string, messageIDs []string) (messageID string) {
	localizedText = strings.TrimSpace(localizedText)

	for _, ID := range messageIDs {
		localizedMessage := l.MustLocalize(&i18n.LocalizeConfig{MessageID: ID})

		if strings.EqualFold(localizedText, localizedMessage) {
			messageID = ID
			break
		}
	}

	return
}

func NewLocalizer(bundle *i18n.Bundle, lang string) Localizer {
	return &localizer{i18n.NewLocalizer(bundle, lang)}
}

type localizerFactory struct {
	bundle      *i18n.Bundle
	localeStore LocaleStore
}

func (l *localizerFactory) New(userID int) Localizer {
	lang := l.localeStore.Locale(userID)

	if lang == "" {
		lang = language.English.String()
	}

	return NewLocalizer(l.bundle, lang)
}

func NewLocalizerFactory(localesPath string, localeStore LocaleStore) (LocalizerFactory, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	files, err := ioutil.ReadDir(localesPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		bundle.MustLoadMessageFile(fmt.Sprintf("%s/%s", localesPath, f.Name()))
	}

	return &localizerFactory{bundle, localeStore}, nil
}
