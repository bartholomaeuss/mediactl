// Package exec wraps external ffmpeg/ffprobe execution for mediactl workflows.
package exec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// ffprobeShowEntries defines the stream metadata fields we persist to sidecars.
const ffprobeShowEntries = "stream=index,codec_name,codec_type,profile:stream_disposition=default:stream_tags=language,title,handler_name"

// ExecuteFFMpegRemux remuxes an input file into the output while applying handler_name tags.
func ExecuteFFMpegRemux(input, output string, handlerNames map[int]string) error {
	args := []string{
		"-v", "info",
		"-i", input,
		"-map", "0",
		"-map", "-0:d",
		"-c", "copy",
		"-map_metadata", "0",
	}

	var indices []int
	for index := range handlerNames {
		indices = append(indices, index)
	}
	sort.Ints(indices)
	for _, index := range indices {
		args = append(args, "-metadata:s:"+strconv.Itoa(index), "handler_name="+handlerNames[index])
	}

	args = append(args, output)

	fmt.Fprintf(os.Stdout, "ffmpeg %s\n", formatCommandArgs(args))

	command := exec.Command("ffmpeg", args...)
	var stderr bytes.Buffer
	command.Stdout = os.Stdout
	command.Stderr = io.MultiWriter(os.Stderr, &stderr)

	if err := command.Run(); err != nil {
		return fmt.Errorf("ffmpeg remux %s failed: %v: %s", input, err, strings.TrimSpace(stderr.String()))
	}

	return nil
}

// ExecuteFFProbe runs ffprobe and writes the normalized JSON output to the given file path.
func ExecuteFFProbe(input, output string) error {
	command := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", ffprobeShowEntries,
		"-print_format", "json",
		input,
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		return fmt.Errorf("ffprobe %s failed: %v: %s", input, err, strings.TrimSpace(stderr.String()))
	}

	normalized, err := normalizeHandlerNameKey(stdout.Bytes())
	if err != nil {
		return fmt.Errorf("normalize handler_name for %s: %w", input, err)
	}

	if err := os.WriteFile(output, normalized, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", output, err)
	}

	return nil
}

// normalizeHandlerNameKey ensures handler_name uses lowercase keys for consistency.
func normalizeHandlerNameKey(payload []byte) ([]byte, error) {
	var data map[string]any
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}

	streams, ok := data["streams"].([]any)
	if !ok {
		return payload, nil
	}

	for _, stream := range streams {
		streamMap, ok := stream.(map[string]any)
		if !ok {
			continue
		}
		tags, ok := streamMap["tags"].(map[string]any)
		if !ok {
			continue
		}

		var handlerValue any
		for key, value := range tags {
			if strings.EqualFold(key, "handler_name") {
				handlerValue = value
				if key != "handler_name" {
					delete(tags, key)
				}
			}
		}

		if handlerValue != nil {
			tags["handler_name"] = handlerValue
		}
	}

	return json.MarshalIndent(data, "", "    ")
}

// formatCommandArgs quotes command arguments for readable logging.
func formatCommandArgs(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = strconv.Quote(arg)
	}
	return strings.Join(quoted, " ")
}
