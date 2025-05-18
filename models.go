// models.go
package main

// structs for json request and responses

type JsonRequest struct {
	Body string `json:"body"`
}

type JsonResponse struct {
	Valid bool `json:"valid"`
}

type JsonResponseError struct {
	Error string `json:"error"`
}
