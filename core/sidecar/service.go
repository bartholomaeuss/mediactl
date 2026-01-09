// Package sidecar provides helpers for generating and naming ffprobe sidecars.
package sidecar

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Runner executes ffprobe for a media file and writes output to the given path.
type Runner func(input, output string) error

// WalkFFProbe walks a directory tree, runs ffprobe for matching media files,
// and writes sidecar JSON next to each input file.
func WalkFFProbe(dir, extension string, out io.Writer, runner Runner) (int, error) {
	if out == nil {
		out = io.Discard
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, fmt.Errorf("read directory %s: %w", dir, err)
	}

	var processed int
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			count, err := WalkFFProbe(path, extension, out, runner)
			if err != nil {
				return 0, err
			}
			processed += count
			continue
		}

		if !strings.EqualFold(filepath.Ext(entry.Name()), extension) {
			continue
		}

		processed++
		output, err := NextAvailableSidecarPath(path)
		if err != nil {
			return 0, err
		}

		if err := runner(path, output); err != nil {
			return 0, err
		}

		fmt.Fprintf(out, "wrote sidecar %s\n", output)
	}

	return processed, nil
}

// NextAvailableSidecarPath picks the next available JSON sidecar filename,
// adding numeric suffixes when a base sidecar already exists.
func NextAvailableSidecarPath(mediaPath string) (string, error) {
	base := mediaPath + ".json"
	if _, err := os.Stat(base); err != nil {
		if os.IsNotExist(err) {
			return base, nil
		}
		return "", fmt.Errorf("stat %s: %w", base, err)
	}

	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s.%d.json", mediaPath, i)
		if _, err := os.Stat(candidate); err != nil {
			if os.IsNotExist(err) {
				return candidate, nil
			}
			return "", fmt.Errorf("stat %s: %w", candidate, err)
		}
	}
}
