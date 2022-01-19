package translation

type Translation struct {
	ID          uint64 `gorm:"primaryKey"`
	Text        string `gorm:"uniqueIndex:idx_translation"`
	Translation string `gorm:"uniqueIndex:idx_translation"`
}
