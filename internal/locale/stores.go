package locale

type LocaleStore interface {
	GetLocale(userID int) string
}
