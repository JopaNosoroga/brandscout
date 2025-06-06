package handlers

import (
	dbwork "brandscout/pkg/db_work"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var postgres dbwork.DB

func StartDb() error {
	configFile, err := os.ReadFile("pkg/db_work/config.json")
	if err != nil {
		return err
	}

	var config dbwork.PostgresDbParams

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return err
	}

	postgres, err = dbwork.CreatePostgresDb(config)
	if err != nil {
		return err
	}

	return nil
}

func CreateQuote(rw http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	quota := dbwork.Quote{}

	err = json.Unmarshal(data, &quota)
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = postgres.CreateQuote(quota)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func GetAllQuotes(rw http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	quotes := make([]dbwork.Quote, 0)
	var err error
	if author == "" {
		quotes, err = postgres.GetAllQuotes()
	} else {
		quotes, err = postgres.GetAuthorQuotes(author)
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(quotes)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
}

func GetRandomQuote(rw http.ResponseWriter, r *http.Request) {
	quote, err := postgres.GetRandomQuote()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(quote)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
}

func DeleteQuote(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = postgres.DeleteQuote(id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
