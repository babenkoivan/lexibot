package translation

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	deeplTimeout     = 3 * time.Second
	deeplRequestErr  = "failed to retrieve translations"
	deeplParseErr    = "failed to parse the response"
	deeplNotFoundErr = "translations not found"
)

type Translator interface {
	Translate(text, langFrom, langTo string) (string, error)
}

type deeplTranslator struct {
	client   *http.Client
	endpoint string
	key      string
}

func (t *deeplTranslator) Translate(text, langFrom, langTo string) (string, error) {
	req, err := t.newRequest(langFrom, langTo, text)
	if err != nil {
		return "", err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(deeplRequestErr)
	}

	translation, err := t.parseResponse(resp)
	if err != nil {
		return "", err
	}

	// when translations are not found, the input text will be the same as the translation
	if strings.EqualFold(strings.TrimSpace(translation), strings.TrimSpace(text)) {
		return "", errors.New(deeplNotFoundErr)
	}

	return translation, nil
}

func (t deeplTranslator) newRequest(langFrom, langTo, text string) (*http.Request, error) {
	u, err := url.Parse(t.endpoint)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("auth_key", t.key)
	u.RawQuery = q.Encode()

	v := url.Values{}
	v.Add("source_lang", langFrom)
	v.Add("target_lang", langTo)
	v.Add("text", text)
	b := strings.NewReader(v.Encode())

	req, err := http.NewRequest("POST", u.String(), b)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (t *deeplTranslator) parseResponse(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var content struct {
		Translations []struct {
			Text string
		}
	}

	if err := json.Unmarshal(body, &content); err != nil {
		return "", err
	}

	if len(content.Translations) == 0 {
		return "", errors.New(deeplParseErr)
	}

	return content.Translations[0].Text, nil
}

type dbTranslator struct {
	translationStore TranslationStore
}

func (t *dbTranslator) Translate(text, langFrom, langTo string) (string, error) {
	translation := t.translationStore.First(
		WhereText(text),
		WhereLangFrom(langFrom),
		WhereLangTo(langTo),
		WhereManual(false),
	)

	if translation != nil {
		return translation.Translation, nil
	}

	return "", nil
}

type compositeTranslator struct {
	translators []Translator
}

func (t *compositeTranslator) Translate(text, langFrom, langTo string) (string, error) {
	var lastErr error

	for _, translator := range t.translators {
		translation, err := translator.Translate(text, langFrom, langTo)

		if translation != "" {
			return translation, err
		}

		lastErr = err
	}

	return "", lastErr
}

func NewTranslator(endpoint string, key string, translationStore TranslationStore) Translator {
	client := &http.Client{Timeout: deeplTimeout}

	return &compositeTranslator{[]Translator{
		&dbTranslator{translationStore},
		&deeplTranslator{client, endpoint, key},
	}}
}
