package slog

import (
	"io"
	"log/slog"
)

type Slog struct {
	logger *slog.Logger
}

func NewSlog(writer io.Writer) *Slog {

	logger :=
		&Slog{
			logger: slog.New(slog.NewTextHandler(writer, nil)),
		}

	return logger

}

func (l *Slog) Info(msg string) {

	l.logger.Info(msg)

}

func (l *Slog) Warn(msg string) {

	l.logger.Warn(msg)

}

func (l *Slog) Error(msg string) {

	l.logger.Error(msg)

}
