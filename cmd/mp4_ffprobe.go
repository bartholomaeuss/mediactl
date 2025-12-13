package cmd

import "github.com/spf13/cobra"

// mp4FFProbeCmd reuses the FFprobe exporter for MP4 files.
var mp4FFProbeCmd = &cobra.Command{
	Use:   "ffprobe [directory]",
	Short: "Export ffprobe JSON sidecars for MP4 files",
	Long: `Runs ffprobe for each *.mp4 file inside the specified directory (current directory by default)
and writes a .json sidecar next to the MP4 file.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) == 1 {
			dir = args[0]
		}

		processed, err := walkFFProbe(dir, ".mp4", cmd, ffprobeRunner)
		if err != nil {
			return err
		}

		if processed == 0 {
			cmd.Printf("No MP4 files found in %s\n", dir)
		}

		return nil
	},
}

func init() {
	mp4Cmd.AddCommand(mp4FFProbeCmd)
}
