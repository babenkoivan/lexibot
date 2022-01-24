package config

import "golang.org/x/text/language"

func SupportedLangUI() []string {
	return []string{
		language.English.String(),
		language.Russian.String(),
	}
}

func SupportedLangDict() []string {
	return []string{
		language.German.String(),
		language.English.String(),
		language.Russian.String(),
	}
}
