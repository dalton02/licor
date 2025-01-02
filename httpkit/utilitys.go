package httpkit

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	dtoRequest "github.com/dalton02/licor/dto"
	"github.com/dalton02/licor/validator"
	"github.com/golang-jwt/jwt"
)

func GetRequestBody(request *http.Request) []byte {
	body, ok := request.Context().Value("original_body").([]byte)
	if !ok {
		return nil
	}
	return body
}

func GetJsonSchema[T any](request *http.Request) map[int][]string {
	body := GetRequestBody(request)
	var data T
	jsonString := body
	json.Unmarshal(body, &data)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(jsonString), &jsonData)
	keysByLevel := make(map[int][]string)
	jsonSchema := validator.ExtractKeysByLevel(jsonData, 1, keysByLevel)
	return jsonSchema.Mapa
}

func GetBearerToken(auth string) string {
	result := strings.Replace(auth, "Bearer", "", -1)
	result = strings.TrimSpace(result)
	return result
}

// Returns a struct with a count of the params and a map[string]string to get the param
func GetUrlParams(request *http.Request) (dtoRequest.Params, error) {
	paramsInterface := request.Context().Value("params")
	params, test := paramsInterface.(dtoRequest.Params)
	if !test {
		return params, errors.New("erro ao obter parametros")
	}
	return params, nil
}

func GetDataToken(request *http.Request) (map[string]interface{}, error) {
	authorization := request.Header.Get("Authorization")
	token := GetBearerToken(authorization)
	tokenData, err := GetJwtInfo(token)
	return tokenData, err
}

func GenerateJwt[T any](data T, minutes int) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["data"] = data
	claims["exp"] = time.Now().Add(time.Duration(minutes) * time.Minute).Unix()
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetJwtInfo(tokenString string) (map[string]interface{}, error) {
	secretKey := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return make(map[string]interface{}), err
	}

	data, _ := claims["data"].(map[string]interface{})
	return data, nil
}
