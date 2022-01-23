package locale

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type UserLocale interface {
	MakeLocalizer(userID int) *i18n.Localizer
}
