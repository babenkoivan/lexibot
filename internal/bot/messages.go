package bot

import (
	"fmt"
)

type Message interface {
	Id() string
	Render() (text string, options []interface{})
}

type ErrorMessage struct {
	err error
}

func (m *ErrorMessage) Id() string {
	return "error"
}

func (m *ErrorMessage) Render() (text string, options []interface{}) {
	text = fmt.Sprintf("❗️ %s", m.err)
	return
}

type PlainTextMessage struct {
	text string
}

func (m *PlainTextMessage) Id() string {
	return "plain_text"
}

func (m *PlainTextMessage) Render() (text string, options []interface{}) {
	text = m.text
	return
}
