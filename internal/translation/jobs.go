package translation

import (
	"time"
)

const AutoDecrementAfter = 24 * 7 * time.Hour

type autoDecrementScoreJob struct {
	translationStore TranslationStore
}

func (j *autoDecrementScoreJob) Run() {
	j.translationStore.AutoDecrementScore(AutoDecrementAfter)
}

func NewAutoDecrementScoreJob(translationStore TranslationStore) *autoDecrementScoreJob {
	return &autoDecrementScoreJob{translationStore}
}
