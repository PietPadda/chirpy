// auth_test.go

package auth

import (
	"strings"
	"testing" // importing testing package for unit tests
	"time"

	"github.com/google/uuid"
)

// PASSWORD HASHING
// test pw hashing
func TestHashPassword(t *testing.T) {
	// test case
	password := "password123"

	// hash pw
	hash, err := HashPassword(password)

	// hash fail check
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err) // fatal, don't continue
	}

	// empty hash check
	if hash == "" {
		t.Errorf("HashPassword returned empty hash") // simple test, continue w errorf
	}

	// unhashed check
	if hash == password {
		t.Errorf("HashPassword returned unhashed password")
	}

	// create new has with same pw
	hash2, _ := HashPassword(password)
	// no need for err, already tested

	// duplicate hash check
	if hash == hash2 {
		t.Errorf("HashPassword produced duplicated hashes")
	}
}

// test pw hash compare
func TestCheckPasswordHash(t *testing.T) {
	// test case
	password := "password123"

	// create hash
	hash, _ := HashPassword(password)
	// no need for err, already tested

	// compare hash & pw
	err := CheckPasswordHash(hash, password)

	// compare failed check
	if err != nil {
		t.Fatalf("CheckPasswordHash failed: %v", err) // fatal, don't continue
	}

	// compare works check
	passwordTypo := "password1234icanguessyourpassword!"
	err = CheckPasswordHash(hash, passwordTypo)
	if err == nil {
		t.Errorf("CheckPasswordHash failed to detect incorrect password")
	}
}

// JSON WEB TOKEN GENERATION
// test jwt generation
func TestMakeJWT(t *testing.T) {
	// test case
	userID := "123e4567-e89b-12d3-a456-426614174000" // fixed stringified uuid
	userUUID, _ := uuid.Parse(userID)                // no need to err check uuid
	tokenSecret := "AllYourBase"
	expiresIn := time.Hour * 24

	// gen jwt token
	tokenString, err := MakeJWT(userUUID, tokenSecret, expiresIn)

	// gen token check
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err) // fatal, don't continue
	}

	// check JWT fundamental structure
	dots := strings.Count(tokenString, ".")

	if dots != 2 { // all JWTs take form "header.payload.signature" ie 2 dots!
		t.Errorf("MakeJWT returned invalid JWT format: %s", tokenString)
	}
}

// test jwt validation
func TestValidateJWT(t *testing.T) {
	// test case
	userID := "123e4567-e89b-12d3-a456-426614174000" // fixed stringified uuid
	userUUID, _ := uuid.Parse(userID)                // no need to err check uuid
	tokenSecret := "AllYourBase"
	expiresIn := time.Hour * 24

	// gen jwt token
	tokenString, _ := MakeJWT(userUUID, tokenSecret, expiresIn) // err checked in other test

	// validated token
	userIDValid, err := ValidateJWT(tokenString, tokenSecret)

	// validate token check
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err) // fatal, don't continue
	}

	// check if userIDs match
	if userIDValid != userUUID {
		t.Errorf("ValidateJWT returned invalid user ID: %s", userIDValid)
	}
}

// test jwt validation with incorrect secret key
func TestValidateJWTWrongKey(t *testing.T) {
	// test case
	userID := "123e4567-e89b-12d3-a456-426614174000" // fixed stringified uuid
	userUUID, _ := uuid.Parse(userID)                // no need to err check uuid
	tokenSecret := "AllYourBase"
	tokenSecretWrong := "HaHaIHackedYou" // wrong key :)
	expiresIn := time.Hour * 24

	// gen jwt token (with wrong key added!)
	tokenString, _ := MakeJWT(userUUID, tokenSecretWrong, expiresIn) // err checked in other test

	// validated token (or attempt to...)
	_, err := ValidateJWT(tokenString, tokenSecret)

	// validate token check
	if err == nil {
		t.Fatalf("ValidateJWT failed to invalidate wrong key: %v", err) // fatal, don't continue
	}
}

// test jwt validation with expired token
func TestValidateJWTExpiredToken(t *testing.T) {
	// test case
	userID := "123e4567-e89b-12d3-a456-426614174000" // fixed stringified uuid
	userUUID, _ := uuid.Parse(userID)                // no need to err check uuid
	tokenSecret := "AllYourBase"
	expiresIn := time.Hour * -48 // expired 2 days ago :)

	// gen jwt token (with expired timer!)
	tokenString, _ := MakeJWT(userUUID, tokenSecret, expiresIn) // err checked in other test

	// validated token (or attempt to...)
	_, err := ValidateJWT(tokenString, tokenSecret)

	// validate token check
	if err == nil {
		t.Fatalf("ValidateJWT failed to invalidate expired key: %v", err) // fatal, don't continue
	}
}
