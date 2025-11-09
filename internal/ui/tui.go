package ui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nodelike/diffloc/internal/model"
)

// Model represents the TUI state
type Model struct {
	stats    *model.Stats
	sortMode model.SortMode
	err      error
}

// NewModel creates a new TUI model
func NewModel(stats *model.Stats) Model {
	return Model{
		stats:    stats,
		sortMode: model.SortByName,
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
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n":
			m.sortMode = model.SortByName
			m.sortFiles()
		case "l":
			m.sortMode = model.SortByLines
			m.sortFiles()
		case "a":
			m.sortMode = model.SortByAdditions
			m.sortFiles()
		case "d":
			m.sortMode = model.SortByDeletions
			m.sortFiles()
		}
	}
	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	var b strings.Builder

	// Header
	b.WriteString(headerStyle.Render("✨ diffloc - Diff Line Counter"))
	b.WriteString("\n\n")

	// Changed files section
	if len(m.stats.ChangedFiles) > 0 {
		b.WriteString(sectionHeaderStyle.Render("Changed Files:"))
		b.WriteString("\n")
		b.WriteString(m.renderFileTable(m.stats.ChangedFiles, true))
		b.WriteString("\n")
	}

	// Unchanged files section
	if len(m.stats.UnchangedFiles) > 0 {
		b.WriteString(sectionHeaderStyle.Render("Unchanged Files:"))
		b.WriteString("\n")
		b.WriteString(m.renderFileTable(m.stats.UnchangedFiles, false))
		b.WriteString("\n")
	}

	// Summary
	b.WriteString(m.renderSummary())
	b.WriteString("\n")

	// Footer with keybindings
	b.WriteString(m.renderFooter())

	return b.String()
}

// renderFileTable renders a table of files
func (m Model) renderFileTable(files []*model.FileInfo, isChanged bool) string {
	if len(files) == 0 {
		return mutedNumberStyle.Render("  (none)")
	}

	var b strings.Builder

	// Table header
	b.WriteString("  ")
	b.WriteString(tableHeaderStyle.Render(fmt.Sprintf("%-8s", "Lines")))
	b.WriteString("  ")
	b.WriteString(tableHeaderStyle.Render(fmt.Sprintf("%-8s", "+Add")))
	b.WriteString("  ")
	b.WriteString(tableHeaderStyle.Render(fmt.Sprintf("%-8s", "-Del")))
	b.WriteString("  ")
	b.WriteString(tableHeaderStyle.Render("File"))
	b.WriteString("\n")

	// Separator
	b.WriteString(mutedNumberStyle.Render("  " + strings.Repeat("─", 80)))
	b.WriteString("\n")

	// File rows
	for _, file := range files {
		b.WriteString("  ")

		// Lines count
		linesStr := fmt.Sprintf("%-8d", file.Lines)
		if isChanged {
			b.WriteString(summaryValueStyle.Render(linesStr))
		} else {
			b.WriteString(mutedNumberStyle.Render(linesStr))
		}
		b.WriteString("  ")

		// Additions
		addStr := fmt.Sprintf("+%-7d", file.Additions)
		if file.Additions > 0 {
			b.WriteString(additionStyle.Render(addStr))
		} else {
			b.WriteString(mutedNumberStyle.Render(addStr))
		}
		b.WriteString("  ")

		// Deletions
		delStr := fmt.Sprintf("-%-7d", file.Deletions)
		if file.Deletions > 0 {
			b.WriteString(deletionStyle.Render(delStr))
		} else {
			b.WriteString(mutedNumberStyle.Render(delStr))
		}
		b.WriteString("  ")

		// File path
		if isChanged {
			b.WriteString(filePathStyle.Render(file.Path))
		} else {
			b.WriteString(unchangedFilePathStyle.Render(file.Path))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// renderSummary renders the summary box
func (m Model) renderSummary() string {
	var content strings.Builder

	// Net change
	netChangeStr := ""
	netChangeStyle := summaryNeutralStyle
	if m.stats.NetChange > 0 {
		netChangeStr = fmt.Sprintf("+%d (increased)", m.stats.NetChange)
		netChangeStyle = summaryPositiveStyle
	} else if m.stats.NetChange < 0 {
		netChangeStr = fmt.Sprintf("%d (decreased)", m.stats.NetChange)
		netChangeStyle = summaryNegativeStyle
	} else {
		netChangeStr = "0 (no change)"
	}

	content.WriteString(summaryLabelStyle.Render("Net Change: "))
	content.WriteString(netChangeStyle.Render(netChangeStr))
	content.WriteString("\n\n")

	// File counts
	content.WriteString(summaryLabelStyle.Render("Total Files: "))
	content.WriteString(summaryValueStyle.Render(fmt.Sprintf("%d", m.stats.TotalFiles)))
	content.WriteString(summaryLabelStyle.Render(fmt.Sprintf(" (%d changed, %d unchanged)", 
		m.stats.ChangedCount, m.stats.UnchangedCount)))
	content.WriteString("\n")

	// Line counts
	content.WriteString(summaryLabelStyle.Render("Total Lines: "))
	content.WriteString(summaryValueStyle.Render(fmt.Sprintf("%d", m.stats.TotalLines)))
	content.WriteString("\n")

	content.WriteString(summaryLabelStyle.Render("Insertions:  "))
	content.WriteString(additionStyle.Render(fmt.Sprintf("+%d", m.stats.TotalAdditions)))
	content.WriteString("\n")

	content.WriteString(summaryLabelStyle.Render("Deletions:   "))
	content.WriteString(deletionStyle.Render(fmt.Sprintf("-%d", m.stats.TotalDeletions)))

	return summaryBoxStyle.Render(content.String())
}

// renderFooter renders the footer with keybindings
func (m Model) renderFooter() string {
	currentSort := fmt.Sprintf("(sorted by: %s)", m.sortMode.String())
	
	keybindings := []string{
		keybindingStyle.Render("n") + "=name",
		keybindingStyle.Render("l") + "=lines",
		keybindingStyle.Render("a") + "=additions",
		keybindingStyle.Render("d") + "=deletions",
		keybindingStyle.Render("q") + "=quit",
	}

	footer := strings.Join(keybindings, separatorStyle.Render(" | "))
	footer += " " + mutedNumberStyle.Render(currentSort)

	return footerStyle.Render(footer)
}

// sortFiles sorts the files based on the current sort mode
func (m *Model) sortFiles() {
	sortFunc := func(files []*model.FileInfo) {
		sort.Slice(files, func(i, j int) bool {
			switch m.sortMode {
			case model.SortByName:
				return files[i].Path < files[j].Path
			case model.SortByLines:
				return files[i].Lines > files[j].Lines
			case model.SortByAdditions:
				return files[i].Additions > files[j].Additions
			case model.SortByDeletions:
				return files[i].Deletions > files[j].Deletions
			default:
				return files[i].Path < files[j].Path
			}
		})
	}

	sortFunc(m.stats.ChangedFiles)
	sortFunc(m.stats.UnchangedFiles)
}

// Run starts the TUI application
func Run(stats *model.Stats) error {
	m := NewModel(stats)
	m.sortFiles() // Initial sort

	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

