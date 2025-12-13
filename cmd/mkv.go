/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// mkvCmd represents the mkv command
var mkvCmd = &cobra.Command{
	Use:   "mkv",
	Short: "Namespace for MKV-related operations",
	Long: `Namespace for Matroska (MKV) commands.
All functionality for inspecting or manipulating MKV media lives in subcommands beneath this node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Keep mkv as a namespace; show help when invoked directly.
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(mkvCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mkvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mkvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
