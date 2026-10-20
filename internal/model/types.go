package model

// FileInfo represents information about a single file
type FileInfo struct {
	Path      string
	Lines     int
	Additions int
	Deletions int
	IsChanged bool
}

// Stats represents aggregated statistics
type Stats struct {
	IsGitRepo      bool
	ChangedFiles   []*FileInfo
	UnchangedFiles []*FileInfo
	TotalFiles     int
	ChangedCount   int
	UnchangedCount int
	TotalLines     int
	TotalAdditions int
	TotalDeletions int
	NetChange      int
	// Skipped lists files that matched the filter but could not be read.
	Skipped []string `json:",omitempty"`
}
