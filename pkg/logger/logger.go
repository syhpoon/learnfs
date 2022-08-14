// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	marshalStack := func(err error) interface{} {
		return fmt.Errorf("%w", err)
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.ErrorStackMarshaler = marshalStack

	var writer io.Writer

	writer = zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}

	loggerBase := zerolog.New(writer).With().Stack().Timestamp().Logger()

	log.Logger = loggerBase.With().Caller().Logger()
}
