// auth_test.go

package auth

import (
	"testing" // importing testing package for unit tests
)

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
