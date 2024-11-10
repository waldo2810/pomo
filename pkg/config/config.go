package config

import "time"

type Config struct {
	FocusDuration     time.Duration
	BreakDuration     time.Duration
	LongBreakDuration time.Duration
	Cycle             int // Number of pomodoros before long break
}

var Default = Config{
	FocusDuration:     25 * time.Minute,
	BreakDuration:     5 * time.Minute,
	LongBreakDuration: 15 * time.Minute,
	Cycle:             4,
}
