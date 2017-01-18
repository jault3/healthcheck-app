package main

import (
	"log"
	"net/http"
)

var logger = healthcheckapp.NewLogger("healthcheck-app")

type handler struct{}
type adminHandler struct{}

var healthcheckPassing = true

func main() {
	go adminServer()
	h := &handler{}
	srv := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}
	log.Fatal(srv.ListenAndServe())
}

func adminServer() {
	h := &adminHandler{}
	srv := &http.Server{
		Addr:    ":8081",
		Handler: h,
	}
	log.Fatal(srv.ListenAndServe())
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logger.Info("Received request")
	if healthcheckPassing {
		w.WriteHeader(204)
	} else {
		w.WriteHeader(500)
	}
}

func (h *adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	healthcheckPassing = r.Method != http.MethodDelete
	str := "passing"
	if !healthcheckPassing {
		str = "failing"
	}
	logger.Info("Received admin request: healthcheck is now %s", str)
	w.WriteHeader(200)
	w.Write([]byte(str))
}
