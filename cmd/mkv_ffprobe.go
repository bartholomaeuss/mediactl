package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const ffprobeShowEntries = "stream=index,codec_name,codec_type,profile:stream_disposition=default:stream_tags=language,title"

// mkvFFProbeCmd creates JSON snapshots for every MKV in a directory using ffprobe.
var mkvFFProbeCmd = &cobra.Command{
	Use:   "ffprobe [directory]",
	Short: "Export ffprobe JSON for every MKV file in a folder",
	Long: `Runs ffprobe for each *.mkv file inside the specified directory (current directory by default)
and writes a .json file with the ffprobe output next to the MKV file.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) == 1 {
			dir = args[0]
		}

		return runFFProbeForDir(dir, cmd)
	},
}

var ffprobeRunner = executeFFProbe

func init() {
	mkvCmd.AddCommand(mkvFFProbeCmd)
}

func runFFProbeForDir(dir string, cmd *cobra.Command) error {
	processed, err := walkFFProbe(dir, ".mkv", cmd, ffprobeRunner)
	if err != nil {
		return err
	}

	if processed == 0 {
		cmd.Printf("No MKV files found in %s\n", dir)
	}

	return nil
}

func walkFFProbe(dir, extension string, cmd *cobra.Command, runner func(input, output string) error) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, fmt.Errorf("read directory %s: %w", dir, err)
	}

	var processed int
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			count, err := walkFFProbe(path, extension, cmd, runner)
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
		output := path + ".json"

		if err := runner(path, output); err != nil {
			return 0, err
		}

		cmd.Printf("wrote sidecar %s\n", output)
	}

	return processed, nil
}

func executeFFProbe(input, output string) error {
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
