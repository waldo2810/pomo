package model

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waldo2810/pomo/pkg/config"
	"github.com/waldo2810/pomo/pkg/ui"
)

type Model struct {
	config    config.Config
	mode      Mode
	timer     timer.Model
	list      list.Model
	selected  Choice
	asking    bool
	message   string
	completed int   // Number of completed Pomodoros
	err       error // For error handling
}

func NewModel() Model {
	items := []list.Item{
		Choice{id: 1, title: "Start default", desc: "25 mins focus and 5 mins break"},
		//item{id: 2, title: "Start custom", desc: "Set your own durations"},
	}

	choices := list.New(items, list.NewDefaultDelegate(), 0, 0)
	choices.Title = "🍅 Welcome to Pomo"

	return Model{
		config: config.Default,
		mode:   ModeIdle,
		list:   choices,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Set window size
	case tea.WindowSizeMsg:
		h, v := ui.DocStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	// On every tick, update the timer
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	// When the timer finishes, handle the timeout
	case timer.TimeoutMsg:
		return m.handleTimeout()

	// Handle key presses
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) handleTimeout() (tea.Model, tea.Cmd) {
	var notificationMsg string
	m.asking = true

	switch m.mode {
	case ModeFocus:
		m.completed++
		if m.completed%m.config.Cycle == 0 {
			notificationMsg = "Time for a long break!"
			m.message = "Take a long break? (y/n)"
		} else {
			notificationMsg = "Time for a break!"
			m.message = "Take a break? (y/n)"
		}
	case ModeBreak, ModeLongBreak:
		notificationMsg = "Break's over! Ready to focus?"
		m.message = "Start focusing? (y/n)"
	}

	go m.sendNotification(notificationMsg)
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "ctrl+c":
		return m, tea.Quit

	case "y":
		if m.asking {
			m.asking = false
			return m.startNextPhase()
		}

	case "n":
		if m.asking {
			m.asking = false
			m.message = ""
			m.mode = ModeIdle
			return m, nil
		}

	case "enter":
		if !m.asking {
			if i, ok := m.list.SelectedItem().(Choice); ok {
				switch i.id {
				case 1:
					m.config = config.Default
					return m.startFocusPhase()
				case 2:
					// Handle custom configuration
					return m, nil
				case 3:
					// Show statistics
					return m, nil
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) startNextPhase() (tea.Model, tea.Cmd) {
	switch m.mode {
	case ModeFocus:
		if m.completed%m.config.Cycle == 0 {
			m.mode = ModeLongBreak
			m.timer = timer.NewWithInterval(m.config.LongBreakDuration, time.Second)
		} else {
			m.mode = ModeBreak
			m.timer = timer.NewWithInterval(m.config.BreakDuration, time.Second)
		}
	case ModeBreak, ModeLongBreak:
		return m.startFocusPhase()
	}
	return m, m.timer.Init()
}

func (m Model) startFocusPhase() (tea.Model, tea.Cmd) {
	m.mode = ModeFocus
	m.timer = timer.NewWithInterval(m.config.FocusDuration, time.Second)
	return m, m.timer.Init()
}

func (m Model) sendNotification(msg string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("notify-send", "Pomodoro Timer", msg).Run()
		exec.Command("spd-say", msg).Run()
	}
}

func (m Model) View() string {
	if m.err != nil {
		return ui.DocStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	if m.asking {
		return ui.DocStyle.Render(ui.TitleStyle.Render(m.message))
	}

	if m.mode != ModeIdle {
		return ui.DocStyle.Render(fmt.Sprintf(
			"%s\n\n%s\n%s",
			ui.TitleStyle.Render("🍅 Pomodoro Timer"),
			ui.StatusStyle.Render(fmt.Sprintf("Mode: %s (Completed: %d)", m.mode, m.completed)),
			ui.TimerStyle.Render(m.timer.View()),
		))
	}

	return ui.DocStyle.Render(m.list.View())
}
