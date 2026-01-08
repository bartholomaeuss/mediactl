package infra

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const ffprobeShowEntries = "stream=index,codec_name,codec_type,profile:stream_disposition=default:stream_tags=language,title,handler_name"

// ExecuteFFProbe runs ffprobe and writes the JSON output to the given file path.
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

	if err := os.WriteFile(output, stdout.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", output, err)
	}

	return nil
}
