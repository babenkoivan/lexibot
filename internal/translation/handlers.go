package translation

const (
	OnCancelTranslation string = "cancel_translation"
	OnSaveTranslation   string = "save_translation"
	OnDeleteTranslation string = "delete_translation"

	textSeparator string = "-"
)

//type translateTextHandler struct {
//	translator Translator
//	store      Store
//}
//
//func (h *translateTextHandler) Handle(b bot.Bot, m *telebot.Message) {
//	input := strings.Split(m.Text, textSeparator)
//	text := strings.TrimSpace(input[0])
//
//	if h.store.Exists(text) {
//		b.Send(m.Sender, bot.NewTextMessage("A translation already exists"))
//		return
//	}
//
//	// save the translation if provided
//	if len(input) > 1 {
//		translation := h.store.Create(text, strings.TrimSpace(input[1]))
//		b.Send(m.Sender, &savedTranslationMessage{translation})
//		return
//	}
//
//	// otherwise, translate the given text
//	// todo take from the user app
//	from := language.German
//	to := language.English
//
//	res, err := h.translator.Translate(from, to, text)
//	if err != nil {
//		b.Send(m.Sender, bot.NewErrorMessage(err))
//		return
//	}
//
//	// suggest choosing a translation if many
//	if len(res) > 1 {
//		b.Send(m.Sender, &chooseTranslationMessage{text, res})
//		return
//	}
//
//	// otherwise, save the result
//	translation := h.store.Create(text, res[0])
//	b.Send(m.Sender, &savedTranslationMessage{translation})
//}
//
//func NewTranslateTextHandler(translator Translator, store Store) *translateTextHandler {
//	return &translateTextHandler{translator: translator, store: store}
//}
//
//func cancelTranslationHandler(b bot.Bot, c *telebot.Callback) {
//	b.Edit(bot.ExtractMessageSig(c.Message), bot.NewTextMessage("The translation has been canceled"))
//}
//
//func NewCancelTranslationHandler() bot.CallbackHandler {
//	return bot.CallbackHandlerFunc(cancelTranslationHandler)
//}
//
//type saveTranslationHandler struct {
//	store Store
//}
//
//func (h *saveTranslationHandler) Handle(b bot.Bot, c *telebot.Callback) {
//	data := strings.Split(c.Data, "|")
//	translation := h.store.Create(data[0], data[1])
//	b.Edit(bot.ExtractMessageSig(c.Message), &savedTranslationMessage{translation})
//}
//
//func NewSaveTranslationHandler(store Store) *saveTranslationHandler {
//	return &saveTranslationHandler{store: store}
//}
//
//type deleteTranslationHandler struct {
//	store Store
//}
//
//func (h *deleteTranslationHandler) Handle(b bot.Bot, c *telebot.Callback) {
//	ID, _ := strconv.ParseUint(c.Data, 10, 64)
//	h.store.Delete(ID)
//	b.Edit(bot.ExtractMessageSig(c.Message), bot.NewTextMessage("The translation has been deleted"))
//}
//
//func NewDeleteTranslationHandler(store Store) *deleteTranslationHandler {
//	return &deleteTranslationHandler{store: store}
//}
