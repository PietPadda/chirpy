// metrics.go
package main

import (
	"fmt"
	"net/http"
)

// Metrics handler that returns the fileserver hits!
// apply cfg receiver to access fileserverHits
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	// set header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// write the header (status code or "health")
	w.WriteHeader(http.StatusOK)
	// tells server to send "200 OK" BEFORE writing the response body
	// can skip, will happen implicitly

	// get fileserverhits
	x := cfg.fileserverHits.Load() // safely load the number of hits

	// write the response body
	// use backticks to handle multiline printing
	fmt.Fprintf(w, `
	<html>
      <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
      </body>
    </html>
	`, x)
}

// MIDDLEWARE
// metric increment fileserverhits middleware
// uses apiConfig as method receiver
// wraps the input handler, count each request before passing control along
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc( // return a new func
		func(w http.ResponseWriter, r *http.Request) { // http.Handler type func
			cfg.fileserverHits.Add(1) // incr the server hits atomically for EACH request (concur safety)
			next.ServeHTTP(w, r)      // pass req to the next handler in the chain
		}, // trailing comma required in last arg of multi-line call
	)
}
