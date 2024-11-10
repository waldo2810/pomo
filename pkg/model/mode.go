package model

type Mode int

const (
	ModeIdle Mode = iota
	ModeFocus
	ModeBreak
	ModeLongBreak
)

func (m Mode) String() string {
	return [...]string{"Idle", "Focus", "Break", "Long Break"}[m]
}
