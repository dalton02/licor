package main

import (
	"fmt"
	"net/http"

	"github.com/dalton02/licor/httpkit"
	"github.com/dalton02/licor/licor"
)

func retrieve(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {
	jwtData := map[string]interface{}{
		"profile": "user",
		"name":    "Dalton",
	}
	token, err := httpkit.GenerateJwt(jwtData)

	if err != nil {
		return httpkit.AppBadRequest(err.Error())
	}

	jsonResponse := map[string]interface{}{
		"token": token,
	}

	return httpkit.AppSucess("Here's your token", jsonResponse)
}
func access(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {

	fmt.Println("If you have a valid token generated in the authorization header, this function will be executed")
	dataToken, _ := httpkit.GetDataToken(request)
	return httpkit.AppSucess("Your token information: ", dataToken)

}
func accessAdmin(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {

	fmt.Println("If your user is a admin user, this function will be executed")
	dataToken, _ := httpkit.GetDataToken(request)
	return httpkit.AppSucess("Your token information: ", dataToken)

}

func main() {

	//Basic Configuration of a licor setup
	licor.SetCustomInvalidTokenMessage("Token invalido/expirado")
	licor.SetCustomNotAuthorizedMessage("Perfil não tem acesso ao conteudo atual")
	licor.SetBearerTokenAuthorizationHeader() //Authorization via bearer token in the authorization header http

	licor.Public[any, any]("/retrieve").Get(retrieve)
	licor.Protected[any, any]("/access").Get(access)
	licor.Protected[any, any]("/access-admin", "admin").Get(accessAdmin)

	licor.Init("3003")
}

func Middle1(response http.ResponseWriter, request *http.Request) (httpkit.HttpMessage, bool) {
	var message httpkit.HttpMessage
	fmt.Println("Everything good around here, will proceed to the next middleware")
	return message, true
}

func Middle2(response http.ResponseWriter, request *http.Request) (httpkit.HttpMessage, bool) {
	var message httpkit.HttpMessage
	message = httpkit.AppBadRequest("Something went wrong, can't proceed operation")
	return message, false
}
