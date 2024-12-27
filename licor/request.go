package licor

import (
	"fmt"
	"net/http"

	"github.com/dalton02/licor/httpkit"
	"github.com/rs/cors"
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

func nothing(response http.ResponseWriter, request *http.Request) bool {
	return true
}

type HandlerRequest[B any, Q any] struct {
	endpoint    string
	security    string
	extraRoute  []any
	middlewares []func(http.ResponseWriter, *http.Request) (httpkit.HttpMessage, bool)
	controller  func(http.ResponseWriter, *http.Request) httpkit.HttpMessage
}

func Init(porta string) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Permitir os métodos HTTP
		AllowedHeaders:   []string{"Content-Type", "Authorization"},           // Permitir cabeçalhos
		AllowCredentials: true,                                                // Permitir cookies ou autenticação
	})

	corsHandler := c.Handler(http.DefaultServeMux)

	err := http.ListenAndServe(":"+porta, corsHandler)
	fmt.Println("Server running in port:" + porta)
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
				httpkit.SendHttpMessage(httpkit.AppInternal("internal error"), response)
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
				httpkit.SendHttpMessage(httpkit.AppInternal("internal error"), response)
				return
			}
		}()

		httpkit.SendHttpMessage(generic[B, Q](response, request, r, "GET"), response)
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
				httpkit.SendHttpMessage(httpkit.AppInternal("internal error"), response)
				return
			}
		}()

		httpkit.SendHttpMessage(generic[B, Q](response, request, r, "POST"), response)
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
				httpkit.SendHttpMessage(httpkit.AppInternal("internal error"), response)
				return
			}
		}()

		httpkit.SendHttpMessage(generic[B, Q](response, request, r, "PUT"), response)
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
				httpkit.SendHttpMessage(httpkit.AppInternal("internal error"), response)
				return
			}
		}()

		httpkit.SendHttpMessage(generic[B, Q](response, request, r, "PATCH"), response)
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
				httpkit.SendHttpMessage(httpkit.AppInternal("internal error"), response)
				return
			}
		}()

		httpkit.SendHttpMessage(generic[B, Q](response, request, r, "DELETE"), response)
	})
}
