package licor

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dalton02/licor/httpkit"
)

type Requests interface {
	Post(handlerFunc func(http.ResponseWriter, *http.Request) httpkit.HttpMessage, middlewares ...func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool))
	Get(handlerFunc func(http.ResponseWriter, *http.Request) httpkit.HttpMessage, middlewares ...func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool))
	Put(handlerFunc func(http.ResponseWriter, *http.Request) httpkit.HttpMessage, middlewares ...func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool))
	Patch(handlerFunc func(http.ResponseWriter, *http.Request) httpkit.HttpMessage, middlewares ...func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool))
	Delete(handlerFunc func(http.ResponseWriter, *http.Request) httpkit.HttpMessage, middlewares ...func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool))
	FormData(handlerFunc func(http.ResponseWriter, *http.Request) httpkit.HttpMessage, middlewares ...func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool))
}

type MiddleRequest interface {
	MiddleWare(middleware ...func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool))
}

type HandlerRequest[B any, Q any] struct {
	endpoint    string
	security    string
	extraRoute  []any
	middlewares []func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool)
	controller  func(http.ResponseWriter, *http.Request) httpkit.HttpMessage
}

func Init(porta string) {

	corsHandler := corsConfig.Handler(http.DefaultServeMux)
	fmt.Println("Licor running in port: " + porta)

	err := http.ListenAndServe(":"+porta, corsHandler)
	if err != nil {
		fmt.Println("Erro no servidor: ", err)
	}
}

func Public[B any, Q any](rota string) Requests {
	return &HandlerRequest[B, Q]{
		endpoint: rota,
		security: "public",
	}
}
func Protected[B any, Q any](rota string, extra ...any) Requests {
	return &HandlerRequest[B, Q]{
		endpoint:   rota,
		security:   "protected",
		extraRoute: extra,
	}
}

func sendHttpMessage(message httpkit.HttpMessage, response http.ResponseWriter) {
	json.Marshal(message)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(message.StatusCode)
	if err := json.NewEncoder(response).Encode(message); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func (r *HandlerRequest[B, Q]) FormData(handlerFunc func(http.ResponseWriter, *http.Request) httpkit.HttpMessage, middlewares ...func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool)) {
	r.controller = handlerFunc
	r.middlewares = middlewares
	http.HandleFunc(r.endpoint, func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if errStr, ok := err.(string); ok && errStr == "common" {
					return
				}
				fmt.Println(err)
				sendHttpMessage(httpkit.AppInternal("internal error"), response)
				return
			}
		}()

		generic[B, Q](response, request, r, "FORMDATA")
	})
}

func (r *HandlerRequest[B, Q]) Get(handlerFunc func(http.ResponseWriter, *http.Request) httpkit.HttpMessage, middlewares ...func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool)) {
	r.controller = handlerFunc
	r.middlewares = middlewares
	http.HandleFunc(r.endpoint, func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if errStr, ok := err.(string); ok && errStr == "common" {
					return
				}
				fmt.Println(err)
				sendHttpMessage(httpkit.AppInternal("internal error"), response)
				return
			}
		}()

		sendHttpMessage(generic[B, Q](response, request, r, "GET"), response)
	})
}
func (r *HandlerRequest[B, Q]) Post(handlerFunc func(http.ResponseWriter, *http.Request) httpkit.HttpMessage, middlewares ...func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool)) {
	r.controller = handlerFunc
	r.middlewares = middlewares
	http.HandleFunc(r.endpoint, func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if errStr, ok := err.(string); ok && errStr == "common" {
					return
				}
				fmt.Println(err)
				sendHttpMessage(httpkit.AppInternal("internal error"), response)
				return
			}
		}()

		sendHttpMessage(generic[B, Q](response, request, r, "POST"), response)
	})
}
func (r *HandlerRequest[B, Q]) Put(handlerFunc func(http.ResponseWriter, *http.Request) httpkit.HttpMessage, middlewares ...func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool)) {
	r.controller = handlerFunc
	r.middlewares = middlewares
	http.HandleFunc(r.endpoint, func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if errStr, ok := err.(string); ok && errStr == "common" {
					return
				}
				fmt.Println(err)
				sendHttpMessage(httpkit.AppInternal("internal error"), response)
				return
			}
		}()

		sendHttpMessage(generic[B, Q](response, request, r, "PUT"), response)
	})
}
func (r *HandlerRequest[B, Q]) Patch(handlerFunc func(http.ResponseWriter, *http.Request) httpkit.HttpMessage, middlewares ...func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool)) {
	r.controller = handlerFunc
	r.middlewares = middlewares
	http.HandleFunc(r.endpoint, func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if errStr, ok := err.(string); ok && errStr == "common" {
					return
				}
				fmt.Println(err)
				sendHttpMessage(httpkit.AppInternal("internal error"), response)
				return
			}
		}()

		sendHttpMessage(generic[B, Q](response, request, r, "PATCH"), response)
	})
}
func (r *HandlerRequest[B, Q]) Delete(handlerFunc func(http.ResponseWriter, *http.Request) httpkit.HttpMessage, middlewares ...func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool)) {
	r.controller = handlerFunc
	r.middlewares = middlewares
	http.HandleFunc(r.endpoint, func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if errStr, ok := err.(string); ok && errStr == "common" {
					return
				}
				fmt.Println(err)
				sendHttpMessage(httpkit.AppInternal("internal error"), response)
				return
			}
		}()

		sendHttpMessage(generic[B, Q](response, request, r, "DELETE"), response)
	})
}
