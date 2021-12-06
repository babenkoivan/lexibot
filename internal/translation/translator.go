package translation

import (
	"bytes"
	"encoding/json"
	"errors"
	"golang.org/x/text/language"
	"io"
	"lexibot/internal/config"
	"net/http"
	"net/url"
)

const (
	maxTranslations int    = 3
	translationErr  string = "failed to translate the text"
)

type Translator interface {
	Translate(from, to language.Tag, text string) ([]string, error)
}

type azureTranslator struct {
	endpoint string
	key      string
	region   string
}

func (a *azureTranslator) Translate(from, to language.Tag, text string) ([]string, error) {
	req, err := a.newRequest(from, to, text)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(translationErr)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data []azureResponseBody
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	var ts []string
	for _, t := range data[0].Translations {
		ts = append(ts, t.NormalizedTarget)
		if len(ts) == maxTranslations {
			break
		}
	}

	if len(ts) == 0 {
		return nil, errors.New(translationErr)
	}

	return ts, nil
}

type azureRequestBody struct {
	Text string
}

type azureResponseBody struct {
	Translations []struct {
		NormalizedTarget string
	}
}

func (a *azureTranslator) newRequest(from, to language.Tag, text string) (*http.Request, error) {
	u, err := url.Parse(a.endpoint)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("from", from.String())
	q.Add("to", to.String())
	u.RawQuery = q.Encode()

	data := []azureRequestBody{{text}}
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Ocp-Apim-Subscription-Key", a.key)
	req.Header.Add("Ocp-Apim-Subscription-Region", a.region)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func NewAzureTranslator(config config.Translator) Translator {
	return &azureTranslator{endpoint: config.Endpoint, key: config.Key, region: config.Region}
}
