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

	httpkit "github.com/dalton02/licor/http"
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

func validation[B any, Q any](response http.ResponseWriter, request *http.Request) (bool, *http.Request) {

	body, err := readBody(request, response)
	if err != nil {
		httpkit.GenerateErrorHttpMessage(400, "Erro ao ler o corpo da requisição", response)
		return false, request
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
		httpkit.GenerateErrorHttpMessage(400, errValidacao+errQuerys, response)
		return false, request
	}

	return true, request
}

func generic[B any, Q any](response http.ResponseWriter, request *http.Request, r *HandlerRequest[B, Q], typeRequest string) {
	validRequest := isSameRequest(typeRequest, request)
	if !validRequest {
		httpkit.AppBadRequest("Tipo de metodo não permitido, rota aceita apenas: "+typeRequest, response)
		return
	}
	contentType := request.Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "multipart/form-data") && typeRequest == "FORMDATA" {
		httpkit.AppBadRequest("Rota aceita apenas conteudo multipart/form-data não vazios", response)
		return
	}

	if strings.HasPrefix(contentType, "multipart/form-data") && typeRequest != "FORMDATA" {
		httpkit.AppBadRequest("Rota aceita apenas conteudo json", response)
		return
	}

	if strings.HasPrefix(contentType, "multipart/form-data") {
		passFormData := limitFormData(10, request)
		if !passFormData {
			httpkit.AppBadRequest("Arquivo excedeu o limite de "+strconv.Itoa(10)+" MegaBytes", response)
			return
		}
	}

	var valid bool
	if r.middleware == "public" {
		valid, request = public[B, Q](response, request)
	} else if r.middleware == "protected" {
		valid, request = protected[B, Q](response, request, r.profiles)
	}
	if !valid {
		return
	}

	params, err := extractParams(r.rota, request.URL.Path)
	if err == nil {
		ctx := context.WithValue(request.Context(), "params", params)
		//Rotas extras de middleware aqui
		validMiddleWare := runMiddlewares[B, Q](response, request.WithContext(ctx), r)
		if !validMiddleWare {
			return
		}
		r.controller(response, request.WithContext(ctx))
		return
	}
	validMiddleWare := runMiddlewares[B, Q](response, request, r)
	if !validMiddleWare {
		return
	}
	r.controller(response, request)
}

func public[B any, Q any](response http.ResponseWriter, request *http.Request) (bool, *http.Request) {

	contentType := request.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		return true, request
	}
	valid, request := validation[B, Q](response, request)
	return valid, request
}

func protected[B any, Q any](response http.ResponseWriter, request *http.Request, profiles []string) (bool, *http.Request) {

	valid := true

	contentType := request.Header.Get("Content-Type")

	auth := request.Header.Get("Authorization")
	auth = httpkit.GetBearerToken(auth)

	jwtInfo, err := httpkit.GetJwtInfo(auth)

	if err != nil {
		httpkit.AppForbidden("Token invalido/Expirado, faça login novamente", response)
		return false, request
	}

	if len(profiles) > 0 {
		pass := false
		perfil, _ := jwtInfo["perfil"].(string)
		for i := 0; i < len(profiles); i++ {
			if profiles[i] == perfil {
				pass = true
			}
		}
		if !pass {
			httpkit.AppUnauthorized("Você não está autorizado a acessar o conteudo", response)
			return false, request
		}
	}

	if !strings.HasPrefix(contentType, "multipart/form-data") {
		valid, request = validation[B, Q](response, request)
	}

	return valid, request
}

func runMiddlewares[B any, Q any](response http.ResponseWriter, request *http.Request, r *HandlerRequest[B, Q]) bool {
	//Rotas extras de middleware aqui
	validMiddleWare := false
	if len(r.overMiddleware) > 0 {
		for i := 0; i < len(r.overMiddleware); i++ {
			validMiddleWare = r.overMiddleware[i](response, request)
			if !validMiddleWare {
				return false
			}
		}
	}
	return true
}
