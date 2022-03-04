package training

import (
	"lexibot/internal/settings"
	"lexibot/internal/translation"
	"lexibot/internal/utils"
)

const (
	FamiliarTermScore = 5
	HintsLimit        = 3
)

type TaskGenerator interface {
	Next(userID int) *Task
}

type translateTaskGenerator struct {
	settingsStore    settings.SettingsStore
	translationStore translation.TranslationStore
	scoreStore       translation.ScoreStore
	taskStore        TaskStore
}

func (g *translateTaskGenerator) Next(userID int) *Task {
	userSettings := g.settingsStore.FirstOrInit(userID)

	score := g.scoreStore.LowestNotTrained(userID, userSettings.LangDict)
	if score == nil {
		return nil
	}

	transl := g.translationStore.First(translation.WithID(score.TranslationID))

	var randTransl []*translation.Translation
	if g.includeHints(score.Score) {
		randTransl = append(randTransl, g.translationStore.Rand(
			translation.WithoutID(score.TranslationID),
			translation.WithUserID(score.UserID),
			translation.WithLangFrom(userSettings.LangDict),
			translation.WithLimit(HintsLimit-1),
		)...)

		if len(randTransl) > 0 {
			randTransl = append(randTransl, transl)

			utils.SourcedRand().Shuffle(len(randTransl), func(i, j int) {
				randTransl[i], randTransl[j] = randTransl[j], randTransl[i]
			})
		}
	}

	var task *Task
	if g.translateToDictLang(score.Score) {
		task = &Task{
			UserID:        userID,
			TranslationID: transl.ID,
			Question:      transl.Translation,
			Answer:        transl.Text,
		}

		for _, t := range randTransl {
			task.Hints = append(task.Hints, t.Text)
		}
	} else {
		task = &Task{
			UserID:        userID,
			TranslationID: transl.ID,
			Question:      transl.Text,
			Answer:        transl.Translation,
		}

		for _, t := range randTransl {
			task.Hints = append(task.Hints, t.Translation)
		}
	}

	return g.taskStore.Save(task)
}

// when the term is familiar we ask to translate to the dict lang
// more often, otherwise we ask for the UI lang translation more
func (g *translateTaskGenerator) translateToDictLang(score int) bool {
	r := utils.SourcedRand().Intn(100)

	if score >= FamiliarTermScore {
		return r <= 70
	}

	return r <= 40
}

// when the term is familiar we rarely include hints,
// otherwise the hints are always included
func (g *translateTaskGenerator) includeHints(score int) bool {
	if score >= FamiliarTermScore {
		r := utils.SourcedRand().Intn(100)
		return r <= 20
	}

	return true
}

func NewTaskGenerator(
	settingsStore settings.SettingsStore,
	translationStore translation.TranslationStore,
	scoreStore translation.ScoreStore,
	taskStore TaskStore,
) TaskGenerator {
	return &translateTaskGenerator{
		settingsStore,
		translationStore,
		scoreStore,
		taskStore,
	}
}
