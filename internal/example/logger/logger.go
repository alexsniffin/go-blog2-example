package logger

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/alexsniffin/go-blog2-example/internal/example/models"
)

func NewLogger(cfg models.Logger) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return zerolog.Logger{}, err
	}

	logger := zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger()

	// use pretty logging
	logger = logger.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})

	return logger, nil
}
