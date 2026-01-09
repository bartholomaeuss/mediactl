/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import (
	"fmt"
	"strings"

	"mediactl/core/sidecar"
	"mediactl/infra/exec"

	"github.com/spf13/cobra"
)

var (
	// sidecarExportFormat captures the --format flag for determining media extensions.
	sidecarExportFormat string
)

// sidecarExportCmd exports ffprobe JSON sidecars for a selected media format.
var sidecarExportCmd = &cobra.Command{
	Use:   "export [directory]",
	Short: "Generate ffprobe JSON sidecars for media files",
	Long: `Runs ffprobe for each matching media file in the specified directory
(current directory by default) and writes a JSON sidecar next to the media file.

Existing sidecars are preserved by writing to the next available numbered file
name, such as *.1.json or *.2.json.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		dir := "."
		if args[0] != "" {
			dir = args[0]
		}

		normalizedSidecarExportFormat := strings.ToLower(strings.TrimSpace(sidecarExportFormat))
		extension, err := sidecarFormatToExtension(normalizedSidecarExportFormat)
		if err != nil {
			return err
		}

		processed, err := sidecar.WalkFFProbe(dir, extension, cmd.OutOrStdout(), exec.ExecuteFFProbe)
		if err != nil {
			return err
		}

		if processed == 0 {
			cmd.Printf("No %s files found in %s\n", strings.ToUpper(normalizedSidecarExportFormat), dir)
		}

		return nil
	},
}

// init registers the export subcommand and its flags.
func init() {
	sidecarCmd.AddCommand(sidecarExportCmd)
	sidecarExportCmd.Flags().StringVar(&sidecarExportFormat, "format", "", "Media format to export (mkv or mp4)")
	_ = sidecarExportCmd.MarkFlagRequired("format")
}

// sidecarFormatToExtension maps a supported format name to its file extension.
func sidecarFormatToExtension(format string) (string, error) {
	switch format {
	case "mkv":
		return ".mkv", nil
	case "mp4":
		return ".mp4", nil
	case "":
		return "", fmt.Errorf("format is required")
	default:
		return "", fmt.Errorf("unsupported format %q (use mkv or mp4)", format)
	}
}
