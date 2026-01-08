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
	Short: "CLI utilities for media inspection and sidecar workflows",
	Long: `mediactl is a focused CLI for inspecting media libraries and generating
JSON sidecars for downstream tooling.

Start with the sidecar namespace to export ffprobe snapshots or build custom
commands around your media workflow.`,
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
