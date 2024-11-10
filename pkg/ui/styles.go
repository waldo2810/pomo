package ui

import "github.com/charmbracelet/lipgloss"

var (
	DocStyle    = lipgloss.NewStyle().Margin(10, 10)
	TitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF75B5"))
	TimerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B5")).Bold(true).Padding(0, 1)
	StatusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
)
