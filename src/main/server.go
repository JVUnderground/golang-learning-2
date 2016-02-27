package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func ServeStatic(router *mux.Router, staticDirectory string) {
	staticPaths := map[string]string{
		"styles":  staticDirectory + "/styles/",
		"images":  staticDirectory + "/images/",
		"scripts": staticDirectory + "/scripts/",
	}
	for pathName, pathValue := range staticPaths {
		pathPrefix := "/" + pathName + "/"
		router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix,
			http.FileServer(http.Dir(pathValue))))
	}
}

func main() {

	staticDirectory := "/static/"
	router := MainRouter()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	ServeStatic(router, staticDirectory)
	log.Fatal(http.ListenAndServe(":8080", router))
}
