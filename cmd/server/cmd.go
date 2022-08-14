// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package server

import (
	"context"

	"learnfs/pkg/device"
	"learnfs/pkg/fs"
	"learnfs/pkg/logger"
	"learnfs/pkg/server"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var flagListen string

var Cmd = &cobra.Command{
	Use:   "server {device}",
	Short: "Run a learnfs server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		logger.Init()

		dev, err := device.NewBlockUring(context.Background(), args[0])
		if err != nil {
			log.Fatal().Err(err).Msg("failed to open device")
		}

		fsObj, err := fs.Load(dev)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to load filesystem")
		}
		defer fsObj.Shutdown()

		if err := server.RunServer(ctx, server.Params{
			Listen: flagListen,
			Fs:     fsObj,
		}); err != nil {
			log.Fatal().Err(err).Msg("failed to run learnfs server")
		}
	},
}

func init() {
	Cmd.Flags().StringVarP(&flagListen, "listen", "l", ":1234", "Server listen address")
}
