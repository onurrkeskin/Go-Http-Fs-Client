package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	FSERVER_PORT = 8081
)

func main() {
	// Simple static webserver:
	applicationPort := fmt.Sprintf("0.0.0.0:%d", FSERVER_PORT)
	log.Printf("Starting to serve files on " + applicationPort)
	log.Fatal(http.ListenAndServe(applicationPort, http.FileServer(http.Dir("/service/simple-dir/"))))
}
