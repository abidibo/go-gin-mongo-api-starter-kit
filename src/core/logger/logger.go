package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var once sync.Once

// Init initializes a thread-safe singleton logger
// This would be called from a main method when the application starts up
// This function would ideally, take zap configuration, but is left out
// in favor of simplicity using the example logger.
func init() {
	// once ensures the singleton is initialized only once
	once.Do(func() {
		consoleDebugging := zapcore.Lock(os.Stdout)
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
		core := zapcore.NewCore(consoleEncoder, consoleDebugging, zapcore.DebugLevel)
		logger := zap.New(core, zap.AddCaller())
		defer logger.Sync() // flushes buffer, if any
		zap.ReplaceGlobals(logger)
	})
}
