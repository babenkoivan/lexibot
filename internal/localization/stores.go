package localization

type LocaleStore interface {
	Locale(userID int) string
}
