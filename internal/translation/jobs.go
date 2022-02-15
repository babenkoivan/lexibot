package translation

import (
	"time"
)

const autoDecrementAfter = 24 * 7 * time.Hour

type autoDecrementScoreJob struct {
	scoreStore ScoreStore
}

func (j *autoDecrementScoreJob) Run() {
	j.scoreStore.AutoDecrement(autoDecrementAfter)
}

func NewAutoDecrementScoreJob(scoreStore ScoreStore) *autoDecrementScoreJob {
	return &autoDecrementScoreJob{scoreStore}
}
