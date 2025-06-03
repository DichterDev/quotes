package quotes

import (
	"encoding/json"
	"net/http"

	. "github.com/segmentio/ksuid"
)

type Quote struct {
	Quote      string `json:"quote"`
	Author     KSUID  `json:"author"`
	Collection KSUID  `json:"collection"`
}

type Author struct {
	Name string `json:"name"`
}

type Collection struct {
	Name        string   `json:"name"`
	Description string   `json:"desc"`
	Authors     []Author `json:"authors[]"`
}

func addQuote(w http.ResponseWriter, r *http.Request) {
	var quote Quote
	err := json.NewDecoder(r.Body).Decode(&quote)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	insertQuote(quote)
}

func createCollection(w http.ResponseWriter, r *http.Request) {
	var collection Collection
	err := json.NewDecoder(r.Body).Decode(&collection)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	insertCollection(collection)
}
