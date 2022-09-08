// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package info

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/syhpoon/learnfs/pkg/device"
	"github.com/syhpoon/learnfs/pkg/fs"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "info {device}",
	Short: "Print filesystem info on a given device",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dev, err := device.NewUring(context.Background(), args[0])
		if err != nil {
			log.Fatalf("failed to open device: %v", err)
		}

		fsObj, err := fs.Load(dev)
		if err != nil {
			log.Fatalf("failed to load filesystem: %v", err)
		}

		sb := fsObj.Superblock()

		b, _ := json.MarshalIndent(sb, "", "  ")
		fmt.Printf("%+v\n", string(b))

		_ = fsObj.Shutdown()
	},
}
