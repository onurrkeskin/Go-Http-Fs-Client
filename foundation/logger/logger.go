package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(application string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()

	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.EpochNanosTimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]any{
		"application": application,
	}

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}
