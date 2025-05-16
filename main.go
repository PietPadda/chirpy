package main

import (
	// std go libraries
	"log"      // for err loggin
	"net/http" // http protocol
)

func main() {
	// create server mux for routing http requests
	mux := http.NewServeMux()

	// apply a fileserver handler to mux
	mux.Handle("/", http.FileServer(http.Dir(".")))
	// http.Dir(".") -  tells server to run files "here"
	// http.FileServer(...) - looks for the index.html
	// mux.Handle("/", ...) -- server handle all requests

	// create Server struct for config
	server := &http.Server{ //ptr is more efficient than new copy
		Addr:    ":8080", //server listens to port 8080 for all incoming requests
		Handler: mux,     // mux will "handle" our http request
	}

	// start the server
	err := server.ListenAndServe() // "listens" to http requests on addr and let's mux handle them

	// server start check
	if err != nil {
		// log the err and terminate server
		log.Fatal(err)
	}
}
