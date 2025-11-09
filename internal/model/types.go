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
	ChangedFiles   []*FileInfo
	UnchangedFiles []*FileInfo
	TotalFiles     int
	ChangedCount   int
	UnchangedCount int
	TotalLines     int
	TotalAdditions int
	TotalDeletions int
	NetChange      int
}

// SortMode defines how files should be sorted
type SortMode int

const (
	SortByName SortMode = iota
	SortByLines
	SortByAdditions
	SortByDeletions
)

func (s SortMode) String() string {
	switch s {
	case SortByName:
		return "name"
	case SortByLines:
		return "lines"
	case SortByAdditions:
		return "additions"
	case SortByDeletions:
		return "deletions"
	default:
		return "name"
	}
}

