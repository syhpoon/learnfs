// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package mkfs

import (
	"context"
	"log"

	"learnfs/pkg/device"
	"learnfs/pkg/fs"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "mkfs {device}",
	Short: "Create a new filesystem on a given device",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dev, err := device.NewUring(context.Background(), args[0])
		if err != nil {
			log.Fatalf("failed to open device: %v", err)
		}

		if err := fs.Create(dev); err != nil {
			log.Fatalf("failed to create filesystem: %v", err)
		}

		dev.Close()
	},
}
