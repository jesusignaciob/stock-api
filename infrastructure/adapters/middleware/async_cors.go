// async_cors.go provides an asynchronous CORS middleware for Gin.
// This middleware sets the appropriate CORS headers for allowed origins
// and handles preflight (OPTIONS) requests asynchronously with a timeout.

package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AsyncCORSMiddleware returns a Gin middleware handler that manages CORS headers asynchronously.
// Parameters:
// - allowedOrigins: a slice of allowed origin strings.
//
// The middleware:
// - Checks if the request's Origin header is in the allowedOrigins list.
// - Sets the appropriate CORS headers if the origin is allowed.
// - Handles preflight (OPTIONS) requests by aborting with HTTP 204 No Content.
// - Runs header setting and preflight handling in a goroutine to avoid blocking the main thread.
// - Uses a timeout to prevent hanging if the goroutine takes too long.
// - Aborts with HTTP 500 if the operation times out, or returns if the client cancels the request.
func AsyncCORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	originSet := make(map[string]struct{})
	for _, o := range allowedOrigins {
		originSet[o] = struct{}{}
	}

	return func(c *gin.Context) {
		// Use a goroutine for potentially blocking operations
		done := make(chan struct{})
		go func() {
			defer close(done)

			if origin := c.Request.Header.Get("Origin"); origin != "" {
				if _, exists := originSet[origin]; exists {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}

			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers",
				"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods",
				"POST, OPTIONS, GET, PUT, DELETE, PATCH")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}()

		select {
		case <-done:
			c.Next()
		case <-time.After(50 * time.Millisecond): // Timeout for CORS operations
			c.AbortWithStatus(http.StatusInternalServerError)
		case <-c.Request.Context().Done(): // If the client cancels
			return
		}
	}
}
