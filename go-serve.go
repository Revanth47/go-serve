package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

type ServeConfig struct {
	Port string
	Dir  string
	Path string
}

func main() {
	args := os.Args[1:]
	config := getServeConfig(args)
	serve(config)
}

func getServeConfig(args []string) ServeConfig {

	dir, _ := os.Getwd()
	config := ServeConfig{"8000", dir, "/"}

	return config
}

func serve(config ServeConfig) {

	http.Handle(config.Path, http.FileServer(http.Dir(config.Dir)))

	fmt.Println("Serve Config \n Directory : " + path.Base(config.Dir) + " \n On        : http://localhost:" + config.Port + config.Path + "\n")
	log.Println("Starting server on port: " + config.Port)

	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatal("Error ListenAndServe", err)
	}
}
