package config

import "time"

type WatcherConfig struct {
	SourceURL string
	Interval  time.Duration
	Jitter    time.Duration
	LockKey   string
	WMKey     string
}
