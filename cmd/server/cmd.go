// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package server

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/syhpoon/learnfs/pkg/device"
	"github.com/syhpoon/learnfs/pkg/fs"
	"github.com/syhpoon/learnfs/pkg/logger"
	"github.com/syhpoon/learnfs/pkg/server"

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

		dev, err := device.NewUring(context.Background(), args[0])
		if err != nil {
			log.Fatal().Err(err).Msg("failed to open device")
		}

		fsObj, err := fs.Load(dev)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to load filesystem")
		}

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go fsObj.RunFlusher(ctx, wg)

		go func() {
			if err := server.RunServer(ctx, server.Params{
				Listen: flagListen,
				Fs:     fsObj,
			}); err != nil {
				log.Fatal().Err(err).Msg("failed to run learnfs server")
			}
		}()

		wait(ctx, cancel, wg)
		_ = fsObj.Shutdown()
	},
}

func wait(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	c := make(chan os.Signal, 1)

	signal.Notify(c,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGABRT,
		syscall.SIGPIPE,
		syscall.SIGBUS,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGQUIT)

LOOP:
	for {
		select {
		case <-ctx.Done():
			break

		case sig := <-c:
			switch sig {
			case os.Interrupt, syscall.SIGTERM, syscall.SIGABRT,
				syscall.SIGPIPE, syscall.SIGBUS, syscall.SIGUSR1, syscall.SIGUSR2,
				syscall.SIGQUIT:

				break LOOP
			}
		}
	}

	cancel()
	wg.Wait()
}

func init() {
	Cmd.Flags().StringVarP(&flagListen, "listen", "l", ":1234", "Server listen address")
}
