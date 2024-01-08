package logging

import (
	"go.uber.org/zap"
	"sync"
)

var once sync.Once
var sugared zap.SugaredLogger

func NewLogger() zap.SugaredLogger {

	once.Do(func() {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		defer logger.Sync()

		sugared = *logger.Sugar()
	})

	return sugared
}
