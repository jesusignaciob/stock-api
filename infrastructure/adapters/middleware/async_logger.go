// middleware/logger.go
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// stringError implements the error interface for a string message
type stringError struct {
	msg string
}

func (e *stringError) Error() string {
	return e.msg
}

type concurrentLogger struct {
	logChan chan logEntry
	zap     *zap.Logger
}

type logEntry struct {
	status  int
	method  string
	path    string
	latency time.Duration
	ip      string
	errors  []error
}

func AsyncLogger(zapLogger *zap.Logger) gin.HandlerFunc {
	cl := &concurrentLogger{
		logChan: make(chan logEntry, 100), // Buffer para logs
		zap:     zapLogger,
	}

	// Worker único para escribir logs (evita contención de I/O)
	go cl.logWorker()
	return func(c *gin.Context) {
		start := time.Now()

		// Procesar la solicitud
		c.Next()

		// Convert []string to []error
		stringErrors := c.Errors.Errors()
		errorSlice := make([]error, len(stringErrors))
		for i, errStr := range stringErrors {
			errorSlice[i] = &stringError{msg: errStr}
		}

		entry := logEntry{
			status:  c.Writer.Status(),
			method:  c.Request.Method,
			path:    c.Request.URL.Path,
			latency: time.Since(start),
			ip:      c.ClientIP(),
			errors:  errorSlice,
		}

		select {
		case cl.logChan <- entry:
		case <-time.After(100 * time.Millisecond): // Timeout para evitar bloqueos
			log.Println("Logger buffer lleno, descartando entrada")
		}
	}
}

func (cl *concurrentLogger) logWorker() {
	for entry := range cl.logChan {
		// Usar logger estructurado
		cl.zap.Info("request",
			zap.String("method", entry.method),
			zap.String("path", entry.path),
			zap.Int("status", entry.status),
			zap.Duration("latency", entry.latency),
			zap.String("ip", entry.ip),
			zap.Errors("errors", entry.errors),
		)
	}
}
