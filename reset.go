// reset.go
package main

import (
	"net/http"
	"strconv"
)

// MetricsResets handler that resets the fileserver hits!
// apply cfg receiver to access fileserverHits
func (cfg *apiConfig) handlerMetricsReset(w http.ResponseWriter, req *http.Request) {
	// set header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// not required here, but good to be explicit

	// write the header (status code or "health")
	w.WriteHeader(http.StatusOK)
	// tells server to send "200 OK" BEFORE writing the response body
	// can skip, will happen implicitly

	// get fileserverhits
	cfg.fileserverHits.Store(0)        // in32 doesn't have .Reset(), but we can "store" it as 0 :)
	x := cfg.fileserverHits.Load()     // safely load Int32 the RESETed number of hits
	serverhits := strconv.Itoa(int(x)) // convert to string, make int as was int32

	// write the response body
	w.Write([]byte("Hits reset to " + serverhits + "\n"))
}
