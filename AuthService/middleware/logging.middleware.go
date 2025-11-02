package middleware

import (
	"log"
	"net/http"
	"time"
)

func RequestLoggingMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Record the start time

		// -- 1. PRE-PROCESSING (Logging) --
		log.Printf(
			"[%s] STARTED: %s %s from %s",
			start.Format("2006/01/02 15:04:05"),
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
		)

		// -- 2. EXECUTE THE NEXT HANDLER --
		next.ServeHTTP(w, r) // Call the original handler function

		// -- 3. POST-PROCESSING (Logging duration) --
		log.Printf(
			"[%s] COMPLETED: %s %s in %v",
			time.Now().Format("2006/01/02 15:04:05"),
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
