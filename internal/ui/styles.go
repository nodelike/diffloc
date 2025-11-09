package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Modern color palette - using hex colors for better consistency
	primaryColor    = lipgloss.Color("#AF87FF") // Vibrant purple
	accentColor     = lipgloss.Color("#06B6D4") // Cyan blue
	successColor    = lipgloss.Color("#10B981") // Modern green
	errorColor      = lipgloss.Color("#EF4444") // Modern red
	warningColor    = lipgloss.Color("#F59E0B") // Amber
	mutedColor      = lipgloss.Color("#6B7280") // Neutral gray
	lightMutedColor = lipgloss.Color("#9CA3AF") // Light gray
	textColor       = lipgloss.Color("#F9FAFB") // Off-white
	highlightColor  = lipgloss.Color("#EC4899") // Pink
	borderColor     = lipgloss.Color("#4B5563") // Dark gray
	backgroundColor = lipgloss.Color("#1F2937") // Dark background

	// Title style - main app header
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 2)

	// Header styles
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

	// Table styles
	tableHeaderStyle = lipgloss.NewStyle().
				Foreground(lightMutedColor).
				Bold(true).
				Padding(0, 1)

	tableRowStyle = lipgloss.NewStyle().
			Padding(0, 1)

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
			Foreground(textColor).
			Bold(false)

	unchangedFilePathStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true)

	// Summary styles
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

	// Stat card style for individual stats
	statCardStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			MarginRight(2)

	// Footer styles
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

	// Error style
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Padding(2, 4)

	// Icon/badge styles
	badgeStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Background(lipgloss.Color("#312E81")).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)
)
