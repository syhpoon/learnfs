// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package main

import (
	"os"

	"learnfs/cmd/info"
	"learnfs/cmd/mkfs"
	"learnfs/cmd/mount"
	"learnfs/cmd/server"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "learnfs",
	Short: "learnfs command line tool",
}

func init() {
	rootCmd.AddCommand(info.Cmd)
	rootCmd.AddCommand(mkfs.Cmd)
	rootCmd.AddCommand(mount.Cmd)
	rootCmd.AddCommand(server.Cmd)
}
