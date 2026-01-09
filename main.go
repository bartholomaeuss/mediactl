/*
Copyright © 2025 Bartholomaeuss
*/
// Package main wires the mediactl CLI entrypoint.
package main

import "mediactl/cmd"

// main boots the CLI command tree and exits on failure.
func main() {
	cmd.Execute()
}
