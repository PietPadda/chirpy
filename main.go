package main

import (
	// std go libraries
	"log"      // for err loggin
	"net/http" // http protocol
)

func main() {
	// set constants
	const filepathRoot = "." // used constant
	const port = "8080"

	// create server mux for routing http requests
	mux := http.NewServeMux()

	// apply a fileserver handler to mux
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))
	// http.Dir(".") -  tells server to run files "here"
	// http.FileServer(...) - looks for the index.html
	// mux.Handle("/app/", ...) -- server handle all requests
	// stripping the /app/ to just . -- /app/ is just there for cleaner url

	// register handlerReadiness, using /healthz system endpoint
	mux.HandleFunc("/healthz", handlerReadiness)
	// /healthz, because "system endpoint" convention!

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
