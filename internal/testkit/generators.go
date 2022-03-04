package testkit

import "lexibot/internal/training"

type taskGeneratorMock struct {
	onNext func(userID int) *training.Task
}

func (m *taskGeneratorMock) OnNext(callback func(userID int) *training.Task) {
	m.onNext = callback
}

func (m *taskGeneratorMock) Next(userID int) *training.Task {
	if m.onNext == nil {
		return nil
	}

	return m.onNext(userID)
}

func MockTaskGenerator() *taskGeneratorMock {
	return &taskGeneratorMock{}
}
