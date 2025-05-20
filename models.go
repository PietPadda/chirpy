// models.go
package main

import (
	"time"

	"github.com/google/uuid"
)

// structs for API json request and responses

// API JSON Request from Client
type JsonRequest struct {
	Body string `json:"body"`
}

type JsonRequestEmail struct {
	Email string `json:"email"`
}

// API JSON Response to Client
type JsonResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type JsonResponseError struct {
	Error string `json:"error"`
}

type JsonResponseEmail struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
