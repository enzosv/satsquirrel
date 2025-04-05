package main

import (
	"encoding/json"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/questions", getQuestions())
	http.ListenAndServe(":8090", nil)
}

func getQuestions() http.HandlerFunc {

	topics := map[string]int{"math": 5, "english": 5}

	return func(w http.ResponseWriter, req *http.Request) {
		questions, err := loadOpenSAT("OpenSAT.json")
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		shuffled := randomize(questions, topics)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shuffled)
	}

}
