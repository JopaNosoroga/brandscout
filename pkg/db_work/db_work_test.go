package dbwork

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestCreatePostgresDb(t *testing.T) {
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	var config PostgresDbParams

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	postgres, err := CreatePostgresDb(config)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println(postgres)
	fmt.Println()
	defer postgres.Close()
}

func TestCreateQote(t *testing.T) {
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	var config PostgresDbParams

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	postgres, err := CreatePostgresDb(config)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer postgres.Close()

	quote := Quote{Author: "Mikhail", Quote: "Yes"}

	err = postgres.CreateQuote(quote)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetAllQotes(t *testing.T) {
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	var config PostgresDbParams

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	postgres, err := CreatePostgresDb(config)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer postgres.Close()

	_, err = postgres.GetAllQuotes()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetRandomQote(t *testing.T) {
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	var config PostgresDbParams

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	postgres, err := CreatePostgresDb(config)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer postgres.Close()

	_, err = postgres.GetRandomQuote()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetAuthorQotes(t *testing.T) {
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	var config PostgresDbParams

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	postgres, err := CreatePostgresDb(config)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer postgres.Close()

	_, err = postgres.GetAuthorQuotes("Mikhail")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestDeleteQote(t *testing.T) {
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	var config PostgresDbParams

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	postgres, err := CreatePostgresDb(config)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer postgres.Close()

	err = postgres.DeleteQuote(1)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
