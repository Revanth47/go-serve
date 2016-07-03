package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"strconv"
	"time"
)

type ServeConfig struct {
	Port string
	Dir  string
	Path string
}


type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	w.length = len(b)
	return w.ResponseWriter.Write(b)
}

// WriteLog Logs the Http Status for a request into fileHandler and returns a httphandler function which is a wrapper to log the requests.
func logger(handle http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		start := time.Now()
		writer := statusWriter{w, 0, 0}
		handle.ServeHTTP(&writer, request)
		end := time.Now()
		latency := end.Sub(start)
		statusCode := writer.status
		log.Println(request.Method, request.URL.Path, statusCode, latency)
	}
}

func getServeConfig(args []string) ServeConfig {

	dir, _ := os.Getwd()
	config := ServeConfig{"8000", dir, "/"}

	// if no args are passed, use default config
	if len(args) > 0 {
		for _, element := range args {

			_, err := strconv.Atoi(element)

			if err == nil {
				config.Port = element
			} else if(strings.Contains(element,":")) {
				s := strings.Split(element,":")
				config.Dir,config.Port = s[0],s[1] 
			} else {
				config.Dir = element	
			}
		}
	}

	return config
}

func serve(config ServeConfig) {
 	
 	pathStat, err := os.Stat(config.Dir)
	if err != nil {
		log.Fatal(err)
	}

	if pathStat.IsDir() {
		http.Handle(config.Path, http.FileServer(http.Dir(config.Dir)))
	} else {
		http.HandleFunc(config.Path, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, config.Dir)
		})
	}

	fmt.Println("\nServe Config \n Directory : " + path.Base(config.Dir) + " \n Path      : http://localhost:" + config.Port + config.Path + "\n")
	log.Println("Starting server on port: " + config.Port)

	log.Fatal(http.ListenAndServe(":"+config.Port, logger(http.DefaultServeMux)))
}

func main() {
	args := os.Args[1:]
	config := getServeConfig(args)
	serve(config)
}
