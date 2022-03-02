package testkit

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"lexibot/internal/localization"
	"testing"
)

const localesPath = "../../locales"

func MockLocalizerFactory(t *testing.T, locale language.Tag) localization.LocalizerFactory {
	settingsStoreMock := MockSettingsStore(t)

	settingsStoreMock.OnLocale(func(userID int) string {
		return locale.String()
	})

	localizerFactory, err := localization.NewLocalizerFactory(localesPath, settingsStoreMock)
	require.NoError(t, err)

	return localizerFactory
}
