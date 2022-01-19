package app

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const DefaultLocalePath = "./locale"

func NewBundle(path string) (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		bundle.MustLoadMessageFile(fmt.Sprintf("%s/%s", path, f.Name()))
	}

	return bundle, nil
}
