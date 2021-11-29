package translations_test

import (
	"cloud.google.com/go/translate"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
	"lexibot/internal/pkg/translations"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGoogleTranslatorTranslate(t *testing.T) {
	server := newGoogleServer(map[string]string{"bunt": "colorful"})
	defer server.Close()

	ctx := context.Background()
	client := newGoogleClient(t, ctx, server.URL)
	defer client.Close()

	translator := translations.NewGoogleTranslator(client)

	t.Run("returns translations when found", func(t *testing.T) {
		want := []string{"colorful"}
		got, err := translator.Translate(ctx, language.German, language.English, "bunt")

		require.NoError(t, err)
		assert.Equal(t, got, want)
	})

	t.Run("returns an error when translations not found", func(t *testing.T) {
		_, err := translator.Translate(ctx, language.German, language.English, "weg")

		assert.Error(t, err)
	})
}

func newGoogleServer(ts map[string]string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		t, ok := ts[q]

		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "{\"data\":{\"translations\":[{\"translatedText\":%q}]}}", t)
	})

	return httptest.NewServer(handler)
}

func newGoogleClient(t *testing.T, ctx context.Context, url string) *translate.Client {
	client, err := translate.NewClient(ctx, option.WithEndpoint(url), option.WithoutAuthentication())

	if err != nil {
		t.Fatal(err.Error())
	}

	return client
}
