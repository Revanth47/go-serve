package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

type config struct {
	Port    string
	Dir     string
	Path    string
	listDir bool
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
		w.status = http.StatusOK
	}
	w.length = len(b)
	return w.ResponseWriter.Write(b)
}

func logger(handle http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		start := time.Now()
		writer := statusWriter{w, 0, 0}
		handle.ServeHTTP(&writer, request)
		latency := time.Since(start)
		statusCode := writer.status
		log.Println(request.Method, request.URL.Path, statusCode, latency)
	}
}

func (c config) serve() {
	file, err := os.Stat(c.Dir)
	if err != nil {
		log.Fatal(err)
	}

	if file.IsDir() {
		http.Handle(c.Path, http.FileServer(http.Dir(c.Dir)))
	} else {
		http.HandleFunc(c.Path, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, c.Dir)
		})
	}

	fmt.Println("\nServe Config \n Directory : " + path.Base(c.Dir) + " \n Path      : http://localhost:" + c.Port + c.Path + "\n")
	log.Println("Starting server on port: " + c.Port)

	log.Fatal(http.ListenAndServe(":"+c.Port, logger(http.DefaultServeMux)))
}

func main() {
	c := config{}
	flag.StringVar(&c.Port, "p", "8000", "Port Number")
	flag.StringVar(&c.Dir, "d", ".", "Serve Directory")
	flag.StringVar(&c.Path, "path", "/", "Public Access Path")
	flag.BoolVar(&c.listDir, "disable-dir", false, "Disable Directory Listing")
	flag.Parse()
	c.serve()
}
