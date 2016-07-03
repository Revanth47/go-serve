package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

type ServeConfig struct {
	Port string
	Dir  string
	Path string
}

func getServeConfig(args []string) ServeConfig {

	dir, _ := os.Getwd()
	config := ServeConfig{"8000", dir, "/"}

	// if no args are passed, use default config
	if(len(args)>0){
		for _,element := range args {
			
			_,err := strconv.Atoi(element)
			
			if(err == nil){
				config.Port = element
			} else {
				if(element[0]!='/') {
					config.Dir = element
				}
			}
		}
	}

	return config
}

func serve(config ServeConfig) {

	pathStat,err := os.Stat(config.Dir)
	if(err!=nil) {
		log.Fatal(err)
	}

	if(pathStat.IsDir()) {
		http.Handle(config.Path, http.FileServer(http.Dir(config.Dir)))
	} else {
		http.HandleFunc(config.Path, func (w http.ResponseWriter, r *http.Request) {
		    http.ServeFile(w, r, config.Dir)
		})
	}

	fmt.Println("Serve Config \n Directory : " + path.Base(config.Dir) + " \n Path      : http://localhost:" + config.Port + config.Path + "\n")
	log.Println("Starting server on port: " + config.Port)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatal("Error ListenAndServe", err)
	}
}

func main() {
	args := os.Args[1:]
	config := getServeConfig(args)
	serve(config)
}