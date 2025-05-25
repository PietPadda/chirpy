package main

import (
	"log"
	"net/http"
	"time"

	"github.com/PietPadda/chirpy/internal/auth"
)

// Refresh handler that reissues access token if refresh is valid
func (apiCfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
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
		WriteJSONError(w, "Login must be POSTed", http.StatusMethodNotAllowed)
		return // early return
	}

	// get request bearer token
	bearerToken, err := auth.GetBearerToken(req.Header) // pass request header

	// bearer token check
	if err != nil {
		log.Printf("Error getting bearer token: %s", err) // log msg with err
		// helper to insert error msg + 401 unauthorised req status
		WriteJSONError(w, "Invalid token", http.StatusUnauthorized) // general error, obscure to client
		return                                                      // early return
	}

	// get user from refresh token
	loginUser, err := apiCfg.db.GetUserFromRefreshToken(req.Context(), bearerToken)

	// get user check
	if err != nil {
		log.Printf("Error getting user from refresh token: %s", err) // log msg with err
		// helper to insert error msg + 401 unauthorised req status
		WriteJSONError(w, "Invalid token", http.StatusUnauthorized) // general error, obscure to client
		return                                                      // early return
	}

	// refresh token expiration check
	if time.Now().UTC().After(loginUser.ExpiresAt) {
		log.Printf("Error refresh token expired") // log msg
		// helper to insert error msg + 401 unauthorised req status code
		WriteJSONError(w, "Invalid token", http.StatusUnauthorized) // general error, obscure to client
		return                                                      // early return
	}

	// refresh token revoked check
	if loginUser.RevokedAt.Valid { // recall that nullable structs have time & valid fields
		log.Printf("Error refresh token revoked") // log msg
		// helper to insert error msg + 401 unauthorised req status code
		WriteJSONError(w, "Invalid token", http.StatusUnauthorized) // general error, obscure to client
		return                                                      // early return
	}

	// set default expiration time for JWT
	expiresDuration := 3600 * time.Second // 1 hour

	// make JWT token
	tokenString, err := auth.MakeJWT(loginUser.UserID, apiCfg.serverKey, expiresDuration)

	// check make jwt
	if err != nil {
		log.Printf("Error making JWT token: %s", err) // log msg with err
		// helper to insert error msg + 500 internal server error status code
		WriteJSONError(w, "Internal server token generation error", http.StatusInternalServerError)
		return // early return
	}

	// json response payload
	respLogin := JsonRefreshResponse{
		Token: tokenString,
	}

	// helper to insert body response + 200 ok  status code
	WriteJSONResponse(w, respLogin, http.StatusOK)
}

// Revoke handler that revokes the refresh token, preventing new access tokens to be issued
func (apiCfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
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
		WriteJSONError(w, "Login must be POSTed", http.StatusMethodNotAllowed)
		return // early return
	}

	// get request bearer token
	bearerToken, err := auth.GetBearerToken(req.Header) // pass request header

	// bearer token check
	if err != nil {
		log.Printf("Error getting bearer token: %s", err) // log msg with err
		// helper to insert error msg + 401 unauthorised req status
		WriteJSONError(w, "Invalid token", http.StatusUnauthorized) // general error, obscure to client
		return                                                      // early return
	}

	// revoke the refresh token
	err = apiCfg.db.RevokeRefreshToken(req.Context(), bearerToken)

	// revoke refresh token check
	if err != nil {
		log.Printf("Error revoking refresh token: %s", err) // log msg with err
		// helper to insert error msg + 500 internal server error
		WriteJSONError(w, "Internal server error", http.StatusInternalServerError)
		return // stop processing req
	}

	// no response struct, just write msg to client
	// write the header with 204 No content
	log.Printf("Revoking refresh token: %s", bearerToken) // server msg
	w.WriteHeader(http.StatusNoContent)                   // resp to client
}
