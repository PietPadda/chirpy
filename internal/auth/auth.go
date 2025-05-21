// auth.go
package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// func to hash a user's password
func HashPassword(password string) (string, error) {
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
