package logger

import "go.uber.org/zap"

func Load() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	return sugar

}
