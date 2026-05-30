package ui

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor    = lipgloss.Color("#AF87FF")
	accentColor     = lipgloss.Color("#06B6D4")
	successColor    = lipgloss.Color("#10B981")
	errorColor      = lipgloss.Color("#EF4444")
	warningColor    = lipgloss.Color("#F59E0B")
	mutedColor      = lipgloss.Color("#6B7280")
	lightMutedColor = lipgloss.Color("#9CA3AF")
	textColor       = lipgloss.Color("#F9FAFB")
	borderColor     = lipgloss.Color("#4B5563")

	headerStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 2).
			MarginBottom(1)

	sectionHeaderStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true).
				Background(lipgloss.Color("#1E293B")).
				Padding(0, 2).
				MarginTop(1).
				MarginBottom(1)

	tableHeaderStyle = lipgloss.NewStyle().
				Foreground(lightMutedColor).
				Bold(true).
				Padding(0, 1)

	additionStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	deletionStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	mutedNumberStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	filePathStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Bold(false)

	summaryBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	summaryPositiveStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	summaryNegativeStyle = lipgloss.NewStyle().
				Foreground(errorColor).
				Bold(true)

	summaryNeutralStyle = lipgloss.NewStyle().
				Foreground(warningColor).
				Bold(true)

	summaryLabelStyle = lipgloss.NewStyle().
				Foreground(lightMutedColor).
				Padding(0, 1, 0, 0)

	summaryValueStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Bold(true)

	separatorStyle = lipgloss.NewStyle().
			Foreground(borderColor)

	badgeStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Background(lipgloss.Color("#312E81")).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)
)
