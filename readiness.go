// readiness.go
package main

import "net/http"

// HANDLERS
// Readiness handler gives server /healthz system endpoint
func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	// set header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// write the header (status code or "health")
	w.WriteHeader(http.StatusOK)

	// write the response body as "OK"
	w.Write([]byte(http.StatusText(http.StatusOK)))
	// instead of hardcoding, use the OK status var
}
