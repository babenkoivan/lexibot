package translation_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"lexibot/internal/testkit"
	"lexibot/internal/translation"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDeeplTranslator_Translate(t *testing.T) {
	server := newDeeplServer("foo", map[string]string{"mund": "mouth"})
	defer server.Close()

	t.Run("invalid auth key", func(t *testing.T) {
		translator := translation.NewDeeplTranslator(server.URL, "bar")
		transl, err := translator.Translate("bunt", "de", "en")

		assert.Error(t, err)
		assert.Len(t, transl, 0)
	})

	t.Run("translation is not found", func(t *testing.T) {
		translator := translation.NewDeeplTranslator(server.URL, "foo")
		transl, err := translator.Translate("bunt", "de", "en")

		assert.Error(t, err)
		assert.Len(t, transl, 0)
	})

	t.Run("translation is found", func(t *testing.T) {
		translator := translation.NewDeeplTranslator(server.URL, "foo")
		transl, err := translator.Translate("mund", "de", "en")

		assert.NoError(t, err)
		assert.Equal(t, "mouth", transl)
	})
}

func TestDBTranslator_Translate(t *testing.T) {
	langFrom, langTo := "de", "en"

	t.Run("translation is not found", func(t *testing.T) {
		translationStoreMock := testkit.MockTranslationStore(t)
		translator := translation.NewDBTranslator(translationStoreMock)
		transl, err := translator.Translate("bunt", langFrom, langTo)

		assert.Error(t, err)
		assert.Len(t, transl, 0)
	})

	t.Run("translation is found", func(t *testing.T) {
		text, expectedTransl := "mund", "mouth"

		translationStoreMock := testkit.MockTranslationStore(t)
		translationStoreMock.OnFirst(func(conds ...translation.TranslationQueryCond) *translation.Translation {
			testkit.AssertEqualTranslationQuery(t, []translation.TranslationQueryCond{
				translation.WithTextStrict(text),
				translation.WithLangFrom(langFrom),
				translation.WithLangTo(langTo),
				translation.WithManual(false),
			}, conds)

			return &translation.Translation{Translation: expectedTransl}
		})

		translator := translation.NewDBTranslator(translationStoreMock)
		actualTransl, err := translator.Translate(text, langFrom, langTo)

		assert.NoError(t, err)
		assert.Equal(t, expectedTransl, actualTransl)
	})
}

func TestCompositeTranslator_Translate(t *testing.T) {
	translatorMockA := testkit.MockTranslator(t)
	translatorMockB := testkit.MockTranslator(t)
	compositeTranslator := translation.NewCompositeTranslator(translatorMockA, translatorMockB)

	text, langFrom, langTo := "bunt", "de", "en"
	transl, err := compositeTranslator.Translate(text, langFrom, langTo)

	translatorMockA.AssertTranslated(text, langFrom, langTo)
	translatorMockB.AssertTranslated(text, langFrom, langTo)

	assert.Error(t, err)
	assert.Len(t, transl, 0)
}

func newDeeplServer(authKey string, dict map[string]string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryKey := r.URL.Query().Get("auth_key")
		if queryKey != authKey {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		params, err := url.ParseQuery(string(body))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		text := params.Get("text")
		transl, ok := dict[text]
		if !ok {
			transl = text
		}

		resp := fmt.Sprintf(`{"Translations":[{"Text":"%s"}]}`, transl)

		_, err = w.Write([]byte(resp))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	return httptest.NewServer(handler)
}
