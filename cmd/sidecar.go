/*
Copyright © 2025 Bartholomaeuss
*/
package cmd

import "github.com/spf13/cobra"

// sidecarCmd represents the sidecar namespace.
var sidecarCmd = &cobra.Command{
	Use:   "sidecar",
	Short: "Create and manage JSON sidecars for media files",
	Long: `Namespace for sidecar workflows, including exporting ffprobe JSON
and preparing metadata files that live alongside your media.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(sidecarCmd)
}
