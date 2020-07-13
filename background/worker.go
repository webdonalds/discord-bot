package background

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Worker struct {
	watchers []*WatcherContext
	mutex    sync.Mutex
	session  *discordgo.Session
}

type WatcherContext struct {
	watcher   Watcher
	channelID string
	mention   string
}

func NewWorker(sess *discordgo.Session) *Worker {
	return &Worker{watchers: []*WatcherContext{}, mutex: sync.Mutex{}, session: sess}
}

func (w *Worker) Start() {
	go func() {
		for range time.Tick(5 * time.Second) {
			w.Do()
		}
	}()
}

func (w *Worker) Do() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	newWatcher := []*WatcherContext{}
	defer func() { w.watchers = newWatcher }()

	for _, ctx := range w.watchers {
		if !ctx.watcher.ShouldRun() {
			newWatcher = append(newWatcher, ctx)
			continue
		}

		msg, watcher, err := ctx.watcher.Run()
		if err == nil {
			if msg != "" {
				_, _ = w.session.ChannelMessageSend(ctx.channelID, fmt.Sprintf("%s %s", ctx.mention, msg))
			}
			if watcher != nil {
				newWatcher = append(newWatcher, &WatcherContext{watcher: watcher, channelID: ctx.channelID, mention: ctx.mention})
			}
		}
	}
}

func (w *Worker) AddWatcher(watcher Watcher, channelID, mention string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.watchers = append(w.watchers, &WatcherContext{watcher: watcher, channelID: channelID, mention: mention})
}
