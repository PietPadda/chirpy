// users.go
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/lib/pq"
)

// ValidateChirp handler that handles json reqs and resp!
// we use apiConfig as receiver to access the database!
func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	// apiConfig check
	if apiCfg == nil {
		// handle gracefully
		log.Printf("Internal server error: apiCfg is nil") // msg to server admin
		// send msg to client code 500
		WriteJSONError(w, "Internal server configuration error", http.StatusInternalServerError)
		return // stop processing req
	}

	// HTTP method check
	if req.Method != "POST" {
		// helper to insert error msg + 405 invalid method status code
		WriteJSONError(w, "User creation must be POSTed", http.StatusMethodNotAllowed)
		return // early return
	}

	// json request from client
	var reqEmail JsonRequestEmail

	// create json req body decoder
	decoder := json.NewDecoder(req.Body)

	// close on exit to prevent mem leak
	defer req.Body.Close()

	// decode the req body
	err := decoder.Decode(&reqEmail)

	// request body missing edge case check (before general error check)
	if err == io.EOF { // end of file
		log.Printf("Error empty request body: %s", err) // log msg with err
		// helper to insert error msg + 400 bad req status code
		WriteJSONError(w, "Email is empty", http.StatusBadRequest)
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

	// reqEmail is now successfully populated

	// check email empty
	if len(reqEmail.Email) == 0 {
		// helper to insert error msg + 400 bad req status code
		WriteJSONError(w, "Email is empty", http.StatusBadRequest)
		return // early return
	}

	// create new user using sqlc function
	newUser, err := apiCfg.db.CreateUser(req.Context(), reqEmail.Email)

	// ENSURE EMAIL IS UNIQUE (to handle error gracefully)
	pqErr, isPQError := err.(*pq.Error)

	// handle specific error first
	// check if url duplication occurred
	if isPQError && pqErr.Code == "23505" {
		// the error exists and it matches the PostgreSQL code for unique duplication
		// graceful degradation
		log.Printf("Error creating new user: %s", err)
		WriteJSONError(w, "User is already registered", http.StatusBadRequest)
		return // early return
	}

	// new user check (general)
	if err != nil {
		log.Printf("Error creating new user: %s", err) // log msg with err
		// helper to insert error msg + 400 bad req status code
		WriteJSONError(w, "Error occurred creating new user", http.StatusBadRequest)
		return // early return
	}

	// json response payload
	respUser := JsonResponseEmail{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}

	// helper to insert body response + 201 created status code
	WriteJSONResponse(w, respUser, http.StatusCreated)
}
