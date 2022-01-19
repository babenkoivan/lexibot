package bot

import (
	"fmt"
)

type Message interface {
	Type() string
	Render() (text string, options []interface{})
}

type ErrorMessage struct {
	Err error
}

func (m *ErrorMessage) Type() string {
	return "error"
}

func (m *ErrorMessage) Render() (text string, options []interface{}) {
	text = fmt.Sprintf("❗️ %s", m.Err)
	return
}

type PlainTextMessage struct {
	Text string
}

func (m *PlainTextMessage) Type() string {
	return "plain_text"
}

func (m *PlainTextMessage) Render() (text string, options []interface{}) {
	text = m.Text
	return
}
