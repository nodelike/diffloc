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
	highlightColor  = lipgloss.Color("#EC4899")
	borderColor     = lipgloss.Color("#4B5563")
	backgroundColor = lipgloss.Color("#1F2937")

	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 2)

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

	tableRowStyle = lipgloss.NewStyle().
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

	unchangedFilePathStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true)

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

	statCardStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			MarginRight(2)

	footerStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1).
			Padding(0, 2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(borderColor)

	keybindingKeyStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Background(lipgloss.Color("#374151")).
				Padding(0, 1)

	keybindingDescStyle = lipgloss.NewStyle().
				Foreground(lightMutedColor)

	separatorStyle = lipgloss.NewStyle().
			Foreground(borderColor)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Padding(2, 4)

	badgeStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Background(lipgloss.Color("#312E81")).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)
)
