package extensions

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func init() {
	encoderConfiguration := zap.NewDevelopmentEncoderConfig()
	encoderConfiguration.TimeKey = "timestamp"
	encoderConfiguration.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfiguration.EncodeLevel = zapcore.CapitalColorLevelEncoder

	configuration := zap.NewDevelopmentConfig()
	configuration.EncoderConfig = encoderConfiguration

	var err error
	logger, err = configuration.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

}

func Info(message string, fields ...zap.Field) {
	logger.Info("\n"+Green+message+Reset+"\n\n", fields...)
}

func Debug(message string, fields ...zap.Field) {
	logger.Debug("\n"+Yellow+message+Reset+"\n\n", fields...)
}

func Error(message string, fields ...zap.Field) {
	// message = Red + message + Reset
	logger.Error("\n"+Red+message+Reset+"\n\n", fields...)
}
