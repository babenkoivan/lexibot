package translation

import (
	"regexp"
	"strings"
)

type TextSanitizer interface {
	Sanitize(text string) string
}

type SanitizeTextFunc func(text string) string

func (s SanitizeTextFunc) Sanitize(text string) string {
	return s(text)
}

func sanitizeGermanText(text string) string {
	text = removeArticles(text, []string{
		"der",
		"die",
		"das",
		"den",
		"dem",
		"des",
		"ein",
		"eine",
		"einen",
		"einem",
		"eines",
		"einer",
	})

	return strings.TrimSpace(text)
}

func NewGermanTextSanitizer() TextSanitizer {
	return SanitizeTextFunc(sanitizeGermanText)
}

func removeArticles(text string, articles []string) string {
	var pattern []string
	for _, a := range articles {
		pattern = append(pattern, `\b`+a+`\b`)
	}

	reg := regexp.MustCompile(`(?i)` + strings.Join(pattern, `|`))
	text = reg.ReplaceAllString(text, "")

	reg = regexp.MustCompile(` +`)
	return reg.ReplaceAllString(text, " ")
}
