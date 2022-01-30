package translation

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	requestErr = "failed to retrieve translations"
	parseErr   = "failed to parse the response"
)

type Translator interface {
	Translate(text, langFrom, langTo string) (string, error)
}

type deeplTranslator struct {
	endpoint string
	key      string
}

func (t *deeplTranslator) Translate(text, langFrom, langTo string) (string, error) {
	req, err := t.newRequest(langFrom, langTo, text)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(requestErr)
	}

	return t.parseResponse(resp)
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
		return "", errors.New(parseErr)
	}

	return content.Translations[0].Text, nil
}

type cachingTranslator struct {
	translationStore TranslationStore
	deeplTranslator  *deeplTranslator
}

func (t *cachingTranslator) Translate(text, langFrom, langTo string) (string, error) {
	cached := t.translationStore.GetAuto(text, langFrom, langTo)
	if cached != nil {
		return cached.Translation, nil
	}

	translation, err := t.deeplTranslator.Translate(text, langFrom, langTo)
	if err != nil {
		return "", err
	}

	t.translationStore.Save(&Translation{
		Text:        text,
		LangFrom:    langFrom,
		LangTo:      langTo,
		Translation: translation,
	})

	return translation, nil
}

func NewTranslator(endpoint string, key string, translationStore TranslationStore) Translator {
	deeplTranslator := &deeplTranslator{endpoint, key}
	return &cachingTranslator{translationStore, deeplTranslator}
}
