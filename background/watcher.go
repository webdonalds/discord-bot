package background

import (
	"time"
)

type Watcher interface {
	ShouldRun() bool
	Run() (string, Watcher, error)
}

type TimerWatcher struct {
	runAt   *time.Time
	message string
}

func NewTimerWatcher(runAt *time.Time, msg string) Watcher {
	return &TimerWatcher{runAt: runAt, message: msg}
}

func (w *TimerWatcher) ShouldRun() bool {
	return time.Now().After(*w.runAt)
}

func (w *TimerWatcher) Run() (string, Watcher, error) {
	return w.message, nil, nil
}
