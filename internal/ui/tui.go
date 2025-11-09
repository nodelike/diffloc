package ui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nodelike/diffloc/internal/model"
)

// Model represents the TUI state
type Model struct {
	stats       *model.Stats
	sortMode    model.SortMode
	sortReverse bool // Track if numeric sort is reversed
	err         error
}

// NewModel creates a new TUI model
func NewModel(stats *model.Stats) Model {
	return Model{
		stats:       stats,
		sortMode:    model.SortByLines,
		sortReverse: false, // Ascending by default
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "n":
			if m.sortMode == model.SortByName {
				m.sortReverse = !m.sortReverse
			} else {
				m.sortMode = model.SortByName
				m.sortReverse = false
			}
			m.sortFiles()
		case "l":
			if m.sortMode == model.SortByLines {
				m.sortReverse = !m.sortReverse
			} else {
				m.sortMode = model.SortByLines
				m.sortReverse = true // Default descending for numbers
			}
			m.sortFiles()
		case "a":
			if m.sortMode == model.SortByAdditions {
				m.sortReverse = !m.sortReverse
			} else {
				m.sortMode = model.SortByAdditions
				m.sortReverse = true // Default descending for numbers
			}
			m.sortFiles()
		case "d":
			if m.sortMode == model.SortByDeletions {
				m.sortReverse = !m.sortReverse
			} else {
				m.sortMode = model.SortByDeletions
				m.sortReverse = true // Default descending for numbers
			}
			m.sortFiles()
		}
	}
	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	if m.err != nil {
		return "\n" + errorStyle.Render(fmt.Sprintf("⚠️  Error: %v", m.err)) + "\n"
	}

	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(headerStyle.Render("✨ diffloc — Diff Line Counter"))
	b.WriteString("\n")

	// Changed files section
	if len(m.stats.ChangedFiles) > 0 {
		changedBadge := badgeStyle.Render(fmt.Sprintf("%d", len(m.stats.ChangedFiles)))
		b.WriteString(sectionHeaderStyle.Render(changedBadge + " Changed Files"))
		b.WriteString("\n")
		b.WriteString(m.renderFileTable(m.stats.ChangedFiles, true))
	}

	// Unchanged files section
	if len(m.stats.UnchangedFiles) > 0 {
		unchangedBadge := badgeStyle.Render(fmt.Sprintf("%d", len(m.stats.UnchangedFiles)))
		b.WriteString(sectionHeaderStyle.Render(unchangedBadge + " Unchanged Files"))
		b.WriteString("\n")
		b.WriteString(m.renderFileTable(m.stats.UnchangedFiles, false))
	}

	// Summary
	b.WriteString(m.renderSummary())

	// Footer with keybindings
	b.WriteString(m.renderFooter())
	b.WriteString("\n")

	return b.String()
}

// renderFileTable renders a table of files
func (m Model) renderFileTable(files []*model.FileInfo, isChanged bool) string {
	if len(files) == 0 {
		return mutedNumberStyle.Render("    (none)") + "\n"
	}

	var b strings.Builder

	// Table header with better spacing
	b.WriteString("    ")
	b.WriteString(tableHeaderStyle.Render(fmt.Sprintf("%-10s", "LINES")))
	b.WriteString("  ")
	b.WriteString(tableHeaderStyle.Render(fmt.Sprintf("%-10s", "ADDED")))
	b.WriteString("  ")
	b.WriteString(tableHeaderStyle.Render(fmt.Sprintf("%-10s", "REMOVED")))
	b.WriteString("  ")
	b.WriteString(tableHeaderStyle.Render("FILE PATH"))
	b.WriteString("\n")

	// Separator line
	b.WriteString("    ")
	b.WriteString(separatorStyle.Render(strings.Repeat("─", 90)))
	b.WriteString("\n")

	// File rows
	for _, file := range files {
		b.WriteString("    ")

		// Lines count - same style for all files
		linesStr := fmt.Sprintf("%-10d", file.Lines)
		b.WriteString(summaryValueStyle.Render(linesStr))
		b.WriteString("  ")

		// Additions with visual indicator
		if file.Additions > 0 {
			addStr := fmt.Sprintf("+%-9d", file.Additions)
			b.WriteString(additionStyle.Render(addStr))
		} else {
			b.WriteString(mutedNumberStyle.Render(fmt.Sprintf("%-10s", "—")))
		}
		b.WriteString("  ")

		// Deletions with visual indicator
		if file.Deletions > 0 {
			delStr := fmt.Sprintf("-%-9d", file.Deletions)
			b.WriteString(deletionStyle.Render(delStr))
		} else {
			b.WriteString(mutedNumberStyle.Render(fmt.Sprintf("%-10s", "—")))
		}
		b.WriteString("  ")

		// File path with visual indicator
		pathPrefix := ""
		if isChanged {
			if file.Additions > 0 && file.Deletions > 0 {
				pathPrefix = "◆ " // Modified
			} else if file.Additions > 0 {
				pathPrefix = "+ " // Added
			} else if file.Deletions > 0 {
				pathPrefix = "- " // Deleted
			}
		}
		b.WriteString(filePathStyle.Render(pathPrefix + file.Path))
		b.WriteString("\n")
	}

	return b.String()
}

