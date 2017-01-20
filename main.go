package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var logger = log.New(os.Stdout, "healthcheck-app", log.LstdFlags)

type handler struct{}
type adminHandler struct{}

var (
	healthcheckPassing = true
	duration           int
)

func main() {
	duration, _ = strconv.Atoi(os.Getenv("SLEEP_DURATION"))
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
	time.Sleep(time.Duration(duration) * time.Second)
	if healthcheckPassing {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
	}
	str := "passing"
	if !healthcheckPassing {
		str = "failing"
	}
	logger.Printf("Received request: %s %s - %s\n", r.Method, r.RequestURI, str)
	w.Write([]byte(str + "-" + os.Getenv("CATALYZE_JOB_ID")))
}

func (h *adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	healthcheckPassing = r.Method != http.MethodDelete
	str := "passing"
	if !healthcheckPassing {
		str = "failing"
	}
	logger.Printf("Received admin request: %s %s - healthcheck is now %s\n", r.Method, r.RequestURI, str)
	w.WriteHeader(200)
	w.Write([]byte(str))
}
