package licor

import (
	"net/http"

	"github.com/dalton02/licor/httpkit"
)

type Protections interface {
	AuthorizationHeader(func(response http.ResponseWriter, request *http.Request, extras ...any) (bool, *http.Request, httpkit.HttpMessage))
}

func AuthorizationHeader(response http.ResponseWriter, request *http.Request, extras []any) (bool, *http.Request, httpkit.HttpMessage) {

	var message httpkit.HttpMessage
	auth := request.Header.Get("Authorization")
	auth = httpkit.GetBearerToken(auth)

	jwtInfo, err := httpkit.GetJwtInfo(auth)

	if err != nil {
		message = httpkit.AppForbidden(invalidToken)
		return false, request, message
	}

	if len(extras) > 0 {
		pass := false
		perfil, _ := jwtInfo["profile"].(string)
		for i := 0; i < len(extras); i++ {
			if extras[i] == perfil {
				pass = true
			}
		}
		if !pass {
			message = httpkit.AppUnauthorized(profile401)
			return false, request, message
		}
	}

	return true, request, message
}

var invalidToken string = "Token not valid or expired"
var profile401 string = "Your profile has no access to this content"

var CustomProtection func(response http.ResponseWriter, request *http.Request, extras ...any) (bool, *http.Request, httpkit.HttpMessage)
var CurrentProtection string

func SetCustomProtection(protection func(response http.ResponseWriter, request *http.Request, extras ...any) (bool, *http.Request, httpkit.HttpMessage)) {
	CustomProtection = protection
	CurrentProtection = "custom"
}

func SetProfileKey(keys ...string) {

}

func SetCustomInvalidTokenMessage(message string) {
	invalidToken = message
}
func SetCustomNotAuthorizedMessage(message string) {
	profile401 = message
}

func SetBearerTokenAuthorizationHeader() {
	CurrentProtection = "bearerTokenAuthorizationHeader"
}
