package logger

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var once sync.Once

// Init initializes the global logger
func Init() {
	once.Do(func() {
		// Configure the global logger
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
		log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()

		// Set the global log level (optional)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a sub-logger with request-specific fields
		logger := log.With().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("remote_addr", r.RemoteAddr).
			Logger()

		// Log the incoming request
		logger.Info().Msg("Request started")

		// Create a response writer to capture the status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Call the next handler
		next.ServeHTTP(ww, r)

		// Log the request completion with the status code and duration
		logger.Debug().
			Int("status", ww.Status()).
			Dur("duration", time.Since(start)).
			Msg("Request completed")
	})
}
