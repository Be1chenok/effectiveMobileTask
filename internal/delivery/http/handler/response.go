package handler

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

type resultResponse struct {
	Message interface{} `json:"message"`
}

func writeJsonResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	resp := resultResponse{
		Message: response,
	}
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func writeJsonErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	resp := errorResponse{
		Message: err.Error(),
	}
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
