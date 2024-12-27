package httpkit

import (
	"encoding/json"
	"net/http"
)

type HttpMessage struct {
	StatusCode int    `json:"statusCode"`
	Data       any    `json:"data"`
	Message    string `json:"message"`
}

func GenerateHttpMessage[T any](statusCode int, data T, message string) HttpMessage {
	var dataResponse = HttpMessage{
		StatusCode: statusCode,
		Data:       data,
		Message:    message,
	}
	return dataResponse
}

func GenerateErrorHttpMessage(statusCode int, message string) HttpMessage {
	var dataResponse = HttpMessage{
		StatusCode: statusCode,
		Data:       nil,
		Message:    message,
	}
	return dataResponse
}

func SendHttpMessage(message HttpMessage, response http.ResponseWriter) {
	json.Marshal(message)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(message.StatusCode)
	if err := json.NewEncoder(response).Encode(message); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}
