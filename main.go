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

		ip := getIPAddressFromRequest(r)
		city := getCityFromIpAddress(ip)

		jsonResponse(w, map[string]string{
			"client_ip": ip,
			"location":  city,
			"greeting":  fmt.Sprintf("Hello, %s!, the temperature is 11 degrees Celsius in %s", name, city),
		})
	})

	log.Println("listening on", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("could not start server:", err)
	}
}

func getIPAddressFromRequest(r *http.Request) string {
	ips := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
	if len(ips) > 0 {
		return net.ParseIP(ips[0]).String()
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func getCityFromIpAddress(ip string) string {
	var client http.Client

	resp, err := client.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	if err != nil {
		return "unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data struct {
			Status string
			City   string
		}
		err := json.NewDecoder(resp.Body).Decode(&data)
		if err == nil && data.Status == "success" && data.City != "" {
			return data.City
		}
	}

	return "unknown"
}

func jsonResponse(w http.ResponseWriter, v map[string]string) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
