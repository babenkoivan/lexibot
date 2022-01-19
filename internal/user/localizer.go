package user

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func NewLocalizer(bundle *i18n.Bundle, userId int) *i18n.Localizer {
	// todo get UI language from settings
	return i18n.NewLocalizer(bundle, language.English.String())
}
