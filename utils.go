// utils.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// ERROR helper to make the API much more DRY
func WriteJSONError(w http.ResponseWriter, message string, statusCode int) {
	respError := JsonResponseError{Error: message}

	// marshal the go response error to json, removing whitespaces
	dat, err := json.Marshal(respError)

	// marshall check
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err) // log msg with err
		w.WriteHeader(http.StatusInternalServerError) // status code
		return                                        // early return
	}

	// respError is now successfully populated
	w.Header().Set("Content-Type", "application/json") // set header to json
	w.WriteHeader(statusCode)                          // status code
	w.Write(dat)                                       // write the response error
}

// RESPONSE helper to make the API much more DRY
// payload to allow ANY type of struct as input
func WriteJSONResponse(w http.ResponseWriter, payload interface{}, statusCode int) {
	// marshal the go response to json, removing whitespaces
	dat, err := json.Marshal(payload)

	// marshall check
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err) // log msg with err
		w.WriteHeader(http.StatusInternalServerError) // status code
		return                                        // early return
	}

	// respBody is now successfully populated

	// send the server response to client
	w.Header().Set("Content-Type", "application/json") // set header to json
	w.WriteHeader(statusCode)                          // status code
	w.Write(dat)                                       // write the response body
}
