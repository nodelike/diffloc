package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Color palette
	primaryColor   = lipgloss.Color("86")   // Cyan
	successColor   = lipgloss.Color("42")   // Green
	errorColor     = lipgloss.Color("196")  // Red
	warningColor   = lipgloss.Color("220")  // Yellow
	mutedColor     = lipgloss.Color("243")  // Gray
	highlightColor = lipgloss.Color("213")  // Pink

	// Title style
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	// Header styles
	headerStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1)

	sectionHeaderStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Underline(true).
				MarginTop(1).
				MarginBottom(0)

	// Table styles
	tableHeaderStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	// File info styles
	additionStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	deletionStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	mutedNumberStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	filePathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	unchangedFilePathStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	// Summary styles
	summaryBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			MarginTop(1)

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
				Foreground(mutedColor)

	summaryValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Bold(true)

	// Footer styles
	footerStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1).
			Italic(true)

	keybindingStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	separatorStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Error style
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Padding(1, 2)
)

