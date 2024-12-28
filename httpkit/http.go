package httpkit

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
