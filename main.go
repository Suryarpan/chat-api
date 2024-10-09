package main

import (
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main()  {

	loggerConf := zap.NewProductionConfig()
	loggerConf.EncoderConfig.TimeKey = "timestamp"
	loggerConf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConf.Build()
	if err != nil {
		log.Fatal("Could not acquire a logger")
	}
	sugar := logger.Sugar()
	log.Println()
}