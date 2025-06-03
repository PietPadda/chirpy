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

	"github.com/PietPadda/chirpy/internal/auth"
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

	// authenticate before decoding request
	token, err := auth.GetBearerToken(req.Header) // get the bearer's token

	// get token check
	if err != nil {
		log.Printf("Error getting bearer token: %s", err) // log msg with err
		// helper to insert error msg + 401 unauthorized status code
		WriteJSONError(w, "Unauthorized access", http.StatusUnauthorized)
		return // early return
	}

	// validate the JWT token after getting bearer's token
	uuidJWTValidated, err := auth.ValidateJWT(token, apiCfg.serverKey) // pass in tokenstring and server secret

	// jwt validation check
	if err != nil {
		log.Printf("Error validating JWT token: %s", err) // log msg with err
		// helper to insert error msg + 401 unauthorized status code
		WriteJSONError(w, "Unauthorized access", http.StatusUnauthorized)
		return // early return
	}

	// create json req body decoder
	decoder := json.NewDecoder(req.Body)

	// close on exit to prevent mem leak
	defer req.Body.Close()

	// decode the req body
	err = decoder.Decode(&reqBody)

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
		Body:   bodyClean,        // add the profanity cleaned chirp body
		UserID: uuidJWTValidated, // get user_id from the VALIDATED JWT!
	}) // we ignore the request's userid and ONLY use the VALIDATED userid!

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

// DeleteChirp handler that deletes a chirp
func (apiCfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	// HTTP method check
	if req.Method != "DELETE" {
		// helper to insert error msg + 405 invalid method status code
		WriteJSONError(w, "Chirp must be DELETEted", http.StatusMethodNotAllowed)
		return // early return
	}

	// no json request from client, get chirp_api from endpoint path
	urlSplit := strings.Split(req.URL.Path, "/") // split by /

	// valid path check
	if len(urlSplit) != 4 { // we explicitly only tolerate /api/chirps/{chirpID}, nothing else
		log.Printf("Error api endpoint path incorrect length: %v", len(urlSplit)) // log msg with err
		// helper to insert error msg + 400 bad request status code
		WriteJSONError(w, "Invalid API endpoint path", http.StatusBadRequest)
		return // early return
	}

	// check empty path chirp id check
	if urlSplit[3] == "" {
		log.Printf("Error api endpoint path chirp id empty: %v", len(urlSplit)) // log msg with err
		// helper to insert error msg + 400 bad request status code
		WriteJSONError(w, "Missing chirp id in endpoint", http.StatusBadRequest)
		return // early return
	}

	// know path is /api/chirps/${chirpID}, so get index 3 (1st / counts)
	chirpID := urlSplit[3] // last item in split path string

	// conv chirpID to int then uuid
	chirpUUID, err := uuid.Parse(chirpID)

	// chirp id to uuid check
	if err != nil {
		log.Printf("Error conv chirp id to uuid: %s", err) // log msg with err
		// helper to insert error msg + 400 bad request status code
		WriteJSONError(w, "Invalid Chirp id", http.StatusBadRequest)
		return // early return
	}

	// authenticate before decoding request
	token, err := auth.GetBearerToken(req.Header) // get the bearer's token

	// get token check
	if err != nil {
		log.Printf("Error getting bearer token: %s", err) // log msg with err
		// helper to insert error msg + 401 unauthorized status code
		WriteJSONError(w, "Unauthorized access", http.StatusUnauthorized)
		return // early return
	}

	// validate the JWT token after getting bearer's token
	uuidJWTValidated, err := auth.ValidateJWT(token, apiCfg.serverKey) // pass in tokenstring and server secret

	// jwt validation check
	if err != nil {
		log.Printf("Error validating JWT token: %s", err) // log msg with err
		// helper to insert error msg + 401 unauthorized status code
		WriteJSONError(w, "Unauthorized access", http.StatusUnauthorized)
		return // early return
	}

	// get chirp author id
	uuidChirpAuthor, err := apiCfg.db.GetUserIDByChirpID(req.Context(), chirpUUID)

	// get chirp author id check
	if err != nil {
		log.Printf("Error chirp not found: %s", err) // log msg with err
		// helper to insert error msg + 404 not found status code
		WriteJSONError(w, "Chirp not found", http.StatusNotFound)
		return // early return
	}

	// confirm chirp author id matches the JWT validated id
	if uuidChirpAuthor != uuidJWTValidated {
		log.Printf("Error chirp author id (%v) doesn't match JWT validated user id (%v)", uuidChirpAuthor, uuidJWTValidated) // log msg with err)", err) // log msg with err
		// helper to insert error msg + 403 forbidden status code
		WriteJSONError(w, "Unauthorized access", http.StatusForbidden)
		return // early return
	}

	// proceed to delete the chirp
	deletedChirp, err := apiCfg.db.DeleteChirp(req.Context(), chirpUUID)

	// delete chirp check
	if err != nil {
		// handle gracefully
		log.Printf("Error could not delete chirp: %s", err) // msg to server admin
		// send msg to client code 500
		WriteJSONError(w, "Server couldn't delete chirp", http.StatusInternalServerError)
		return // stop processing req
	}

	// write to server and client that chirp deleted
	log.Printf("Chirp has been deleted: ID = %s", deletedChirp.ID) // log msg with err
	w.WriteHeader(http.StatusNoContent)                            // status code 204 to client
}

