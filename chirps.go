// chirps.go
package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PietPadda/chirpy/internal/database"
	"github.com/google/uuid"
)

// CreateChirp handler that creates a chirp (keep ValidateChirp logic)
func (apiCfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	// consts
	const maxMessageLimit = 140

	// HTTP method check
	if req.Method != "POST" {
		// helper to insert error msg + 405 invalid method status code
		WriteJSONError(w, "Chirp must be POSTed", http.StatusMethodNotAllowed)
		return // early return
	}

	// json request from client
	var reqBody JsonChirpRequest

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

	// create chirp
	newChirp, err := apiCfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   bodyClean,      // add the profanity cleaned chirp body
		UserID: reqBody.UserID, // get user_id for fk from the json request body
	})

	// create chirp check
	if err != nil {
		log.Printf("Error creating chirp: %s", err) // log msg with err
		// helper to insert error msg + 400 bad req status code
		WriteJSONError(w, "Error occurred creating new chirp", http.StatusBadRequest)
		return // early return
	}

	// json response payload
	respChirp := JsonChirpResponse{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}

	// helper to insert body response + 201 created status code
	WriteJSONResponse(w, respChirp, http.StatusCreated)
}

// GetChirps handler that returns all chirps!
func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	// apiConfig check
	if apiCfg == nil {
		// handle gracefully
		log.Printf("Internal server error: apiCfg is nil") // msg to server admin
		// send msg to client code 500
		WriteJSONError(w, "Internal server configuration error", http.StatusInternalServerError)
		return // stop processing req
	}

	// HTTP method check
	if req.Method != "GET" {
		// helper to insert error msg + 405 invalid method status code
		WriteJSONError(w, "Chirps must be GETted", http.StatusMethodNotAllowed)
		return // early return
	}

	// get (all) chirps!
	dbChirps, err := apiCfg.db.GetChirps(req.Context())

	// get chirps check
	if err != nil {
		log.Printf("Error getting all chirps: %s", err) // log msg with err
		// helper to insert error msg + 400 bad req status code
		WriteJSONError(w, "Error occurred getting all chirps", http.StatusBadRequest)
		return // early return
	}

	// init a pre-allocated slice (not nil or zero!) all fields 0'd to store all chirps
	chirpResponses := make([]JsonChirpResponse, len(dbChirps))
	// pre-alloc mem to # of chirps, no less, no more -- mem efficient!

	// loop thru chirps and add each as an element
	for i, dbChirp := range dbChirps {
		// make a response for each chirp at index i
		chirpResponses[i] = JsonChirpResponse{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}

	// helper to insert body response + 200 OK status code
	WriteJSONResponse(w, chirpResponses, http.StatusOK)
}

// GetChirp handler that returns one chirp by id
func (apiCfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	// apiConfig check
	if apiCfg == nil {
		// handle gracefully
		log.Printf("Internal server error: apiCfg is nil") // msg to server admin
		// send msg to client code 500
		WriteJSONError(w, "Internal server configuration error", http.StatusInternalServerError)
		return // stop processing req
	}

	// HTTP method check
	if req.Method != "GET" {
		// helper to insert error msg + 405 invalid method status code
		WriteJSONError(w, "Chirp must be GETted", http.StatusMethodNotAllowed)
		return // early return
	}

	// get chirp id from api endpoint path string
	chirpIDStr := req.PathValue("chirpID")

	// conv str to UUID to use GetChirp
	chirpUUID, err := uuid.Parse(chirpIDStr)

	// uuid conv check
	if err != nil {
		log.Printf("Error getting chirp ID: %s", err) // log msg with err
		// helper to insert error msg + 400 bad req status code
		WriteJSONError(w, "Invalid chirp ID format", http.StatusBadRequest)
		return // early return
	}

	// get (one) chirps
	dbChirp, err := apiCfg.db.GetChirp(req.Context(), chirpUUID)

	// Check if this is a 404 "not found" error
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Error could not find chirp: %s", err) // log msg with err
		// helper to insert error msg + 404 not found status code
		WriteJSONError(w, "Chirp not found", http.StatusNotFound)
		return
	}

	// get chirps check
	if err != nil {
		log.Printf("Error getting chirp: %s", err) // log msg with err
		// helper to insert error msg + 500 internal error status code
		WriteJSONError(w, "Error occurred getting chirp", http.StatusInternalServerError)
		return // early return
	}

	// build the chirp response
	chirpResp := JsonChirpResponse{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	// helper to insert body response + 200 OK status code
	WriteJSONResponse(w, chirpResp, http.StatusOK)
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
