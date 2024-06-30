package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func jsonResponse(w http.ResponseWriter, v map[string]string) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/hello/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		name := r.URL.Query().Get("visitor_name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse(w, map[string]string{
				"message": "visitor_name field is required",
			})
			return
		}

		jsonResponse(w, map[string]string{
			"client_ip": r.RemoteAddr,
			"location":  "New York",
			"greeting":  fmt.Sprintf("Hello, %s!, the temperature is 11 degrees Celsius in New York", name),
		})
	})

	log.Println("listening on", port)
	http.ListenAndServe(":"+port, nil)
}
