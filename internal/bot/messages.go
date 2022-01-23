package bot

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Message interface {
	Type() string
	Render(localizer *i18n.Localizer) (text string, options []interface{})
}
