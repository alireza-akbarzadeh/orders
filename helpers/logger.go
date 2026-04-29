package helpers

import (
	"sync"

	"go.uber.org/zap"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Get returns a singleton zap.Logger instance
func Get() *zap.Logger {
	once.Do(func() {
		var err error
		log, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	})
	return log
}
