package analyzer

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"golang.org/x/sync/errgroup"
)

const readBufferSize = 64 * 1024

var bufPool = sync.Pool{New: func() any {
	b := make([]byte, readBufferSize)
	return &b
}}

// CountLines counts lines the way git does: a final line without a trailing
// newline still counts. Binary files (NUL byte in the first chunk) count as 0.
func CountLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	bp := bufPool.Get().(*[]byte)
	defer bufPool.Put(bp)
	buf := *bp

	count := 0
	first := true
	var last byte
	for {
		n, err := file.Read(buf)
		if n > 0 {
			chunk := buf[:n]
			if first {
				if bytes.IndexByte(chunk, 0) != -1 {
					return 0, nil
				}
				first = false
			}
			count += bytes.Count(chunk, []byte{'\n'})
			last = chunk[n-1]
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}

	if !first && last != '\n' {
		count++
	}
	return count, nil
}

// countResult is the line count of one file, or why it could not be read.
type countResult struct {
	lines   int
	missing bool
	err     error
}

// countAll counts lines of root-relative paths in parallel. Results line up with paths.
func countAll(ctx context.Context, root string, paths []string) ([]countResult, error) {
	results := make([]countResult, len(paths))

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(runtime.NumCPU() * 2)

	for i, p := range paths {
		eg.Go(func() error {
			if err := ctx.Err(); err != nil {
				return err
			}
			lines, err := CountLines(filepath.Join(root, filepath.FromSlash(p)))
			switch {
			case errors.Is(err, fs.ErrNotExist):
				results[i] = countResult{missing: true}
			case err != nil:
				results[i] = countResult{err: err}
			default:
				results[i] = countResult{lines: lines}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}
