package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dalton02/licor/httpkit"
	"github.com/dalton02/licor/licor"
	"github.com/rs/cors"
)

func retrieve(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {
	jwtData := map[string]interface{}{
		"profile": "user",
		"name":    "Dalton",
	}
	token, err := httpkit.GenerateJwt(jwtData, 30)

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
	dataToken, _ := httpkit.GetDataToken(request)
	return httpkit.AppSucess("Your token information: ", dataToken)
}

func custom(response http.ResponseWriter, request *http.Request, extras ...any) (bool, *http.Request, httpkit.HttpMessage) {
	var message httpkit.HttpMessage
	ctx := context.WithValue(request.Context(), "someData", "The main function will recieve this information")
	return true, request.WithContext(ctx), message
}

type q struct {
	Queryzinha string `query:"queryzinha" validator:"dateString"`
}

func main() {

	licor.SetMaxSizeFormData(10)
	licor.SetCustomInvalidTokenMessage("Token invalido/expirado")
	licor.SetCustomNotAuthorizedMessage("Perfil não tem acesso ao conteudo atual")
	licor.SetCustomProtection(custom)

	licor.SetCors(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	licor.Public[any, any]("/retrieve").Get(retrieve)
	licor.Protected[any, any]("/access").Get(access)
	licor.Protected[any, q]("/access-admin", "admin").Get(accessAdmin)

	licor.Init("3003")
}

func middle1(response http.ResponseWriter, request *http.Request) (httpkit.HttpMessage, bool) {
	var message httpkit.HttpMessage
	fmt.Println("Everything good around here, will proceed to the next middleware")
	return message, true
}

func middle2(response http.ResponseWriter, request *http.Request) (httpkit.HttpMessage, bool) {
	var message httpkit.HttpMessage
	message = httpkit.AppBadRequest("Something went wrong, can't proceed operation")
	return message, false
}
