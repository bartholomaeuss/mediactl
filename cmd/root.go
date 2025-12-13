/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mediactl",
	Short: "Utilities for inspecting media libraries",
	Long: `mediactl is a lightweight helper for working with Matroska (MKV) collections.

Use the mkv namespace to export ffprobe JSON sidecars or build custom subcommands
for other media workflows.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
