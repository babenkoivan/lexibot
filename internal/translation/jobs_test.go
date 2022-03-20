package translation_test

import (
	"lexibot/internal/testkit"
	"lexibot/internal/translation"
	"testing"
)

func TestAutoDecrementScoreJob_Run(t *testing.T) {
	translationStoreMock := testkit.MockTranslationStore(t)

	job := translation.NewAutoDecrementScoreJob(translationStoreMock)
	job.Run()

	translationStoreMock.AssertAutoDecremented(translation.AutoDecrementAfter)
}
