package licor

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/dalton02/licor/httpkit"
	"github.com/dalton02/licor/validator"
)

func readBody(request *http.Request, response http.ResponseWriter) ([]byte, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(request.Body, &buf)
	body, err := ioutil.ReadAll(tee)
	if err != nil {
		return body, err
	}
	if len(body) == 0 {
		body = []byte("{}")
	}
	request.Body = ioutil.NopCloser(bytes.NewBuffer(buf.Bytes()))
	return body, nil
}

func validation[B any, Q any](response http.ResponseWriter, request *http.Request) (bool, *http.Request, httpkit.HttpMessage) {

	var message httpkit.HttpMessage
	body, err := readBody(request, response)
	if err != nil {
		message = httpkit.GenerateErrorHttpMessage(400, "Erro ao ler o corpo da requisição")
		return false, request, message
	}

	var dataR B
	jsonString := body
	json.Unmarshal(body, &dataR)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(jsonString), &jsonData)
	keysByLevel := make(map[int][]string)
	ctx := context.WithValue(request.Context(), "original_body", body)
	request = request.WithContext(ctx)
	errValidacao, hasError := validator.CheckPropretys[B](dataR, validator.ExtractKeysByLevel(jsonData, 1, keysByLevel))
	params, has, maping := extractQueryParams[Q](request)
	errQuerys := ""
	hasErrorQ := false

	if has {
		errQuerys, hasErrorQ = validator.CheckPropretys[Q](params, validator.QueryMap(maping))
	}
	if hasError || hasErrorQ {
		return false, request, httpkit.GenerateErrorHttpMessage(400, errValidacao+errQuerys)
	}

	return true, request, message
}

func generic[B any, Q any](response http.ResponseWriter, request *http.Request, r *HandlerRequest[B, Q], typeRequest string) httpkit.HttpMessage {
	validRequest := isSameRequest(typeRequest, request)
	if !validRequest {
		return httpkit.AppBadRequest("Tipo de metodo não permitido, rota aceita apenas: " + typeRequest)

	}
	contentType := request.Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "multipart/form-data") && typeRequest == "FORMDATA" {
		return httpkit.AppBadRequest("Rota aceita apenas conteudo multipart/form-data não vazios")
	}

	if strings.HasPrefix(contentType, "multipart/form-data") && typeRequest != "FORMDATA" {

		return httpkit.AppBadRequest("Rota aceita apenas conteudo json")
	}

	if strings.HasPrefix(contentType, "multipart/form-data") {
		passFormData := limitFormData(maxSizeFormData, request)
		if !passFormData {
			return httpkit.AppBadRequest("Arquivo excedeu o limite de " + strconv.Itoa(maxSizeFormData) + " MegaBytes")
		}
	}

	var valid bool
	var message httpkit.HttpMessage
	if r.security == "public" {
		valid, request, message = public[B, Q](response, request)
	} else if r.security == "protected" {
		valid, request, message = protected[B, Q](response, request, r.extraRoute)
	}
	if !valid {
		return message
	}
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		valid, request, message = validation[B, Q](response, request)
		if !valid {
			return message
		}
	}

	params, err := extractParams(r.endpoint, request.URL.Path)
	if err == nil {
		ctx := context.WithValue(request.Context(), "params", params)
		message, validMiddleWare := runMiddlewares[B, Q](response, request.WithContext(ctx), r)
		if !validMiddleWare {
			return message
		}
		message = r.controller(response, request.WithContext(ctx))
		return message
	}
	message, validMiddleWare := runMiddlewares[B, Q](response, request, r)
	if !validMiddleWare {
		return message
	}
	message = r.controller(response, request)
	return message
}

func public[B any, Q any](response http.ResponseWriter, request *http.Request) (bool, *http.Request, httpkit.HttpMessage) {

	var message httpkit.HttpMessage
	return true, request, message
}

// When making custom protected routes, one must return a boolean and a httpkit message
func protected[B any, Q any](response http.ResponseWriter, request *http.Request, extras []any) (bool, *http.Request, httpkit.HttpMessage) {

	var valid bool
	var message httpkit.HttpMessage

	switch currentProtection {
	case "bearerTokenAuthorizationHeader":
		valid, request, message = authorizationHeader(response, request, extras)
	case "custom":
		valid, request, message = customProtection(response, request, extras)
	default:
		valid, request, message = authorizationHeader(response, request, extras)

	}

	return valid, request, message
}

func runMiddlewares[B any, Q any](response http.ResponseWriter, request *http.Request, r *HandlerRequest[B, Q]) (httpkit.HttpMessage, bool) {

	var message httpkit.HttpMessage
	validMiddleWare := false
	if len(r.middlewares) > 0 {
		for i := 0; i < len(r.middlewares); i++ {
			message, validMiddleWare = r.middlewares[i](response, request)
			if !validMiddleWare {
				return message, false
			}
		}
	}
	return message, true
}
