// models.go
package main

import (
	"time"

	"github.com/google/uuid"
)

// structs for API json request and responses

// REQUESTS
type JsonRequest struct {
	Body string `json:"body"`
}

// CreateUser request
type JsonUserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

// CreateChirp request
type JsonChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

// UserLogin request
type JsonLoginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

// RESPONSES
// API JSON Response to Client
type JsonResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

// Client error response
type JsonResponseError struct {
	Error string `json:"error"`
}

// Client user created response
type JsonUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// Client login successful response
type JsonLoginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// Client chirp response
type JsonChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}
