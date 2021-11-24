package translation_test

import (
	"cloud.google.com/go/translate"
	"context"
	"fmt"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
	"lexibot/internal/translation"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGoogleTranslator(t *testing.T) {
	server := newGoogleServer(map[string]string{"bunt": "colorful"})
	defer server.Close()

	ctx := context.Background()
	client := newGoogleClient(t, ctx, server.URL)
	defer client.Close()

	translator := translation.NewGoogleTranslator(client)

	t.Run("returns translations when found", func(t *testing.T) {
		want := []string{"colorful"}
		got, _ := translator.Translate(ctx, language.German, language.English, "bunt")

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v translations got %v", want, got)
		}
	})

	t.Run("returns an error when translations not found", func(t *testing.T) {
		_, err := translator.Translate(ctx, language.German, language.English, "weg")

		if err == nil {
			t.Error("Expected an error got nil")
		}
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
