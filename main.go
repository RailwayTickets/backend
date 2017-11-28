package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/RailwayTickets/backend-go/controller"
	"github.com/RailwayTickets/backend-go/entity"
	h "github.com/RailwayTickets/backend-go/handler"
)

func main() {
	http.HandleFunc("/", hello)
	http.Handle("/register", h.Chain(http.HandlerFunc(registerHandler),
		h.SetContentTypeJSON,
		h.RequiredPost))
	http.Handle("/login", h.Chain(http.HandlerFunc(loginHandler),
		h.SetContentTypeJSON,
		h.RequiredPost))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	registrationInfo := new(entity.RegistrationInfo)
	err := json.NewDecoder(r.Body).Decode(registrationInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := registrationInfo.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	creds, err := controller.Register(registrationInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(creds)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	loginInfo := new(entity.LoginInfo)
	err := json.NewDecoder(r.Body).Decode(loginInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := loginInfo.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	creds, err := controller.Login(loginInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(creds)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}
