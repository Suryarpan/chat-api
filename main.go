package main

import (
	"fmt"
	"log/slog"
	"os"
)

func loggerConfig() *slog.HandlerOptions {
	level, ok := os.LookupEnv("CHAT_API_CONFIG")
	if !ok {
		level = "prod"
	}

	var logLevel slog.Level
	switch level {
	case "dev":
		logLevel = slog.LevelDebug
	case "prod":
		logLevel = slog.LevelInfo
	default:
		panic(fmt.Sprintf("could not understand log level: %s\n", logLevel))
	}
	fmt.Printf("Log level set to %s\n", logLevel)
	return &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, loggerConfig()))
	slog.SetDefault(logger)

	slog.Info("starting chat api server")
	slog.Debug("using debug level")
	slog.Info("using info level")
}
