package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nodelike/diffloc/internal/model"
)

// PrintStatic renders the analysis results to stdout.
func PrintStatic(stats *model.Stats) {
	sortByLines(stats)
	fmt.Print(renderStatic(stats))
}

// sortByLines sorts changed and unchanged files by line count (ascending),
// keeping the analyzer's path order for ties.
func sortByLines(stats *model.Stats) {
	byLines := func(files []*model.FileInfo) {
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].Lines < files[j].Lines
		})
	}
	byLines(stats.ChangedFiles)
	byLines(stats.UnchangedFiles)
}

// renderStatic renders the complete report.
func renderStatic(stats *model.Stats) string {
	var b strings.Builder

	b.WriteString(headerStyle.Render("✨ diffloc — Diff Line Counter"))
	isGitRepo := stats.IsGitRepo

	if isGitRepo {
		if len(stats.UnchangedFiles) > 0 {
			unchangedBadge := badgeStyle.Render(fmt.Sprintf("%d", len(stats.UnchangedFiles)))
			b.WriteString(sectionHeaderStyle.Render(unchangedBadge + " Unchanged Files"))
			b.WriteString("\n")
			b.WriteString(renderFileTable(stats.UnchangedFiles, false, true))
		}

		if len(stats.ChangedFiles) > 0 {
			changedBadge := badgeStyle.Render(fmt.Sprintf("%d", len(stats.ChangedFiles)))
			b.WriteString(sectionHeaderStyle.Render(changedBadge + " Changed Files"))
			b.WriteString("\n")
			b.WriteString(renderFileTable(stats.ChangedFiles, true, true))
		}
	} else {
		allFiles := stats.UnchangedFiles
		if len(allFiles) > 0 {
			filesBadge := badgeStyle.Render(fmt.Sprintf("%d", len(allFiles)))
			b.WriteString(sectionHeaderStyle.Render(filesBadge + " Files"))
			b.WriteString("\n")
			b.WriteString(renderFileTable(allFiles, false, false))
		}
	}

	b.WriteString(renderSummary(stats, isGitRepo))
	b.WriteString("\n")

	return b.String()
}

// renderFileTable renders a table of files.
func renderFileTable(files []*model.FileInfo, isChanged, showGitColumns bool) string {
	if len(files) == 0 {
		return mutedNumberStyle.Render("    (none)") + "\n"
	}

	var b strings.Builder

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

	b.WriteString("    ")
	sepLength := 90
	if !showGitColumns {
		sepLength = 60
	}
	b.WriteString(separatorStyle.Render(strings.Repeat("─", sepLength)))
	b.WriteString("\n")

	for _, file := range files {
		b.WriteString("    ")

		linesStr := fmt.Sprintf("%-10d", file.Lines)
		b.WriteString(summaryValueStyle.Render(linesStr))
		b.WriteString("  ")

		if showGitColumns {
			if file.Additions > 0 {
				addStr := fmt.Sprintf("+%-9d", file.Additions)
				b.WriteString(additionStyle.Render(addStr))
			} else {
				b.WriteString(mutedNumberStyle.Render(fmt.Sprintf("%-10s", "—")))
			}
			b.WriteString("  ")

			if file.Deletions > 0 {
				delStr := fmt.Sprintf("-%-9d", file.Deletions)
				b.WriteString(deletionStyle.Render(delStr))
			} else {
				b.WriteString(mutedNumberStyle.Render(fmt.Sprintf("%-10s", "—")))
			}
			b.WriteString("  ")
		}

		pathPrefix := ""
		if isChanged && showGitColumns {
			if file.Additions > 0 && file.Deletions > 0 {
				pathPrefix = "◆ "
			} else if file.Additions > 0 {
				pathPrefix = "+ "
			} else if file.Deletions > 0 {
				pathPrefix = "- "
			}
		}
		b.WriteString(filePathStyle.Render(pathPrefix + file.Path))
		b.WriteString("\n")
	}

	return b.String()
}

// renderSummary renders the summary box.
func renderSummary(stats *model.Stats, isGitRepo bool) string {
	var content strings.Builder

	content.WriteString(tableHeaderStyle.Render("📊 SUMMARY"))
	content.WriteString("\n")
	content.WriteString(separatorStyle.Render(strings.Repeat("─", 60)))
	content.WriteString("\n")

	if isGitRepo {
		var netChangeStr string
		var netChangeIcon string
		var netChangeStyle lipgloss.Style
		if stats.NetChange > 0 {
			netChangeIcon = "▲"
			netChangeStr = fmt.Sprintf("+%d lines", stats.NetChange)
			netChangeStyle = summaryPositiveStyle
		} else if stats.NetChange < 0 {
			netChangeIcon = "▼"
			netChangeStr = fmt.Sprintf("%d lines", stats.NetChange)
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

		accentStyle := lipgloss.NewStyle().Foreground(accentColor).Bold(true)
		content.WriteString(summaryLabelStyle.Render("Files:"))
		content.WriteString("       ")
		content.WriteString(summaryValueStyle.Render(fmt.Sprintf("%d", stats.TotalFiles)))
		content.WriteString(summaryLabelStyle.Render(" total  •  "))
		content.WriteString(accentStyle.Render(fmt.Sprintf("%d", stats.ChangedCount)))
		content.WriteString(summaryLabelStyle.Render(" changed  •  "))
		content.WriteString(mutedNumberStyle.Render(fmt.Sprintf("%d", stats.UnchangedCount)))
		content.WriteString(summaryLabelStyle.Render(" unchanged"))
		content.WriteString("\n")

		content.WriteString(summaryLabelStyle.Render("Total Lines:"))
		content.WriteString(" ")
		content.WriteString(summaryValueStyle.Render(fmt.Sprintf("%d", stats.TotalLines)))
		content.WriteString("\n")

		content.WriteString(summaryLabelStyle.Render("Changes:"))
		content.WriteString("    ")
		content.WriteString(additionStyle.Render(fmt.Sprintf("+%d", stats.TotalAdditions)))
		content.WriteString(summaryLabelStyle.Render(" added  •  "))
		content.WriteString(deletionStyle.Render(fmt.Sprintf("-%d", stats.TotalDeletions)))
		content.WriteString(summaryLabelStyle.Render(" removed"))
	} else {
		content.WriteString(summaryLabelStyle.Render("Total Files:"))
		content.WriteString(" ")
		content.WriteString(summaryValueStyle.Render(fmt.Sprintf("%d", stats.TotalFiles)))
		content.WriteString("\n")
		content.WriteString(summaryLabelStyle.Render("Total Lines:"))
		content.WriteString(" ")
		content.WriteString(summaryValueStyle.Render(fmt.Sprintf("%d", stats.TotalLines)))
	}

	return summaryBoxStyle.Render(content.String())
}
