package licor

import "github.com/rs/cors"

var maxSizeFormData = 15

var corsConfig = cors.New(cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	AllowedHeaders:   []string{"Content-Type", "Authorization"},
	AllowCredentials: true,
})

func SetMaxSizeFormData(max int) {
	maxSizeFormData = max
}

func SetCors(cors *cors.Cors) {
	corsConfig = cors
}
