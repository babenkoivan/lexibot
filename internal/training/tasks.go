package training

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"lexibot/internal/settings"
	"lexibot/internal/translation"
	"lexibot/internal/utils"
)

const (
	familiarWordScore = 5
	hintsLimit        = 3
)

type Task struct {
	Question *i18n.LocalizeConfig
	Answer   string
	Hints    []string
}

type taskGenerator struct {
	settingsStore    settings.SettingsStore
	translationStore translation.TranslationStore
	scoreStore       translation.ScoreStore
}

func (f *taskGenerator) Next(userID int) *Task {
	userSettings := f.settingsStore.FirstOrInit(userID)
	// todo exclude translation that were trained
	score := f.scoreStore.LowestNotTrained(userID, userSettings.LangDict)
	// todo exit with error if there is no score
	transl := f.translationStore.First(translation.WithID(score.TranslationID))

	var randTransl []*translation.Translation
	if f.includeHints(score.Score) {
		randTransl = append(randTransl, f.translationStore.Rand(
			translation.WithoutID(score.TranslationID),
			translation.WithUserID(score.UserID),
			translation.WithLangFrom(userSettings.LangDict),
			translation.WithLimit(hintsLimit-1),
		)...)

		if len(randTransl) > 0 {
			randTransl = append(randTransl, transl)

			utils.NewRand().Shuffle(len(randTransl), func(i, j int) {
				randTransl[i], randTransl[j] = randTransl[j], randTransl[i]
			})
		}
	}

	if f.translateToDictLang(score.Score) {
		task := &Task{
			Question: &i18n.LocalizeConfig{
				MessageID: "training.task",
				TemplateData: map[string]interface{}{
					"Text": transl.Translation,
				},
			},
			Answer: transl.Text,
		}

		for _, t := range randTransl {
			task.Hints = append(task.Hints, t.Text)
		}

		return task
	}

	task := &Task{
		Question: &i18n.LocalizeConfig{
			MessageID: "training.task",
			TemplateData: map[string]interface{}{
				"Text": transl.Text,
			},
		},
		Answer: transl.Translation,
	}

	for _, t := range randTransl {
		task.Hints = append(task.Hints, t.Translation)
	}

	return task
}

// when the word is familiar we ask to translate to the dict lang
// more often, otherwise we ask for the UI lang translation more
func (f *taskGenerator) translateToDictLang(score int) bool {
	r := utils.NewRand().Intn(100)

	if score > familiarWordScore {
		return r <= 70
	}

	return r <= 40
}

// when the word is familiar we rarely include hints,
// otherwise the hints are always included
func (f *taskGenerator) includeHints(score int) bool {
	if score > familiarWordScore {
		r := utils.NewRand().Intn(100)
		return r <= 20
	}

	return true
}

func NewTaskGenerator(
	settingsStore settings.SettingsStore,
	translationStore translation.TranslationStore,
	scoreStore translation.ScoreStore,
) *taskGenerator {
	return &taskGenerator{
		settingsStore,
		translationStore,
		scoreStore,
	}
}
