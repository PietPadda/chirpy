// validate_chirp.go
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// ValidateChirp handler that handles json reqs and resps!
func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	// consts
	const maxMessageLimit = 140

	// HTTP method check
	if req.Method != "POST" {
		// helper to insert error msg + 405 invalid method status code
		writeJSONError(w, "Chirp must be POSTed", http.StatusMethodNotAllowed)
		return // early return
	}

	// json request from client
	var reqBody JsonRequest

	// create json req body decoder
	decoder := json.NewDecoder(req.Body)

	// close on exit to prevent mem leak
	defer req.Body.Close()

	// decode the req body
	err := decoder.Decode(&reqBody)

	// request body missing edge case check (before general error check)
	if err == io.EOF { // end of file
		log.Printf("Error empty request body: %s", err) // log msg with err
		// helper to insert error msg + 400 bad req status code
		writeJSONError(w, "Chirp is empty", http.StatusBadRequest)
		return // early return
	}

	// decode check
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value

		log.Printf("Error decoding parameters: %s", err) // log msg with err
		// helper to insert error msg + 400 bad req status code
		writeJSONError(w, "Something went wrong", http.StatusBadRequest)
		return // early return
	}

	// reqBody is now successfully populated

	// check chirp empty
	if len(reqBody.Body) == 0 {
		// helper to insert error msg + 400 bad req status code
		writeJSONError(w, "Chirp is empty", http.StatusBadRequest)
		return // early return
	}

	// check chirp too long
	if len(reqBody.Body) > maxMessageLimit {
		// helper to insert error msg + 400 bad req status code
		writeJSONError(w, "Chirp is too long", http.StatusBadRequest)
		return // early return
	}

	// json response to
	respBody := JsonResponse{Valid: true} // set resp to valid as req is successful

	// helper to insert body response + 200 ok status code
	writeJSONResponse(w, respBody, http.StatusOK)
}

// ERROR helper to make the API much more DRY
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	respError := JsonResponseError{Error: message}

	// marshal the go responserror to json, removing whitespaces
	dat, err := json.Marshal(respError)

	// marhsal check
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
func writeJSONResponse(w http.ResponseWriter, payload interface{}, statusCode int) {
	// marshal the go response to json, removing whitespaces
	dat, err := json.Marshal(payload)

	// marhsal check
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
