/**************************************************
 * @author  Revanth M(@Revanth47)                 *
 * @github  https://github.com/Revanth47/go-serve *
 * @license MIT                                   * 
 **************************************************/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type config struct {
	port       string
	dir        string
	public     string
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
	/*********************************************************
     * Status is set to http.StatusOK initially              *
     * It is then replaced by correct status in WriteHeader  *
	 *********************************************************/
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

func (c *config) clean() {
	/**************************************************
	 * Guard clause to validate if serve dir is valid *
	 **************************************************/
	_, err := os.Stat(c.dir)
	if err != nil {
		log.Fatal(err)
	}

	/**************************************************************
	 * public path must be of the form "/something/"              *
	 * directory path must be clean to ensure proper file serving *
	 **************************************************************/
	c.public = path.Join("/", c.public)
	if c.public != "/" {
		c.public = c.public + "/"
	}

	c.dir = path.Clean(c.dir)
	c.port = ":" + c.port
}

func (c config) serve() {

	http.HandleFunc(c.public, func(w http.ResponseWriter, r *http.Request) {

		/********************************************
		 * mockup of http's StripPrefix, used to    *
		 * avoid writing a handler around this func *
		 ********************************************/
		if(c.public!="/") {
			r.URL.Path = strings.TrimPrefix(r.URL.Path,c.public)
		}

		p := path.Clean(c.dir + r.URL.Path)
		file, err := os.Stat(p)
		
		if err != nil {
			http.NotFound(w, r)
		} else if file.IsDir() && c.disableDir {
			index, err := filepath.Glob("index.htm*")
			if err == nil && len(index) > 0 {
				http.ServeFile(w, r, index[0])
			}
			http.NotFound(w, r)
		} else {
			http.ServeFile(w, r, p)
		}
	})

	fmt.Println("\nServe Config \n Directory : " + path.Base(c.dir) + " \n Path      : http://localhost" + c.port + c.public + "\n")
	log.Println("Starting server")

	log.Fatal(http.ListenAndServe(c.port, logger((http.DefaultServeMux))))
}

/**************************************
 * Set default config for all values  *
 * to allow lazy execution of command *
 **************************************/
func main() {
	c := config{}
	flag.StringVar(&c.port, "p", "8000", "Port Number")
	flag.StringVar(&c.dir, "d", ".", "Serve Directory")
	flag.StringVar(&c.public, "public", "/", "Public Access Path")
	flag.BoolVar(&c.disableDir, "disable-dir", false, "Disable Directory Listing (useful for asset serving .etc)")
	flag.Parse()
	/*****************************************************
	 * Arguements are cleaned and validated to ensure    *
	 * proper arguements were passed                     *
	 *****************************************************/
	c.clean()
	c.serve()
}
