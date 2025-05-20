// validate_chirp.go
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

// ValidateChirp handler that handles json reqs and resp!
func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	// consts
	const maxMessageLimit = 140

	// HTTP method check
	if req.Method != "POST" {
		// helper to insert error msg + 405 invalid method status code
		WriteJSONError(w, "Chirp must be POSTed", http.StatusMethodNotAllowed)
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
		WriteJSONError(w, "Chirp is empty", http.StatusBadRequest)
		return // early return
	}

	// decode check
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value

		log.Printf("Error decoding parameters: %s", err) // log msg with err
		// helper to insert error msg + 400 bad req status code
		WriteJSONError(w, "Something went wrong", http.StatusBadRequest)
		return // early return
	}

	// reqBody is now successfully populated

	// check chirp empty
	if len(reqBody.Body) == 0 {
		// helper to insert error msg + 400 bad req status code
		WriteJSONError(w, "Chirp is empty", http.StatusBadRequest)
		return // early return
	}

	// check chirp too long
	if len(reqBody.Body) > maxMessageLimit {
		// helper to insert error msg + 400 bad req status code
		WriteJSONError(w, "Chirp is too long", http.StatusBadRequest)
		return // early return
	}

	// clean the request body
	bodyClean := cleanProfanity(reqBody.Body)

	// json response to
	respBody := JsonResponse{CleanedBody: bodyClean} // set resp to valid as req is successful

	// helper to insert body response + 200 ok status code
	WriteJSONResponse(w, respBody, http.StatusOK)
}

// HELPER FUNCS

// RESPONSE helper to clean profanity before passing payload to response
func cleanProfanity(body string) string {
	// split the body
	words := strings.Split(body, " ") // split to list by space " "

	// profanity list
	profanityList := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	// check each word in body
	for i, word := range words {
		// lowercase the word
		wordLower := strings.ToLower(word)

		// compare lowercase word with lowercase profanity
		for _, profanity := range profanityList {

			// if it matches any of the profane words
			if wordLower == profanity {
				words[i] = "****" // replace it
				break             // stop checking other profanity, already matched
			}
		}
	}

	// rejoin the cleaned slice
	cleanedWords := strings.Join(words, " ")

	// return the cleaned slice of words
	return cleanedWords
}
