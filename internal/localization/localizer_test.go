package localization_test

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"lexibot/internal/localization"
	"testing"
)

func TestLocalizer_MatchMessage(t *testing.T) {
	msgYes := &i18n.Message{ID: "yes", Other: "да"}
	msgNo := &i18n.Message{ID: "no", Other: "нет"}

	bundle := i18n.NewBundle(language.English)
	err := bundle.AddMessages(language.English, msgYes, msgNo)
	require.NoError(t, err)

	localizer := localization.NewLocalizer(bundle, language.English.String())
	matchID := localizer.MatchMessage(msgYes.Other, []string{msgNo.ID, msgYes.ID})

	assert.Equal(t, msgYes.ID, matchID)
}
