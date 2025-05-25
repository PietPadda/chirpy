// auth.go
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// PASSWORD HASHING
// func to hash a user's password
func HashPassword(password string) (string, error) {
	// handle empty password
	if password == "" {
		return "", errors.New("empty password") // early return
	}

	// use bcrypt's pw gen -- accepts a byte (max 72) and cost (1-31)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// hash pw string check
	if err != nil {
		return "", err // early return
	}
	// low-level funcs shouldn't log or send client responses

	// successfully hash the pw
	return string(hashedPassword), nil
}

// compare hashed pw with user login string
func CheckPasswordHash(hash, password string) error {
	// use bcrypt's compare func
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	// we return the resulting error (and let the func caller handle it)
	return err
}

// GENERATE AND VERIFY JWT TOKENS
// generate jwt token on server to send to user
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// the signing key
	signingKey := []byte(tokenSecret)

	// HMAC HS256 signing method
	signingMethod := jwt.SigningMethodHS256

	// create registered claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",                                      // issuer = our application
		IssuedAt:  jwt.NewNumericDate(time.Now()),                // current time
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)), // current time + expiration time
		Subject:   userID.String(),                               // stringified version of user id
	}

	// create JWT token using signing method and registered claim
	token := jwt.NewWithClaims(signingMethod, claims)

	// sign the token with secret key
	tokenString, err := token.SignedString(signingKey)

	// token sign check
	if err != nil {
		return "", err // failed to sign token with secret ke
	}
	// low level func, let high levels handle the error

	// successfully created signed jwt token
	return tokenString, nil
}

// validate jwt token returned from user
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// create empty claims struct to be populated
	claims := &jwt.RegisteredClaims{} // ptr because ParseWithClaims requires it

	// parse the user's jwt token claims, validates it, then returns the parsed token
	tokenParse, err := jwt.ParseWithClaims(tokenString, claims, func(tokenParse *jwt.Token) (interface{}, error) {
		// return secret key a byte slice
		return []byte(tokenSecret), nil
	}) // anon func takes token from ParseWithClaims and returns the secret key as byte slice

	// check token claims parse
	if err != nil {
		return uuid.Nil, err // nil id
	}

	// Use a type assertion to get the claims as *jwt.RegisteredClaims
	token, ok := tokenParse.Claims.(*jwt.RegisteredClaims)
	// tokenParse.Claims is of type jwt.Claims (interface)

	// type assertion check
	if !ok {
		return uuid.Nil, errors.New("invalid token claims") // nil id
	} // use custom err, not just "err" (from previous check...)

	// token expiration check
	if time.Now().After(token.ExpiresAt.Time) {
		return uuid.Nil, errors.New("token has expired") // nil id
	} // use custom err, not just "err" (from previous check...)

	// convert userid (subject) to uuid
	userID, err := uuid.Parse(token.Subject)

	// uuid parse check
	if err != nil {
		return uuid.Nil, err // nil id
	}

	// return userid (subject) as uuid from the populate claims
	// validation confirms which user this JWT belongs to
	return userID, nil
}

// BEARER TOKEN
// checks the client's token from the header!
func GetBearerToken(headers http.Header) (string, error) {
	// get that header type
	authHeader := headers.Get("Authorization") // .Get() is a method that works on http.Header

	// if not "Get"ed, .Get() will return an empty string
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	// split the authHeader to get all terms
	headerFields := strings.Fields(authHeader) // get each term regardless of whitespace
	// NOTE: string.Split(...," ") doesn't handle multiple spaces
	// NOTE: string.TrimSpace(...) only handles whitespace at start adn end

	// check that header contains 2 words!
	if len(headerFields) != 2 {
		return "", errors.New("invalid authorization header")
	} // else we get indexerror

	// check auth type
	if headerFields[0] != "Bearer" {
		return "", errors.New("missing authorization header")
	}

	// check token string length not empty
	if headerFields[1] == "" {
		return "", errors.New("missing token string")
	}

	// get client's JWT token string auth type
	tokenString := headerFields[1]

	// return string and success
	return tokenString, nil
}

// REFRESH TOKEN
// makes a refresh token for user
func MakeRefreshToken() (string, error) {
	// make zero'd slice wtih 32 bytes (which is 256 bits)
	key := make([]byte, 32)
	// each byte is 8 bits, meaning 0-255 or 2^8

	// fill the slice with random raw bytes 0-255
	_, err := rand.Read(key) // bytes space efficient, hard to transmit

	// random check
	if err != nil {
		return "", err // early return
	}

	// encode the key to hex string
	encodedKey := hex.EncodeToString(key) // hex easier to transmit, takes bit more space

	// return successfully made refresh token string
	return encodedKey, nil
}
