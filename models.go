// models.go
package main

// structs for json request and responses

type JsonRequest struct {
	Body string `json:"body"`
}

type JsonResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type JsonResponseError struct {
	Error string `json:"error"`
}
