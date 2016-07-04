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
	port       string
	dir        string
	path       string
	disableDir bool
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
	/**************************************************
	 * Guard clause to validate if serve dir is valid *
	 **************************************************/
	file, err := os.Stat(c.dir)
	if err != nil {
		log.Fatal(err)
	}

	if file.IsDir() {
		if c.disableDir {
			http.HandleFunc(c.path, func(w http.ResponseWriter, r *http.Request){
				p := path.Clean(c.dir+r.URL.Path)
				file,err = os.Stat(p)

				/***************************************************
				 * Send a 404 if file not found or if path is dir  *
				 * Implemented if disable-dir flag is set to true  *
				 ***************************************************/
				if err!=nil || file.IsDir(){
					http.NotFound(w,r)
				} else {
					http.ServeFile(w,r,p)
				}
			})
		} else {
			/*************************************************************************
		     * Prefer to use http's default directory                                *
			 * if no custom modifications are required(eg disable directory listing) *
			 *************************************************************************/ 
			http.Handle(c.path, http.FileServer(http.Dir(c.dir)))
		}
	} else { 
		http.HandleFunc(c.path, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, c.dir)
		})
	}

	fmt.Println("\nServe Config \n Directory : " + path.Base(c.dir) + " \n Path      : http://localhost:" + c.port + c.path + "\n")
	log.Println("Starting server on port: " + c.port)

	log.Fatal(http.ListenAndServe(":"+c.port, logger(http.DefaultServeMux)))
}

/**************************************
 * Set default config for all values  *
 * to allow lazy execution of command *
 **************************************/
func main() {
	c := config{}
	flag.StringVar(&c.port, "p", "8000", "Port Number")
	flag.StringVar(&c.dir, "d", ".", "Serve Directory")
	flag.StringVar(&c.path, "path", "/", "Public Access Path")
	flag.BoolVar(&c.disableDir, "disable-dir", false, "Disable Directory Listing(useful for asset serving etc)")
	flag.Parse()
	c.serve()
}
