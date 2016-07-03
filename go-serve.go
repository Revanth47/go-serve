package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"path"
)

type ServeConfig struct {
	Port string
	Dir  string
	Path string
}

func main() {
	args   := os.Args[1:]
	config := getServeConfig(args)
	serve(config)
}

func getServeConfig(args []string)(ServeConfig) {

	var config ServeConfig
	if(len(args)==0){
		dir,_ := os.Getwd()
		config = ServeConfig{"8000",dir,"/"}
	}

	return config
}

func serve(config ServeConfig) {
	
	http.Handle("/",http.FileServer(http.Dir(config.Dir)))

	fmt.Println("Serve directory: ",path.Base(config.Dir))
	log.Println("Starting server on port: "+config.Port)

	err := http.ListenAndServe(":"+config.Port,nil)
	if err != nil {
		log.Fatal("Error ListenAndServe", err)
	} 
}