package main

import (
	"encoding/json"
	"net/http"
)

func startServer(port string) error {
	http.HandleFunc("/api/daily", daily())
	http.Handle("/", http.FileServer(http.Dir("./public")))

	return http.ListenAndServe(port, nil)
}

// always the same per day
func daily() http.HandlerFunc {

	topics := map[string]int{"math": 2, "english": 2}

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
