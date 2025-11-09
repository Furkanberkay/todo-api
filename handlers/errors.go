package handlers

import (
	"log"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ValidationError(s string) ErrorResponse {
	return ErrorResponse{
		Error: s,
	}
}
func InternalError(err error) ErrorResponse {
	if err != nil {
		log.Printf("internal error: %v", err)
	}
	return ErrorResponse{
		Error: "internal server error",
	}
}
