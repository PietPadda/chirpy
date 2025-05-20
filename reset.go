// reset.go
package main

import (
	"log"
	"net/http"
	"strconv"
)

// AdminMetricsResets handler that resets the fileserver hits!
// apply cfg receiver to access fileserverHits
func (cfg *apiConfig) handlerAdminMetricsReset(w http.ResponseWriter, req *http.Request) {
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

// AdminUsersReset handler that resets the fileserver hits!
// apply cfg receiver to access fileserverHits
func (cfg *apiConfig) handlerAdminUsersReset(w http.ResponseWriter, req *http.Request) {
	// check if env PLATFORM is dev
	if cfg.platform != "dev" {
		// handle gracefully
		log.Printf("Unauthorised to reset users") // msg to server admin
		// send msg to client code 403
		WriteJSONError(w, "Unauthorised to reset users", http.StatusForbidden)
		return // stop processing req
	}

	// HTTP method check
	if req.Method != "POST" {
		// helper to insert error msg + 405 invalid method status code
		WriteJSONError(w, "Users reset must be POSTed", http.StatusMethodNotAllowed)
		return // early return
	}

	// apiConfig check
	if cfg == nil {
		// handle gracefully
		log.Printf("Internal server error: apiCfg is nil") // msg to server admin
		// send msg to client code 500
		WriteJSONError(w, "Internal server configuration error", http.StatusInternalServerError)
		return // stop processing req
	}

	// remove users from database
	err := cfg.db.ResetUsers(req.Context())

	// remove users check
	if err != nil {
		// handle gracefully
		log.Printf("Internal server error: couldn't reset users %v", err) // msg to server admin
		// send msg to client code 500
		WriteJSONError(w, "Error resetting users", http.StatusInternalServerError)
		return // stop processing req
	}

	// successfully reset db users
	// handle gracefully
	log.Printf("DB users have been reset") // msg to server admin
	// send msg to client code 200
	w.WriteHeader(http.StatusOK)
}
