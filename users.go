// users.go
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	// our internal package
	// postgresql db access
	"github.com/PietPadda/chirpy/internal/auth"
	"github.com/PietPadda/chirpy/internal/database"
	"github.com/lib/pq" // postgresql driver
)

// CreateUser handler that creates a new user
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
	var reqEmail JsonUserRequest

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

	// hash the password
	hash, err := auth.HashPassword(reqEmail.Password)

	// hash check
	if err != nil {
		log.Printf("Error hashing password: %s", err) // log msg with err
		// helper to insert error msg + 500 internal error status code
		WriteJSONError(w, "Error hashing password", http.StatusInternalServerError)
		return // early return
	}

	// create new user using sqlc function
	newUser, err := apiCfg.db.CreateUser(req.Context(), database.CreateUserParams{
		HashedPassword: hash,           // hashed password
		Email:          reqEmail.Email, // directly  user input
	})

	// new user check (general)
	if err != nil {
		log.Printf("Error creating new user: %s", err) // log msg with err
		// helper to insert error msg + 400 bad req status code
		WriteJSONError(w, "Error occurred creating new user", http.StatusBadRequest)
		return // early return
	}

	// json response payload
	respUser := JsonUserResponse{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}

	// helper to insert body response + 201 created status code
	WriteJSONResponse(w, respUser, http.StatusCreated)
}

// LoginUser handler that allows user login by password and email
func (apiCfg *apiConfig) handlerUserLogin(w http.ResponseWriter, req *http.Request) {
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

	// json request from client
	var reqLogin JsonLoginRequest

	// create json req body decoder
	decoder := json.NewDecoder(req.Body)

	// close on exit to prevent mem leak
	defer req.Body.Close()

	// decode the req body
	err := decoder.Decode(&reqLogin)

	// request body missing edge case check (before general error check)
	if err == io.EOF { // end of file
		log.Printf("Error empty request body: %s", err) // log msg with err
		// helper to insert error msg + 400 bad req status code
		WriteJSONError(w, "Request body is empty", http.StatusBadRequest)
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

	// set default expiration time for JWT
	expiresDuration := 3600 * time.Second // 1 hour

	// reqLogin is now successfully populated

	// get the user by email
	loginUser, err := apiCfg.db.GetUserByEmail(req.Context(), reqLogin.Email)

	// get user check
	if err != nil {
		log.Printf("Error getting user by email: %s", err) // log msg with err
		// helper to insert error msg + 401 unauthorised req status code
		WriteJSONError(w, "Incorrect email or password", http.StatusUnauthorized)
		return // early return
	}

	// get user hashed password
	hash := loginUser.HashedPassword

	// compare hashed password with input password
	err = auth.CheckPasswordHash(hash, reqLogin.Password)

	// hash check
	if err != nil {
		log.Printf("Error invalid password: %s", err) // log msg with err
		// helper to insert error msg + 401 unauthorised error status code
		WriteJSONError(w, "Incorrect email or password", http.StatusUnauthorized)
		return // early return
	}

	// make JWT token
	tokenString, err := auth.MakeJWT(loginUser.ID, apiCfg.serverKey, expiresDuration)

	// check make jwt
	if err != nil {
		log.Printf("Error making JWT token: %s", err) // log msg with err
		// helper to insert error msg + 500 internal server error status code
		WriteJSONError(w, "Internal server token generation error", http.StatusInternalServerError)
		return // early return
	}

	// generate a refresh token
	tokenRefreshString, err := auth.MakeRefreshToken()

	// check generate refresh token
	if err != nil {
		log.Printf("Error making refresh token: %s", err) // log msg with err
		// helper to insert error msg + 500 internal server error status code
		WriteJSONError(w, "Internal server refresh token generation error", http.StatusInternalServerError)
		return // early return
	}

	// set default expiration time for refresh token
	expiresRefreshDuration := 60 * 24 * time.Hour                           // 60 days
	expiresRefreshTimestamp := time.Now().UTC().Add(expiresRefreshDuration) // conv to timestamp for PostgreSQL

	// add refresh token to the db
	_, err = apiCfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     tokenRefreshString,      // add refresh token string
		UserID:    loginUser.ID,            // set to correct user
		ExpiresAt: expiresRefreshTimestamp, // conv to postgresql time
	})

	// create refresh token check
	if err != nil {
		log.Printf("Error adding refresh token to database: %s", err) // log msg with err
		// helper to insert error msg + 500 internal server error status code
		WriteJSONError(w, "Internal server refresh token adding to database error", http.StatusInternalServerError)
		return // early return
	}

	// json response payload
	respLogin := JsonLoginResponse{
		ID:           loginUser.ID,
		CreatedAt:    loginUser.CreatedAt,
		UpdatedAt:    loginUser.UpdatedAt,
		Email:        loginUser.Email,
		Token:        tokenString,
		RefreshToken: tokenRefreshString,
	}

	// helper to insert body response + 200 ok  status code
	WriteJSONResponse(w, respLogin, http.StatusOK)
}
