package httpkit

import (
	"encoding/json"
	"net/http"
)

type HttpMessage[T any] struct {
	StatusCode int    `json:"statusCode"`
	Data       T      `json:"data"`
	Message    string `json:"message"`
}

func GenerateHtmlResponse(statusCode int, html string, response http.ResponseWriter) {
	response.Header().Set("Content-Type", "text/html")
	response.WriteHeader(statusCode)
	_, err := response.Write([]byte(html))
	if err != nil {
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
	}
}

func GenerateHttpMessage[T any](statusCode int, data T, message string, response http.ResponseWriter) {
	var dataResponse = HttpMessage[T]{
		StatusCode: statusCode,
		Data:       data,
		Message:    message,
	}
	json.Marshal(dataResponse)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	if err := json.NewEncoder(response).Encode(dataResponse); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func GenerateErrorHttpMessage(statusCode int, message string, response http.ResponseWriter) {
	var dataResponse = HttpMessage[any]{
		StatusCode: statusCode,
		Data:       nil,
		Message:    message,
	}
	json.Marshal(dataResponse)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	if err := json.NewEncoder(response).Encode(dataResponse); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}