// GetChirps handler that returns all chirps!
func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	// apiConfig check: Always validate essential dependencies first.
	if apiCfg == nil {
		log.Printf("Internal server error: apiCfg is nil") // Log to server admin
		WriteJSONError(w, "Internal server configuration error", http.StatusInternalServerError)
		return // Stop processing
	}

	// HTTP method check: Ensure the correct HTTP verb is used.
	if req.Method != http.MethodGet { // Use http.MethodGet constant for clarity
		WriteJSONError(w, "Chirps must be GETted", http.StatusMethodNotAllowed)
		return // Early return
	}

	authorIDStr := req.URL.Query().Get("author_id") // Get optional query parameter
	var dbChirps []database.Chirp                   // Initialize an empty chirp slice
	var err error                                   // Declare error variable once

	// Handle requests with an author_id query parameter
	if len(authorIDStr) > 0 { // if a char is input
		authorUUID, parseErr := uuid.Parse(authorIDStr)

		// uuid check
		if parseErr != nil {
			log.Printf("Error converting author ID '%s' to UUID: %s", authorIDStr, parseErr)
			WriteJSONError(w, "Invalid author ID format", http.StatusBadRequest)
			return
		}

		// get chirps from optional author id para,
		dbChirps, err = apiCfg.db.GetChirpsByAuthorID(req.Context(), authorUUID)

		// get chirps by author check
		if err != nil {
			// If no rows are found, it means the author either doesn't exist or has no chirps.
			// Returning 200 OK with an empty array is a common and user-friendly approach for no results.
			if errors.Is(err, sql.ErrNoRows) { // check if no records
				log.Printf("No chirps found for author ID: %s", authorUUID)
				WriteJSONResponse(w, []JsonChirpResponse{}, http.StatusOK) // Return empty array
				return
			}
			// For any other database error, it's an internal server issue.
			log.Printf("Error getting chirps for author ID %s: %s", authorUUID, err)
			WriteJSONError(w, "Failed to retrieve chirps by author", http.StatusInternalServerError)
			return
		}
		// otherwise, no author id is provided
	} else {
		// Handle requests without an author_id query parameter (get all chirps)
		dbChirps, err = apiCfg.db.GetChirps(req.Context())

		// get chirps check
		if err != nil {
			log.Printf("Error getting all chirps: %s", err)
			WriteJSONError(w, "Failed to retrieve all chirps", http.StatusInternalServerError)
			return
		}
	}

	// Transform database chirps into JSON response format
	chirpResponses := make([]JsonChirpResponse, len(dbChirps))
	for i, dbChirp := range dbChirps { // loop through field in each chirp
		chirpResponses[i] = JsonChirpResponse{ // then populate the response
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}

	// Send successful response
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