// renderSummary renders the summary box
func (m Model) renderSummary() string {
	var content strings.Builder

	// Title
	content.WriteString(tableHeaderStyle.Render("📊 SUMMARY"))
	content.WriteString("\n")
	content.WriteString(separatorStyle.Render(strings.Repeat("─", 60)))
	content.WriteString("\n")

	// Net change with visual indicator
	netChangeStr := ""
	netChangeIcon := ""
	netChangeStyle := summaryNeutralStyle
	if m.stats.NetChange > 0 {
		netChangeIcon = "▲"
		netChangeStr = fmt.Sprintf("+%d lines", m.stats.NetChange)
		netChangeStyle = summaryPositiveStyle
	} else if m.stats.NetChange < 0 {
		netChangeIcon = "▼"
		netChangeStr = fmt.Sprintf("%d lines", m.stats.NetChange)
		netChangeStyle = summaryNegativeStyle
	} else {
		netChangeIcon = "●"
		netChangeStr = "no change"
		netChangeStyle = summaryNeutralStyle
	}

	content.WriteString(summaryLabelStyle.Render("Net Change:"))
	content.WriteString("  ")
	content.WriteString(netChangeStyle.Render(netChangeIcon + " " + netChangeStr))
	content.WriteString("\n")

	// File counts
	accentStyle := lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	content.WriteString(summaryLabelStyle.Render("Files:"))
	content.WriteString("       ")
	content.WriteString(summaryValueStyle.Render(fmt.Sprintf("%d", m.stats.TotalFiles)))
	content.WriteString(summaryLabelStyle.Render(" total  •  "))
	content.WriteString(accentStyle.Render(fmt.Sprintf("%d", m.stats.ChangedCount)))
	content.WriteString(summaryLabelStyle.Render(" changed  •  "))
	content.WriteString(mutedNumberStyle.Render(fmt.Sprintf("%d", m.stats.UnchangedCount)))
	content.WriteString(summaryLabelStyle.Render(" unchanged"))
	content.WriteString("\n")

	// Line counts
	content.WriteString(summaryLabelStyle.Render("Total Lines:"))
	content.WriteString(" ")
	content.WriteString(summaryValueStyle.Render(fmt.Sprintf("%d", m.stats.TotalLines)))
	content.WriteString("\n")

	// Changes breakdown
	content.WriteString(summaryLabelStyle.Render("Changes:"))
	content.WriteString("    ")
	content.WriteString(additionStyle.Render(fmt.Sprintf("+%d", m.stats.TotalAdditions)))
	content.WriteString(summaryLabelStyle.Render(" added  •  "))
	content.WriteString(deletionStyle.Render(fmt.Sprintf("-%d", m.stats.TotalDeletions)))
	content.WriteString(summaryLabelStyle.Render(" removed"))

	return summaryBoxStyle.Render(content.String())
}

// renderFooter renders the footer with keybindings
func (m Model) renderFooter() string {
	var footer strings.Builder

	// Keybindings
	keybindings := []string{
		keybindingKeyStyle.Render("n") + " " + keybindingDescStyle.Render("name"),
		keybindingKeyStyle.Render("l") + " " + keybindingDescStyle.Render("lines"),
		keybindingKeyStyle.Render("a") + " " + keybindingDescStyle.Render("additions"),
		keybindingKeyStyle.Render("d") + " " + keybindingDescStyle.Render("deletions"),
		keybindingKeyStyle.Render("q") + " " + keybindingDescStyle.Render("quit"),
	}

	footer.WriteString(strings.Join(keybindings, separatorStyle.Render("  •  ")))

	// Sort indicator
	sortDir := ""
	sortIcon := ""
	if m.sortMode != model.SortByName {
		if m.sortReverse {
			sortDir = "desc"
			sortIcon = "↓"
		} else {
			sortDir = "asc"
			sortIcon = "↑"
		}
	} else {
		sortDir = "A→Z"
		sortIcon = "⇅"
	}

	accentStyle := lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	footer.WriteString("\n")
	footer.WriteString(mutedNumberStyle.Render("Sort: "))
	footer.WriteString(accentStyle.Render(m.sortMode.String()))
	footer.WriteString(mutedNumberStyle.Render(" " + sortIcon + " " + sortDir))

	return footerStyle.Render(footer.String())
}

// sortFiles sorts the files based on the current sort mode and direction
func (m *Model) sortFiles() {
	sortFunc := func(files []*model.FileInfo) {
		sort.Slice(files, func(i, j int) bool {
			var less bool
			switch m.sortMode {
			case model.SortByName:
				less = files[i].Path < files[j].Path
			case model.SortByLines:
				less = files[i].Lines < files[j].Lines
			case model.SortByAdditions:
				less = files[i].Additions < files[j].Additions
			case model.SortByDeletions:
				less = files[i].Deletions < files[j].Deletions
			default:
				less = files[i].Path < files[j].Path
			}

			// Reverse if needed (for numeric sorts, default is descending)
			if m.sortReverse {
				return !less
			}
			return less
		})
	}

	sortFunc(m.stats.ChangedFiles)
	sortFunc(m.stats.UnchangedFiles)
}

// Run starts the TUI application
func Run(stats *model.Stats) error {
	m := NewModel(stats)
	m.sortFiles() // Initial sort

	// Use alternate screen for proper interactive TUI
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
