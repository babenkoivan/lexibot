package translation_test

import (
	"lexibot/internal/testkit"
	"lexibot/internal/translation"
	"testing"
)

func TestAutoDecrementScoreJob_Run(t *testing.T) {
	scoreStoreMock := testkit.MockScoreStore(t)

	job := translation.NewAutoDecrementScoreJob(scoreStoreMock)
	job.Run()

	scoreStoreMock.AssertAutoDecremented(translation.AutoDecrementAfter)
}
