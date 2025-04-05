package main

import (
	"encoding/json"
	"net/http"
)

func startServer() error {
	http.HandleFunc("/daily", daily())
	return http.ListenAndServe(":8090", nil)
}

// always the same per day
// return 5 english and 5 math questions with difficulty based on the day of the week
func daily() http.HandlerFunc {

	topics := map[string]int{"math": 5, "english": 5}

	return func(w http.ResponseWriter, req *http.Request) {
		questions, err := loadOpenSAT("OpenSAT.json")
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		shuffled := shuffleSubset(req.Context(), questions, topics)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(shuffled)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}
