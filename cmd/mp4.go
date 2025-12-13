package cmd

import "github.com/spf13/cobra"

// mp4Cmd represents the mp4 command
var mp4Cmd = &cobra.Command{
	Use:   "mp4",
	Short: "Namespace for MP4-related operations",
	Long: `Namespace for MP4 commands.
All functionality for inspecting or manipulating MP4 media lives in subcommands beneath this node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(mp4Cmd)
}
