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

	t.Run("translation not found", func(t *testing.T) {
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
	storeMock := testkit.MockTranslationStore()
	translator := translation.NewDBTranslator(storeMock)

	t.Run("translation is not found", func(t *testing.T) {
		storeMock.OnFirst(func(conds ...translation.TranslationQueryCond) *translation.Translation {
			return nil
		})

		transl, err := translator.Translate("bunt", "de", "en")

		assert.Error(t, err)
		assert.Len(t, transl, 0)
	})

	t.Run("translation is found", func(t *testing.T) {
		storeMock.OnFirst(func(conds ...translation.TranslationQueryCond) *translation.Translation {
			return &translation.Translation{Translation: "mouth"}
		})

		transl, err := translator.Translate("mund", "de", "en")

		assert.NoError(t, err)
		assert.Equal(t, "mouth", transl)
	})
}

func TestCompositeTranslator_Translate(t *testing.T) {
	translatorSpyA := testkit.NewTranslatorSpy(t)
	translatorSpyB := testkit.NewTranslatorSpy(t)
	compositeTranslator := translation.NewCompositeTranslator(translatorSpyA, translatorSpyB)

	text, langFrom, langTo := "bunt", "de", "en"
	transl, err := compositeTranslator.Translate(text, langFrom, langTo)

	translatorSpyA.AssertTranslated(text, langFrom, langTo)
	translatorSpyB.AssertTranslated(text, langFrom, langTo)

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
