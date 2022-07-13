package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"resty/config"
	"resty/routing"
)

func init() {
	config.Init()
}

func main() {
	config.Init()
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(routing.UnknownMethod)
	router.HandleFunc("/signup", routing.SignUp).Methods("PUT")
	router.HandleFunc("/verify_user", routing.VerifyUser).Methods("POST")
	router.HandleFunc("/send_verify_code", routing.SendVerifyCode).Methods("POST")
	router.HandleFunc("/signin", routing.SignIn).Methods("POST")
	router.HandleFunc("/signout", routing.SignOut).Methods("DELETE")
	router.HandleFunc("/user/{id}", routing.GetUser).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: false,
	})
	handler := c.Handler(router)
	fmt.Println("Listening port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
