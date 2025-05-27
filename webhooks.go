// webhooks.go
package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// PolkaWebhook handler that sets user to premiums
func (apiCfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, req *http.Request) {
	// HTTP method check
	if req.Method != "POST" {
		log.Printf("Invalid HTTP method used: %s", req.Method) // log msg
		w.WriteHeader(http.StatusMethodNotAllowed)             // status code 405 to client
		return                                                 // early return
	}

	// json request from client
	var reqChirpyRed JsonPolkaWebhookRequest

	// create json req body decoder
	decoder := json.NewDecoder(req.Body)

	// close on exit to prevent mem leak
	defer req.Body.Close()

	// decode the req body
	err := decoder.Decode(&reqChirpyRed)

	// request body missing edge case check (before general error check)
	if err == io.EOF { // end of file
		log.Printf("Error empty request body: %s", err) // log msg with err
		w.WriteHeader(http.StatusBadRequest)            // status code 400 to client
		return                                          // early return
	}

	// decode check
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value

		log.Printf("Error decoding parameters: %s", err) // log msg with err
		w.WriteHeader(http.StatusBadRequest)             // status code 400 to client
		return                                           // early return
	}

	// wrong event check
	if reqChirpyRed.Event != "user.upgraded" {
		// write to server and client that chirp deleted
		log.Printf("Incorrect event request: Event = %s", reqChirpyRed.Event) // log msg with err
		w.WriteHeader(http.StatusNoContent)                                   // status code 204 to client
		return                                                                // early return
	}

	// no user id provided check
	if reqChirpyRed.Data.UserID.String() == "" { // stringified UUID
		// write to server and client that chirp deleted
		log.Printf("Data does not include user id") // log msg with err
		w.WriteHeader(http.StatusBadRequest)        // status code 400 to client
		return                                      // early return
	}

	// reqChirpyRed is now successfully populated

	// proceed to update the user's chirpy status to red "premium"
	chirpyRed, err := apiCfg.db.SetIsChirpyRedTrue(req.Context(), reqChirpyRed.Data.UserID)

	// user doesn't exist check (no updates!)
	if err == sql.ErrNoRows { // no rows updated
		log.Printf("User could not be found as no record updated: ID = %s",
			reqChirpyRed.Data.UserID) // msg to server admin
		w.WriteHeader(http.StatusNotFound) // status code 404 to client
		return                             // stop processing req
	}

	// set chirpy red true check
	if err != nil {
		// handle gracefully
		log.Printf("Error set chirpy red true: %s", err) // msg to server admin
		w.WriteHeader(http.StatusInternalServerError)    // status code 500 to client
		return                                           // stop processing req
	}

	// write to server and client that user upgrade to chirpy red
	log.Printf("User has been upgraded to chirpy red: ID = %s", chirpyRed.ID) // log msg with err
	w.WriteHeader(http.StatusNoContent)                                       // status code 204 to client
}
