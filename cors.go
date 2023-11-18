package resty

import (
	"github.com/rs/cors"
	"net/http"
)

var corsAllowedOrigins []string
var corsAllowedMethods []string
var corsAllowedHeaders []string

func SetCors(allowedOrigins, allowedMethods, allowedHeaders []string) {
	corsAllowedOrigins = allowedOrigins
	corsAllowedMethods = allowedMethods
	corsAllowedHeaders = allowedHeaders
}

func setCors(handler *handler) http.Handler {
	co := cors.New(cors.Options{
		AllowedOrigins:   corsAllowedOrigins,
		AllowedMethods:   corsAllowedMethods,
		AllowedHeaders:   corsAllowedHeaders,
		AllowCredentials: true,
	})

	return co.Handler(handler)
}
