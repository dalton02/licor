package licor

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

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
	const maxTentativas = 30

	// Converte a porta inicial para inteiro
	portaInt, err := strconv.Atoi(porta)
	if err != nil {
		fmt.Printf("Erro: porta inválida %s: %v\n", porta, err)
		return
	}

	for tentativa := 0; tentativa < maxTentativas; tentativa++ {
		portaAtual := portaInt + tentativa
		portaStr := strconv.Itoa(portaAtual)

		corsHandler := corsConfig.Handler(http.DefaultServeMux)

		listener, err := net.Listen("tcp", ":"+portaStr)
		if err != nil {
			if strings.Contains(err.Error(), "address already in use") {
				fmt.Printf("Porta %s em uso, tentando porta %d\n", portaStr, portaAtual+1)
				continue
			}
			fmt.Printf("Erro ao tentar usar a porta %s: %v\n", portaStr, err)
			return
		}

		fmt.Printf("Licor rodando na porta: %s\n", portaStr)
		err = http.Serve(listener, corsHandler)
		if err != nil {
			fmt.Printf("Erro no servidor na porta %s: %v\n", portaStr, err)
			return
		}
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
