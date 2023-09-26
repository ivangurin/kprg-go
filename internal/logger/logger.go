package logger

import (
	"io"
	"kprg/internal/logger/slog"
)

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

func NewLogger(writer io.Writer) Logger {

	return slog.NewSlog(writer)

}
