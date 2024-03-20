package errors

import (
	"go.uber.org/zap"
)

func LogError(method string, err error, withStack bool, logger *zap.Logger) {
	logger = logger.With(zap.String("method", method))

	if withStack {
		logger.Error("error", zap.String("error message", err.Error()), zap.Stack("stack"))
	} else {
		logger.Error("error", zap.String("error message", err.Error()))
	}
}
