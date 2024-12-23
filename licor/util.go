package licor

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"

	dtoRequest "github.com/dalton02/licor/dto"
)

func extractQueryParams[Q any](r *http.Request) (Q, bool, map[string]string) {
	var params Q
	maping := make(map[string]string)
	t := reflect.TypeOf(params)
	if t == nil {
		return params, false, maping
	}
	v := reflect.ValueOf(&params).Elem()
	query := r.URL.Query()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		currentQuery := query.Get(field.Tag.Get("query"))
		if query.Get(field.Tag.Get("query")) != "" && value.Kind() == reflect.String {
			value.SetString(currentQuery)
			maping[field.Tag.Get("query")] = currentQuery
		}
	}
	return params, true, maping
}

func extractParams(pattern string, actual string) (dtoRequest.Params, error) {

	allParams := dtoRequest.Params{
		Param: make(map[string]string),
		Count: 0,
	}

	currentName := false
	var name string = ""
	var allNames []string
	var allPositions []int
	var countBarras int = 0
	for i := 0; i < len(pattern); i++ {
		if currentName && pattern[i] != '}' {
			name = name + string(pattern[i])
		}
		if pattern[i] == '/' {
			countBarras++
		}
		if pattern[i] == '{' {
			currentName = true
			allPositions = append(allPositions, countBarras)
		} else if pattern[i] == '}' {
			currentName = false
			allNames = append(allNames, name)
			name = ""
		}
	}

	if len(allPositions) == 0 {
		return allParams, errors.New("nenhum parametro encontrado")
	}

	var allValues []string
	currentBarra := allPositions[0]
	countBarras = 0
	currentCount := 0
	name = ""

	for i := 0; i < len(actual); i++ {
		if countBarras == currentBarra && actual[i] != '/' {
			name += string(actual[i])
		} else if countBarras == currentBarra && actual[i] == '/' && currentCount < len(allPositions)-1 {
			currentCount++
			currentBarra = allPositions[currentCount]
		}
		if actual[i] == '/' || i == len(actual)-1 {
			countBarras++
			if name != "" {
				allValues = append(allValues, name)
			}
			name = ""
		}

	}
	for i := 0; i < len(allNames); i++ {
		allParams.Param[allNames[i]] = allValues[i]
	}
	allParams.Count = len(allNames)
	return allParams, nil

}

func isSameRequest(expectedMethod string, req *http.Request) bool {
	if expectedMethod == "FORMDATA" && req.Method != "GET" {
		return true
	} else if expectedMethod == "FORMDATA" {
		return false
	}
	if req.Method != expectedMethod {
		return false
	}
	return true
}

func limitFormData(limit int, request *http.Request) bool {
	maxSize := (10 << limit)
	contentLength := request.Header.Get("Content-Length")
	if contentLength != "" {
		fileSize, _ := strconv.Atoi(contentLength)
		if fileSize > maxSize {
			return false
		}
	}
	return true
}
