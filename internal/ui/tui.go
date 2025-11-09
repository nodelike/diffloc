package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
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
	viewport    viewport.Model
	ready       bool
}

// NewModel creates a new TUI model
func NewModel(stats *model.Stats) Model {
	return Model{
		stats:       stats,
		sortMode:    model.SortByLines,
		sortReverse: false, // Ascending by default
		ready:       false,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := 3
		if !m.ready {
			// Initialize viewport
			m.viewport = viewport.New(msg.Width, msg.Height-headerHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.renderFullContent())
			m.viewport.GotoBottom() // Start at bottom
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - headerHeight
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		// Sorting controls
		case "n":
			if m.sortMode == model.SortByName {
				m.sortReverse = !m.sortReverse
			} else {
				m.sortMode = model.SortByName
				m.sortReverse = false
			}
			m.sortFiles()
			m.viewport.SetContent(m.renderFullContent())
			m.viewport.GotoBottom()
			return m, nil
		case "l":
			if m.sortMode == model.SortByLines {
				m.sortReverse = !m.sortReverse
			} else {
				m.sortMode = model.SortByLines
				m.sortReverse = true
			}
			m.sortFiles()
			m.viewport.SetContent(m.renderFullContent())
			m.viewport.GotoBottom()
			return m, nil
		case "a":
			if m.sortMode == model.SortByAdditions {
				m.sortReverse = !m.sortReverse
			} else {
				m.sortMode = model.SortByAdditions
				m.sortReverse = true
			}
			m.sortFiles()
			m.viewport.SetContent(m.renderFullContent())
			m.viewport.GotoBottom()
			return m, nil
		case "d":
			if m.sortMode == model.SortByDeletions {
				m.sortReverse = !m.sortReverse
			} else {
				m.sortMode = model.SortByDeletions
				m.sortReverse = true
			}
			m.sortFiles()
			m.viewport.SetContent(m.renderFullContent())
			m.viewport.GotoBottom()
			return m, nil
		}
	}

	// Pass through to viewport for scrolling
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the TUI
func (m Model) View() string {
	if m.err != nil {
		return "\n" + errorStyle.Render(fmt.Sprintf("⚠️  Error: %v", m.err)) + "\n"
	}

	if !m.ready {
		return "\nInitializing..."
	}

	// Render viewport and footer
	isGitRepo := m.stats.TotalAdditions > 0 || m.stats.TotalDeletions > 0 || m.stats.ChangedCount > 0

	return fmt.Sprintf("%s\n%s\n", m.viewport.View(), m.renderFooter(isGitRepo))
}

// renderFullContent renders the complete content without scrolling
func (m Model) renderFullContent() string {
	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(headerStyle.Render("✨ diffloc — Diff Line Counter"))
	b.WriteString("\n")

	// Check if this is a git repo (has any changes tracked)
	isGitRepo := m.stats.TotalAdditions > 0 || m.stats.TotalDeletions > 0 || m.stats.ChangedCount > 0

	if isGitRepo {
		// Git repo: Show changed and unchanged files separately
		if len(m.stats.ChangedFiles) > 0 {
			changedBadge := badgeStyle.Render(fmt.Sprintf("%d", len(m.stats.ChangedFiles)))
			b.WriteString(sectionHeaderStyle.Render(changedBadge + " Changed Files"))
			b.WriteString("\n")
			b.WriteString(m.renderFileTable(m.stats.ChangedFiles, true, true))
		}

		if len(m.stats.UnchangedFiles) > 0 {
			unchangedBadge := badgeStyle.Render(fmt.Sprintf("%d", len(m.stats.UnchangedFiles)))
			b.WriteString(sectionHeaderStyle.Render(unchangedBadge + " Unchanged Files"))
			b.WriteString("\n")
			b.WriteString(m.renderFileTable(m.stats.UnchangedFiles, false, true))
		}
	} else {
		// Non-git: Show all files in one section without git-specific columns
		allFiles := append(m.stats.ChangedFiles, m.stats.UnchangedFiles...)
		if len(allFiles) > 0 {
			filesBadge := badgeStyle.Render(fmt.Sprintf("%d", len(allFiles)))
			b.WriteString(sectionHeaderStyle.Render(filesBadge + " Files"))
			b.WriteString("\n")
			b.WriteString(m.renderFileTable(allFiles, false, false))
		}
	}

	// Summary
	b.WriteString(m.renderSummary(isGitRepo))

	return b.String()
}

// renderFileTable renders a table of files
func (m Model) renderFileTable(files []*model.FileInfo, isChanged bool, showGitColumns bool) string {
	if len(files) == 0 {
		return mutedNumberStyle.Render("    (none)") + "\n"
	}

	var b strings.Builder

	// Table header
	b.WriteString("    ")
	b.WriteString(tableHeaderStyle.Render(fmt.Sprintf("%-10s", "LINES")))
	b.WriteString("  ")

	if showGitColumns {
		b.WriteString(tableHeaderStyle.Render(fmt.Sprintf("%-10s", "ADDED")))
		b.WriteString("  ")
		b.WriteString(tableHeaderStyle.Render(fmt.Sprintf("%-10s", "REMOVED")))
		b.WriteString("  ")
	}

	b.WriteString(tableHeaderStyle.Render("FILE PATH"))
	b.WriteString("\n")

	// Separator line
	b.WriteString("    ")
	sepLength := 90
	if !showGitColumns {
		sepLength = 60
	}
	b.WriteString(separatorStyle.Render(strings.Repeat("─", sepLength)))
	b.WriteString("\n")

	// File rows
	for _, file := range files {
		b.WriteString("    ")

		// Lines count
		linesStr := fmt.Sprintf("%-10d", file.Lines)
		b.WriteString(summaryValueStyle.Render(linesStr))
		b.WriteString("  ")

		if showGitColumns {
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
		}

		// File path with visual indicator
		pathPrefix := ""
		if isChanged && showGitColumns {
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
func (m Model) renderSummary(isGitRepo bool) string {
	var content strings.Builder

	// Title
	content.WriteString(tableHeaderStyle.Render("📊 SUMMARY"))
	content.WriteString("\n")
	content.WriteString(separatorStyle.Render(strings.Repeat("─", 60)))
	content.WriteString("\n")

	if isGitRepo {
		// Git repo: Show net change
		var netChangeStr string
		var netChangeIcon string
		var netChangeStyle lipgloss.Style
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

		// File counts with changed/unchanged breakdown
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
	} else {
		// Non-git: Simple summary
		content.WriteString(summaryLabelStyle.Render("Total Files:"))
		content.WriteString(" ")
		content.WriteString(summaryValueStyle.Render(fmt.Sprintf("%d", m.stats.TotalFiles)))
		content.WriteString("\n")
		content.WriteString(summaryLabelStyle.Render("Total Lines:"))
		content.WriteString(" ")
		content.WriteString(summaryValueStyle.Render(fmt.Sprintf("%d", m.stats.TotalLines)))
	}

	return summaryBoxStyle.Render(content.String())
}

// renderFooter renders the footer with keybindings
func (m Model) renderFooter(isGitRepo bool) string {
	var footer strings.Builder

	// Keybindings
	keybindings := []string{
		keybindingKeyStyle.Render("↑↓/j/k") + " " + keybindingDescStyle.Render("scroll"),
		keybindingKeyStyle.Render("n") + " " + keybindingDescStyle.Render("name"),
		keybindingKeyStyle.Render("l") + " " + keybindingDescStyle.Render("lines"),
	}

	if isGitRepo {
		keybindings = append(keybindings,
			keybindingKeyStyle.Render("a")+" "+keybindingDescStyle.Render("additions"),
			keybindingKeyStyle.Render("d")+" "+keybindingDescStyle.Render("deletions"),
		)
	}

	keybindings = append(keybindings, keybindingKeyStyle.Render("q")+" "+keybindingDescStyle.Render("quit"))

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
	footer.WriteString("\n\n")
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

	// Use alternate screen for clean TUI experience
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	
	_, err := p.Run()
	return err
}

// PrintStatic prints the output once without interactivity
func PrintStatic(stats *model.Stats) error {
	m := NewModel(stats)
	m.sortFiles() // Initial sort

	// Print the view without footer (no interactivity needed)
	fmt.Print(m.renderStatic())
	return nil
}

// renderStatic renders the static view without footer
func (m Model) renderStatic() string {
	if m.err != nil {
		return "\n" + errorStyle.Render(fmt.Sprintf("⚠️  Error: %v", m.err)) + "\n"
	}

	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(headerStyle.Render("✨ diffloc — Diff Line Counter"))
	b.WriteString("\n")

	// Check if this is a git repo (has any changes tracked)
	isGitRepo := m.stats.TotalAdditions > 0 || m.stats.TotalDeletions > 0 || m.stats.ChangedCount > 0

	if isGitRepo {
		// Git repo: Show changed and unchanged files separately
		if len(m.stats.ChangedFiles) > 0 {
			changedBadge := badgeStyle.Render(fmt.Sprintf("%d", len(m.stats.ChangedFiles)))
			b.WriteString(sectionHeaderStyle.Render(changedBadge + " Changed Files"))
			b.WriteString("\n")
			b.WriteString(m.renderFileTable(m.stats.ChangedFiles, true, true))
		}

		if len(m.stats.UnchangedFiles) > 0 {
			unchangedBadge := badgeStyle.Render(fmt.Sprintf("%d", len(m.stats.UnchangedFiles)))
			b.WriteString(sectionHeaderStyle.Render(unchangedBadge + " Unchanged Files"))
			b.WriteString("\n")
			b.WriteString(m.renderFileTable(m.stats.UnchangedFiles, false, true))
		}
	} else {
		// Non-git: Show all files in one section without git-specific columns
		allFiles := append(m.stats.ChangedFiles, m.stats.UnchangedFiles...)
		if len(allFiles) > 0 {
			filesBadge := badgeStyle.Render(fmt.Sprintf("%d", len(allFiles)))
			b.WriteString(sectionHeaderStyle.Render(filesBadge + " Files"))
			b.WriteString("\n")
			b.WriteString(m.renderFileTable(allFiles, false, false))
		}
	}

	// Summary
	b.WriteString(m.renderSummary(isGitRepo))
	b.WriteString("\n")

	return b.String()
}
