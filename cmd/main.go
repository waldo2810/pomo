package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(10, 10)

type item struct {
	id    int
	title string
	desc  string
}

// Getters
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type pomo struct {
	mode     string
	timer    timer.Model
	list     list.Model
	selected item
	asking   bool   // Tracks if we are currently asking for user input
	message  string // Message to show when prompting user for input
}

func (state pomo) Init() tea.Cmd {
	return nil
}

func (state pomo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Set the size of the list
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		state.list.SetSize(msg.Width-h, msg.Height-v)

	// On every tick, update the timer
	case timer.TickMsg:
		var cmd tea.Cmd
		state.timer, cmd = state.timer.Update(msg)
		return state, cmd

	// When the timer is done, notify the user and prompt for the next step
	case timer.TimeoutMsg:
		var notificationMessage string
		if state.mode == "focus" {
			notificationMessage = "Time to take a break!"
			state.message = "Take a break? (y/n)"
		} else if state.mode == "break" {
			notificationMessage = "Get back to work!"
			state.message = "Get back to work? (y/n)"
		}
		state.asking = true

		// Display system notifications
		if os := runtime.GOOS; os == "linux" {
			exec.Command("spd-say", notificationMessage).Run()
			exec.Command("notify-send", notificationMessage).Run()
		}
		return state, nil

	// Handle key presses
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return state, tea.Quit

		// Handle "y" and "n" answers when in asking mode
		case "y":
			if state.asking {
				state.asking = false
				if state.mode == "focus" {
					state.mode = "break"
					state.timer = timer.NewWithInterval(5*time.Second, time.Second)
				} else {
					state.mode = "focus"
					state.timer = timer.NewWithInterval(25*time.Minute, time.Second)
				}
				return state, state.timer.Init()
			}

		case "n":
			if state.asking {
				state.asking = false
				state.message = ""
				// Return to idle mode or main list view
				return state, nil
			}

		case "enter":
			if !state.asking {
				found, exists := state.list.SelectedItem().(item)
				if exists && found.id == 1 {
					if state.mode == "focus" {
						state.mode = "break"
						state.timer = timer.NewWithInterval(5*time.Second, time.Second)
					} else {
						state.mode = "focus"
						state.timer = timer.NewWithInterval(25*time.Minute, time.Second)
					}
					state.selected = found
					return state, state.timer.Init()
				}
			}
		}
	}

	var cmd tea.Cmd
	state.list, cmd = state.list.Update(msg)
	return state, cmd
}

func (state pomo) View() string {
	if state.asking {
		return docStyle.Render("\n" + state.message)
	}

	if state.selected.id == 1 {
		ui := "\nUsing default configuration\n"
		ui = ui + "\nYou are now in " + state.mode + " mode"
		ui = ui + "\nBreak in " + state.timer.View()
		return docStyle.Render(ui)
	}

	return docStyle.Render(state.list.View())
}

func main() {
	items := []list.Item{
		item{id: 1, title: "Start default", desc: "25 mins focus and 5 mins break"},
		item{id: 2, title: "Start custom", desc: "Set your own durations"},
	}

	state := pomo{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
		mode: "focus",
	}
	state.list.Title = "Welcome to Pomo"

	app := tea.NewProgram(state, tea.WithAltScreen())

	if _, err := app.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
