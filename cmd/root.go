/*
Copyright © 2025 Bartholomaeuss
*/
// Package cmd defines the Cobra command tree for mediactl.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the top-level command used when no subcommand is provided.
var rootCmd = &cobra.Command{
	Use:   "mediactl",
	Short: "CLI utilities for media inspection and sidecar workflows",
	Long: `mediactl is a focused CLI for inspecting media libraries and generating
JSON sidecars for downstream tooling.

Start with the sidecar namespace to export ffprobe snapshots or build custom
commands around your media workflow.`,
}

// Execute builds the command tree and runs the CLI, exiting on command failure.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init registers subcommands attached to the root command.
func init() {
}
