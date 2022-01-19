package translation

type Translation struct {
	ID          uint64 `gorm:"primarykey"`
	Text        string `gorm:"uniqueIndex:idx_translation"`
	Translation string `gorm:"uniqueIndex:idx_translation"`
}
