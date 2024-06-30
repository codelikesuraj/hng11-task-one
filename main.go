package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
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
			"client_ip": getIPAddressFromRequest(r),
			"location":  "New York",
			"greeting":  fmt.Sprintf("Hello, %s!, the temperature is 11 degrees Celsius in New York", name),
		})
	})

	log.Println("listening on", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("could not start server:", err)
	}
}

func getIPAddressFromRequest(r *http.Request) string {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return ""
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		return netIP.String()
	}

	return ""
}

func jsonResponse(w http.ResponseWriter, v map[string]string) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
