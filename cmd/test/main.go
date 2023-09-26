package main

import (
	"kprg/internal/logger"
	"os"
)

func main() {

	logger := logger.NewLogger(os.Stdout)

	logger.Info("Info")
	logger.Warn("Warn")
	logger.Error("Error")

}
