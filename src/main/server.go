package main

import (
	"log"
	"net/http"
)

func main() {
	styleHandler := http.FileServer(http.Dir("/Users/joao/Code/developer-shop/src/script/style"))
	http.Handle("/style/", http.StripPrefix("/style/", styleHandler))

	router := MainRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
