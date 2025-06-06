package main

import (
	"brandscout/pkg/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	err := handlers.StartDb()
	if err != nil {
		fmt.Println(err)
		return
	}

	router := mux.NewRouter()

	router.HandleFunc("/quotes", handlers.CreateQuote).Methods("POST")
	router.HandleFunc("/quotes", handlers.GetAllQuotes).Methods("GET")
	router.HandleFunc("/quotes/random", handlers.GetRandomQuote).Methods("GET")
	router.HandleFunc("/quotes/{id}", handlers.DeleteQuote).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
