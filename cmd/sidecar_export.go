/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import (
	"fmt"
	"strings"

	"mediactl/core"
	"mediactl/infra"

	"github.com/spf13/cobra"
)

// sidecarExportCmd exports ffprobe JSON sidecars for a specific media format.
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
		if len(args) == 1 {
			dir = args[0]
		}

		format := strings.ToLower(strings.TrimSpace(sidecarExportFormat))
		ext, err := sidecarFormatToExtension(format)
		if err != nil {
			return err
		}

		processed, err := core.WalkFFProbe(dir, ext, cmd.OutOrStdout(), infra.ExecuteFFProbe)
		if err != nil {
			return err
		}

		if processed == 0 {
			cmd.Printf("No %s files found in %s\n", strings.ToUpper(format), dir)
		}

		return nil
	},
}

var sidecarExportFormat string

func init() {
	sidecarCmd.AddCommand(sidecarExportCmd)
	sidecarExportCmd.Flags().StringVar(&sidecarExportFormat, "format", "", "Media format to export (mkv or mp4)")
	_ = sidecarExportCmd.MarkFlagRequired("format")
}

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
