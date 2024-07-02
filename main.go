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
		temp := getTempFromCity(city)

		jsonResponse(w, map[string]string{
			"client_ip": ip,
			"location":  city,
			"greeting":  fmt.Sprintf("Hello, %s!, the temperature is %.1f degrees Celsius in %s", name, temp, city),
		})
	})

	log.Println("listening on", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("could not start server:", err)
	}
}

func getTempFromCity(city string) float64 {
	var client http.Client

	resp, err := client.Get(fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", os.Getenv("WEATHER_API_KEY"), city))
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data map[string]interface{}

		err := json.NewDecoder(resp.Body).Decode(&data)
		log.Println(data)
		if err == nil {
			if curr, ok := data["current"].(map[string]interface{}); ok {
				if temp, ok := curr["temp_c"].(float64); ok {
					return temp
				}
			}
		}
	}

	return 99.9
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
