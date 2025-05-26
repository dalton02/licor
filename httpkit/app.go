package httpkit

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
	AppConflict(handlerFunc func(message string) HttpMessage)
	AppBadRequest(handlerFunc func(message string) HttpMessage)
	AppUnauthorized(handlerFunc func(message string) HttpMessage)
	AppForbidden(handlerFunc func(message string) HttpMessage)
	AppNotFound(handlerFunc func(message string) HttpMessage)
	AppInternal(handlerFunc func(message string) HttpMessage)
	AppNotImplemented(handlerFunc func(message string) HttpMessage)
	AppSuccess(handlerFunc func(message string, data any) HttpMessage)
	AppSuccessCreate(handlerFunc func(message string, data any) HttpMessage)
}

func AppSucess(message string, data any) HttpMessage {
	return GenerateHttpMessage(200, data, message)
}

func AppSucessCreate(message string, data any) HttpMessage {
	return GenerateHttpMessage(201, data, message)
}
func AppConflict(message string) HttpMessage {
	return GenerateErrorHttpMessage(409, message)
}
func AppBadRequest(message string) HttpMessage {
	return GenerateErrorHttpMessage(400, message)
}

func AppUnauthorized(message string) HttpMessage {
	return GenerateErrorHttpMessage(401, message)
}

func AppForbidden(message string) HttpMessage {
	return GenerateErrorHttpMessage(403, message)
}

func AppNotFound(message string) HttpMessage {
	return GenerateErrorHttpMessage(404, message)
}

func AppInternal(message string) HttpMessage {
	return GenerateErrorHttpMessage(500, message)
}

func AppNotImplemented(message string) HttpMessage {
	return GenerateErrorHttpMessage(501, message)
}
