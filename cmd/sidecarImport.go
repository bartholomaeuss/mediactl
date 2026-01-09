/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"mediactl/infra/exec"

	"github.com/spf13/cobra"
)

// sidecarData mirrors the minimal fields from a mediactl sidecar JSON file.
type sidecarData struct {
	Streams []sidecarStream `json:"streams"`
}

// sidecarStream captures the stream metadata needed to apply handler_name updates.
type sidecarStream struct {
	Index     int               `json:"index"`
	CodecType string            `json:"codec_type"`
	Tags      map[string]string `json:"tags"`
}

// sidecarImportCmd remuxes a media file into MKV while applying handler_name updates.
var sidecarImportCmd = &cobra.Command{
	Use:   "import [sidecar.json]",
	Short: "Remux a media file into MKV using handler_name values from a sidecar",
	Long: `Reads a mediactl-generated sidecar JSON file, derives the input media file
path, and remuxes it into an MKV container while applying handler_name updates
for any audio streams that differ.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sidecarPath := args[0]

		mediaPath, err := deriveMediaPath(sidecarPath)
		if err != nil {
			return err
		}

		output := remuxOutputPath(mediaPath)
		if _, err := os.Stat(output); err == nil {
			return fmt.Errorf("output already exists: %s", output)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("stat %s: %w", output, err)
		}

		desired, err := readSidecar(sidecarPath)
		if err != nil {
			return err
		}

		current, err := probeSidecar(mediaPath)
		if err != nil {
			return err
		}

		if hasDataStreams(current) {
			cmd.Printf("Note: data streams will be dropped for MKV remux: %s\n", mediaPath)
		}

		updates, indices, err := handlerUpdates(desired, current)
		if err != nil {
			return err
		}

		if len(updates) == 0 {
			cmd.Printf("No handler_name changes detected for %s\n", mediaPath)
			return nil
		}

		if err := exec.ExecuteFFMpegRemux(mediaPath, output, updates); err != nil {
			return err
		}

		cmd.Printf("wrote remux %s (updated streams: %s)\n", output, indices)
		return nil
	},
}

// init registers the import subcommand under the sidecar namespace.
func init() {
	sidecarCmd.AddCommand(sidecarImportCmd)
}

// deriveMediaPath resolves the media file path implied by a sidecar JSON filename.
func deriveMediaPath(sidecarPath string) (string, error) {
	if !strings.HasSuffix(strings.ToLower(sidecarPath), ".json") {
		return "", fmt.Errorf("sidecar path must end with .json: %s", sidecarPath)
	}

	base := strings.TrimSuffix(sidecarPath, ".json")
	if _, err := os.Stat(base); err == nil {
		return base, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat %s: %w", base, err)
	}

	re := regexp.MustCompile(`^(.*)\.(\d+)$`)
	if matches := re.FindStringSubmatch(base); len(matches) == 3 {
		candidate := matches[1]
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("stat %s: %w", candidate, err)
		}
	}

	return "", fmt.Errorf("could not derive media file from sidecar %s", sidecarPath)
}

// remuxOutputPath builds the output MKV path for a given media file path.
func remuxOutputPath(mediaPath string) string {
	ext := filepath.Ext(mediaPath)
	base := strings.TrimSuffix(mediaPath, ext)
	return base + ".remux.mkv"
}

// readSidecar loads a sidecar JSON file and decodes the stream entries.
func readSidecar(path string) (*sidecarData, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read sidecar %s: %w", path, err)
	}

	var data sidecarData
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("parse sidecar %s: %w", path, err)
	}

	return &data, nil
}

// probeSidecar runs ffprobe against a media file and returns the parsed sidecar data.
func probeSidecar(mediaPath string) (*sidecarData, error) {
	tmp, err := os.CreateTemp("", "mediactl-ffprobe-*.json")
	if err != nil {
		return nil, fmt.Errorf("create temp ffprobe file: %w", err)
	}
	tmpPath := tmp.Name()
	if err := tmp.Close(); err != nil {
		return nil, fmt.Errorf("close temp ffprobe file: %w", err)
	}
	defer os.Remove(tmpPath)

	if err := exec.ExecuteFFProbe(mediaPath, tmpPath); err != nil {
		return nil, err
	}

	return readSidecar(tmpPath)
}

// handlerUpdates computes handler_name changes by comparing desired vs current streams.
func handlerUpdates(desired, current *sidecarData) (map[int]string, string, error) {
	currentHandlers := make(map[int]string, len(current.Streams))
	for _, stream := range current.Streams {
		if stream.Tags == nil {
			continue
		}
		if handlerName, ok := findHandlerName(stream.Tags); ok {
			currentHandlers[stream.Index] = handlerName
		}
	}

	updates := make(map[int]string)
	var updatedIndices []int

	for _, stream := range desired.Streams {
		if stream.Tags == nil {
			continue
		}
		handlerName, ok := findHandlerName(stream.Tags)
		if !ok || handlerName == "" {
			continue
		}

		currentHandler := currentHandlers[stream.Index]
		if handlerName == currentHandler {
			continue
		}

		if stream.CodecType != "" && stream.CodecType != "audio" {
			return nil, "", fmt.Errorf("handler_name update requested for non-audio stream index %d", stream.Index)
		}

		updates[stream.Index] = handlerName
		updatedIndices = append(updatedIndices, stream.Index)
	}

	sort.Ints(updatedIndices)
	indices := make([]string, len(updatedIndices))
	for i, index := range updatedIndices {
		indices[i] = fmt.Sprintf("%d", index)
	}

	return updates, strings.Join(indices, ", "), nil
}

// hasDataStreams reports whether any stream is marked with codec_type "data".
func hasDataStreams(data *sidecarData) bool {
	for _, stream := range data.Streams {
		if stream.CodecType == "data" {
			return true
		}
	}
	return false
}

// findHandlerName returns the handler_name tag value, matching keys case-insensitively.
func findHandlerName(tags map[string]string) (string, bool) {
	for key, value := range tags {
		if strings.EqualFold(key, "handler_name") {
			return value, true
		}
	}
	return "", false
}
