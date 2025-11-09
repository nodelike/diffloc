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
	stats          *model.Stats
	sortMode       model.SortMode
	sortReverse    bool // Track if numeric sort is reversed
	err            error
	scrollOffset   int  // Current scroll position (line-based)
	viewportHeight int  // Available height for viewport
	contentHeight  int  // Total content height
}

// NewModel creates a new TUI model
func NewModel(stats *model.Stats) Model {
	return Model{
		stats:          stats,
		sortMode:       model.SortByLines,
		sortReverse:    false, // Ascending by default
		scrollOffset:   0,
		viewportHeight: 40, // Default, will be updated based on terminal size
		contentHeight:  0,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Calculate content height for scroll bounds
	// We need to render content to know its height
	content := m.renderFullContent()
	lines := strings.Split(content, "\n")
	m.contentHeight = len(lines)
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update viewport height - leave room for footer
		m.viewportHeight = msg.Height - 3
		if m.viewportHeight < 10 {
			m.viewportHeight = 10
		}
		
	case tea.KeyMsg:
		maxScroll := m.contentHeight - m.viewportHeight
		if maxScroll < 0 {
			maxScroll = 0
		}

		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
			
		// Scrolling controls
		case "up", "k":
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}
		case "down", "j":
			if m.scrollOffset < maxScroll {
				m.scrollOffset++
			}
		case "pgup":
			m.scrollOffset -= m.viewportHeight / 2
			if m.scrollOffset < 0 {
				m.scrollOffset = 0
			}
		case "pgdown":
			m.scrollOffset += m.viewportHeight / 2
			if m.scrollOffset > maxScroll {
				m.scrollOffset = maxScroll
			}
		case "home", "g":
			m.scrollOffset = 0
		case "end", "G":
			m.scrollOffset = maxScroll

		// Sorting controls
		case "n":
			if m.sortMode == model.SortByName {
				m.sortReverse = !m.sortReverse
			} else {
				m.sortMode = model.SortByName
				m.sortReverse = false
			}
			m.sortFiles()
			m.scrollOffset = 0 // Reset scroll on sort change
		case "l":
			if m.sortMode == model.SortByLines {
				m.sortReverse = !m.sortReverse
			} else {
				m.sortMode = model.SortByLines
				m.sortReverse = true // Default descending for numbers
			}
			m.sortFiles()
			m.scrollOffset = 0
		case "a":
			if m.sortMode == model.SortByAdditions {
				m.sortReverse = !m.sortReverse
			} else {
				m.sortMode = model.SortByAdditions
				m.sortReverse = true // Default descending for numbers
			}
			m.sortFiles()
			m.scrollOffset = 0
		case "d":
			if m.sortMode == model.SortByDeletions {
				m.sortReverse = !m.sortReverse
			} else {
				m.sortMode = model.SortByDeletions
				m.sortReverse = true // Default descending for numbers
			}
			m.sortFiles()
			m.scrollOffset = 0
		}
	}
	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	if m.err != nil {
		return "\n" + errorStyle.Render(fmt.Sprintf("⚠️  Error: %v", m.err)) + "\n"
	}

	// Render the full content first
	content := m.renderFullContent()
	
	// Split into lines
	lines := strings.Split(content, "\n")
	
	// Calculate visible range
	maxScroll := m.contentHeight - m.viewportHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	
	scrollOffset := m.scrollOffset
	if scrollOffset > maxScroll {
		scrollOffset = maxScroll
	}
	if scrollOffset < 0 {
		scrollOffset = 0
	}
	
	start := scrollOffset
	end := scrollOffset + m.viewportHeight
	if end > len(lines) {
		end = len(lines)
	}
	if start < 0 {
		start = 0
	}
	
	// Get visible lines
	visibleLines := lines[start:end]
	
	// Add footer with scroll indicators
	var result strings.Builder
	result.WriteString(strings.Join(visibleLines, "\n"))
	result.WriteString("\n")
	
	// Render footer
	isGitRepo := m.stats.TotalAdditions > 0 || m.stats.TotalDeletions > 0 || m.stats.ChangedCount > 0
	result.WriteString(m.renderFooter(isGitRepo))
	
	// Add scroll indicators
	if scrollOffset > 0 {
		result.WriteString(mutedNumberStyle.Render(fmt.Sprintf(" • Line %d/%d", scrollOffset+1, m.contentHeight)))
	}
	if scrollOffset < maxScroll {
		result.WriteString(mutedNumberStyle.Render(" • More below ↓"))
	}
	result.WriteString("\n")
	
	return result.String()
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
			keybindingKeyStyle.Render("a") + " " + keybindingDescStyle.Render("additions"),
			keybindingKeyStyle.Render("d") + " " + keybindingDescStyle.Render("deletions"),
		)
	}
	
	keybindings = append(keybindings, keybindingKeyStyle.Render("q") + " " + keybindingDescStyle.Render("quit"))

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
