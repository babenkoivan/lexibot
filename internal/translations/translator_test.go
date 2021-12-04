package translations_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"io"
	"lexibot/internal/configs"
	"lexibot/internal/translations"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAzureTranslator(t *testing.T) {
	server := newAzureServer(map[string][]string{"weit": {"far", "widely", "well", "long"}})
	defer server.Close()

	config := configs.Azure{Endpoint: server.URL}
	translator := translations.NewAzureTranslator(config)

	t.Run("returns translations on success", func(t *testing.T) {
		want := []string{"far", "widely", "well", "long"}
		got, err := translator.Translate(language.German, language.English, "weit")

		require.NoError(t, err)
		assert.Len(t, got, translations.MaxTranslations)
		assert.Equal(t, got, want)
	})

	t.Run("returns an error on error", func(t *testing.T) {
		_, err := translator.Translate(language.German, language.English, "bunt")

		assert.Error(t, err, translations.TranslationErr)
	})
}

func newAzureServer(dict map[string][]string) *httptest.Server {
	badRequest := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusBadRequest)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			badRequest(w)
			return
		}

		var input []translations.AzureRequestBody
		if err := json.Unmarshal(b, &input); err != nil {
			badRequest(w)
			return
		}

		text := input[0].Text
		ts, ok := dict[text]

		if !ok {
			badRequest(w)
			return
		}

		output := []translations.AzureResponseBody{{}}
		for _, t := range ts {
			nt := struct{ NormalizedTarget string }{NormalizedTarget: t}
			output[0].Translations = append(output[0].Translations, nt)
		}

		b, err = json.Marshal(output)
		if err != nil {
			badRequest(w)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			badRequest(w)
			return
		}
	})

	return httptest.NewServer(handler)
}
