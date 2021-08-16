package main

import (
	"fmt"
	"log"
	"net/http"
	"silence"
	"time"
)

func main() {

	fmt.Println("Welcome to Silent HN Client")
	c, err := silence.NewConfiguration()
	if err != nil {
		log.Println("Failed to get configuration")
	}

	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
	server := http.Server{
		Addr:         addr,
		Handler:      getHandler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Printf("Listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())

}

func getHandler() http.Handler {
	app := silence.NewApp()
	router := http.NewServeMux()
	router.Handle("/", app)
	return app

}
