package translation

import (
	"time"
)

const AutoDecrementAfter = 24 * 7 * time.Hour

type autoDecrementScoreJob struct {
	scoreStore ScoreStore
}

func (j *autoDecrementScoreJob) Run() {
	j.scoreStore.AutoDecrement(AutoDecrementAfter)
}

func NewAutoDecrementScoreJob(scoreStore ScoreStore) *autoDecrementScoreJob {
	return &autoDecrementScoreJob{scoreStore}
}
