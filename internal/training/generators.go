package training

import (
	"lexibot/internal/settings"
	"lexibot/internal/translation"
	"lexibot/internal/utils"
)

const (
	FamiliarWordScore = 5
	HintsLimit        = 3
)

type TaskGenerator interface {
	Next(userID int) *Task
}

type translateTaskGenerator struct {
	settingsStore    settings.SettingsStore
	translationStore translation.TranslationStore
	taskStore        TaskStore
}

func (g *translateTaskGenerator) Next(userID int) *Task {
	userSettings := g.settingsStore.FirstOrInit(userID)
	trainedTranslIDs := g.taskStore.TranslationIDs(userID)

	transl := g.translationStore.First(
		translation.WithoutIDs(trainedTranslIDs),
		translation.WithUserID(userID),
		translation.WithLowestScore(),
	)

	if transl == nil {
		return nil
	}

	var randTransl []*translation.Translation
	if g.includeHints(transl.Score) {
		randTransl = append(randTransl, g.translationStore.Rand(
			translation.WithoutIDs([]int{transl.ID}),
			translation.WithUserID(transl.UserID),
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
	if g.translateToDictLang(transl.Score) {
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

// when the word is familiar we ask to translate to the dict lang
// more often, otherwise we ask for the UI lang translation more
func (g *translateTaskGenerator) translateToDictLang(score int) bool {
	r := utils.SourcedRand().Intn(100)

	if score >= FamiliarWordScore {
		return r <= 70
	}

	return r <= 40
}

// when the word is familiar we rarely include hints,
// otherwise the hints are always included
func (g *translateTaskGenerator) includeHints(score int) bool {
	if score >= FamiliarWordScore {
		r := utils.SourcedRand().Intn(100)
		return r <= 20
	}

	return true
}

func NewTaskGenerator(
	settingsStore settings.SettingsStore,
	translationStore translation.TranslationStore,
	taskStore TaskStore,
) TaskGenerator {
	return &translateTaskGenerator{
		settingsStore,
		translationStore,
		taskStore,
	}
}
