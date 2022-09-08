// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package mount

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/syhpoon/learnfs/pkg/fuse"
	"github.com/syhpoon/learnfs/pkg/logger"
	"github.com/syhpoon/learnfs/proto"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var flagServerAddress string

var Cmd = &cobra.Command{
	Use:   "mount {mount-point}",
	Short: "Mount learnfs from the given server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		logger.Init()

		mount := args[0]
		conn, err := grpc.DialContext(ctx, flagServerAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Fatal().Err(err).Str("server", flagServerAddress).
				Msg("failed to dial learnfs server")
		}

		defer conn.Close()

		cl := proto.NewLearnFSClient(conn)
		ses, err := fuse.NewSession(mount, cl)
		if err != nil {
			log.Fatal().Err(err).Msg("fuse fatal error")
		}

		go ses.Start(ctx)
		defer ses.Stop()

		wait(ctx, cancel)
	},
}

func wait(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

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
}

func init() {
	Cmd.Flags().StringVarP(&flagServerAddress, "server", "s", "127.0.0.1:1234",
		"learnfs server address")
}
