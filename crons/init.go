package crons

type Cron interface {
	Pattern() string
	ChannelID() string
	Execute() string
}
