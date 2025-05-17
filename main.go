package main

import (
	// std go libraries
	"fmt"         // for printing
	"log"         // for err loggin
	"net/http"    // http protocol
	"strconv"     // for conv itoa or atoi
	"sync/atomic" // allows safe incr + read of ints for goroutines
)

// STRUCTS
// stateful struct
type apiConfig struct {
	fileserverHits atomic.Int32
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

// MAIN
func main() {
	// set constants
	const filepathRoot = "." // used constant
	const port = "8080"

	// create server mux for routing http requests
	mux := http.NewServeMux()

	// create apiConfig instance
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{}, // explicitly set to 0
	}

	// apply a fileserver handler to mux
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))),
	))
	// http.Dir(".") -  tells server to run files "here"
	// http.FileServer(...) - looks for the index.html
	// mux.Handle("/app/", ...) -- server handle all requests
	// stripping the /app/ to just . -- /app/ is just there for cleaner url

	// REGISTER HANDLERS
	// register handlerReadiness, using /healthz system endpoint
	mux.HandleFunc("/healthz", handlerReadiness)
	// /healthz, because "system endpoint" convention!

	// register handlerMetrics, using /metrics system endpoint
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics) // register func that receives apiCfg
	// /metrics, no z as this is a conventional name!

	// register handlerMetricsReset, using /reset system endpoint
	mux.HandleFunc("/reset", apiCfg.handlerMetricsReset) // register func that receives apiCfg
	// /reset, no z as this is a conventional name!

	// create Server struct for config
	server := &http.Server{ //ptr is more efficient than new copy
		Addr:    ":" + port, //server listens to port 8080 for all incoming requests
		Handler: mux,        // mux will "handle" our http request
	}

	// print server running msg
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	// print msg before blocking with "Listen"

	// start the server
	err := server.ListenAndServe() // "listens" to http requests on addr and let's mux handle them
	// listen blocks the server to prevent ending main func

	// server start check
	if err != nil {
		// log the err and terminate server
		log.Fatal(err)
	}
}

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

// Metrics handler that returns the fileserver hits!
// apply cfg receiver to access fileserverHits
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	// set header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// write the header (status code or "health")
	w.WriteHeader(http.StatusOK)
	// tells server to send "200 OK" BEFORE writing the response body
	// can skip, will happen implicitly

	// get fileserverhits
	x := cfg.fileserverHits.Load()     // safely load the number of hits
	serverhits := strconv.Itoa(int(x)) // convert to string, make int as was int32

	// write the response body
	w.Write([]byte(fmt.Sprintf("Hits: %v", serverhits)))
}

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
