package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prashantgupta24/firewalld-rest/route"
)

func main() {
	fmt.Println("starting application")

	router := route.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
