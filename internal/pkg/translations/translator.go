package translations

import (
	"cloud.google.com/go/translate"
	"context"
	"errors"
	"golang.org/x/text/language"
)

const (
	NoTranslationsErr = "no translations are found"
)

type Translator interface {
	Translate(ctx context.Context, from, to language.Tag, text string) ([]string, error)
}

type googleTranslator struct {
	client *translate.Client
}

func (g *googleTranslator) Translate(ctx context.Context, from, to language.Tag, text string) ([]string, error) {
	opts := &translate.Options{Source: from, Format: translate.Text}
	resp, err := g.client.Translate(ctx, []string{text}, to, opts)

	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New(NoTranslationsErr)
	}

	var ts []string

	for _, t := range resp {
		ts = append(ts, t.Text)
	}

	return ts, nil
}

func NewGoogleTranslator(client *translate.Client) *googleTranslator {
	return &googleTranslator{client: client}
}
