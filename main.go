package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sajallimbu/go_securing_api/routes"
)

func main() {
	// // our route handler
	// http.Handle("/", routes.Handlers())
	fmt.Println("Listening on port 8080: ...")
	log.Fatal(http.ListenAndServe(":8080", routes.Handlers()))
}
