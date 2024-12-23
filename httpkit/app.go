package httpkit

import (
	"net/http"
)

type HttpMessageDto struct {
	Message string
	Status  int
}

func NewHttpMessage(message string, status int) HttpMessageDto {
	return HttpMessageDto{
		Message: message,
		Status:  status,
	}
}

type App interface {
	AppConflict(handlerFunc func(message string, response http.ResponseWriter))
	AppBadRequest(handlerFunc func(message string, response http.ResponseWriter))
	AppUnauthorized(handlerFunc func(message string, response http.ResponseWriter))
	AppForbidden(handlerFunc func(message string, response http.ResponseWriter))
	AppNotFound(handlerFunc func(message string, response http.ResponseWriter))
	AppInternal(handlerFunc func(message string, response http.ResponseWriter))
	AppNotImplemented(handlerFunc func(message string, response http.ResponseWriter))
	AppSucess(handlerFunc func(message string, data any, response http.ResponseWriter))
	AppSucessCreate(handlerFunc func(message string, data any, response http.ResponseWriter))
}

func AppSucess(message string, data any, response http.ResponseWriter) {
	GenerateHttpMessage(200, data, message, response)
}

func AppSucessCreate(message string, data any, response http.ResponseWriter) {
	GenerateHttpMessage(201, data, message, response)
}
func AppConflict(message string, response http.ResponseWriter) {
	GenerateErrorHttpMessage(409, message, response)
}
func AppBadRequest(message string, response http.ResponseWriter) {
	GenerateErrorHttpMessage(400, message, response)
}

func AppUnauthorized(message string, response http.ResponseWriter) {
	GenerateErrorHttpMessage(401, message, response)
}

func AppForbidden(message string, response http.ResponseWriter) {
	GenerateErrorHttpMessage(403, message, response)
}

func AppNotFound(message string, response http.ResponseWriter) {
	GenerateErrorHttpMessage(404, message, response)
}

func AppInternal(message string, response http.ResponseWriter) {
	GenerateErrorHttpMessage(500, message, response)
}

func AppNotImplemented(message string, response http.ResponseWriter) {
	GenerateErrorHttpMessage(501, message, response)
}
