package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"syscall"
	"time"
)

type ByteSize uint64

const (
	_           = iota
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
)

func (b ByteSize) String() string {
	switch {
	case b >= TB:
		return fmt.Sprintf("%dTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%dGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%dMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%dKB", b/KB)
	}
	return fmt.Sprintf("%dB", b)
}

var logger = log.New(os.Stdout, "healthcheck-app", log.LstdFlags)

type settingsHandler struct{}
type rootHandler struct{}
type helloHandler struct{}
type mountHandler struct{}
type mountPermissionsHandler struct{}
type hostNameHandler struct{}

type Settings struct {
	SleepDuration int    `json:"sleepDuration"`
	Healthy       bool   `json:"healthy"`
	Name          string `json:"name"`
}

var (
	mutex    sync.RWMutex
	settings = &Settings{}
)

func main() {
	settings.SleepDuration = 0
	settings.Healthy = true
	settings.Name = os.Getenv("CATALYZE_JOB_ID")
	http.Handle("/settings", &settingsHandler{})
	http.Handle("/", &rootHandler{})
	http.Handle("/hello", &helloHandler{})
	http.Handle("/mount", &mountHandler{})
	http.Handle("/mountPermissions", &mountPermissionsHandler{})
	http.Handle("/hostname", &hostNameHandler{})
	srv := &http.Server{
		Addr: ":" + os.Getenv("PORT"),
	}
	log.Fatal(srv.ListenAndServe())
}

func (s *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}
	mutex.RLock()
	healthy := settings.Healthy
	duration := settings.SleepDuration
	mutex.RUnlock()

	time.Sleep(time.Duration(duration) * time.Second)
	str := "passing"
	if healthy {
		w.WriteHeader(200)
	} else {
		str = "failing"
		w.WriteHeader(500)
	}
	logger.Printf("Received request: %s %s - %s\n", r.Method, r.RequestURI, str)
	w.Write([]byte(str + "-" + os.Getenv("CATALYZE_JOB_ID")))
}

func (s *settingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	logger.Printf("Received admin request: %s %s\n", r.Method, r.RequestURI)
	if r.Method == http.MethodGet {
		mutex.RLock()
		defer mutex.RUnlock()
		b, _ := json.Marshal(settings)
		w.WriteHeader(200)
		w.Write(b)
	} else if r.Method == http.MethodPost {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("data: %+v\n", data)
		fmt.Printf("err: %+v\n", err)

		tempSettings := Settings{}
		err = json.Unmarshal(data, &tempSettings)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		mutex.Lock()
		defer mutex.Unlock()
		settings.Healthy = tempSettings.Healthy
		settings.SleepDuration = tempSettings.SleepDuration

		b, _ := json.Marshal(settings)
		w.WriteHeader(200)
		w.Write(b)
	}
}

func (s *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(200)
	str := fmt.Sprintf("Hello %s", r.Header.Get("X-Forwarded-For"))
	logger.Printf("Received request: %s %s - %s\n", r.Method, r.RequestURI, str)
	w.Write([]byte(str))
}

func (s *mountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	res := new(syscall.Statfs_t)
	err := syscall.Statfs("/data", res)
	if err != nil {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("%s", ByteSize(res.Bavail*uint64(res.Bsize)))))
	}
}

func (s *mountPermissionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	res := new(syscall.Stat_t)
	err := syscall.Stat("/data", res)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	if res.Uid != 0 {
		w.WriteHeader(500)
		w.Write([]byte("Uid was not root"))
	}
	if res.Gid != 0 {
		w.WriteHeader(500)
		w.Write([]byte("Gid was not root"))
	}
	pathToTmpFile := "/data/test.txt"
	newFile, err := os.Create(pathToTmpFile)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	_, err = newFile.WriteString("nailed it")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	err = os.Remove(pathToTmpFile)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(204)
}

func (s *hostNameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	name, err := os.Hostname()
	if err != nil {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
		w.Write([]byte(name))
	}
}
