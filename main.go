package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func serveInfo(w http.ResponseWriter, r *http.Request) {
	infoStructs := ReadJson()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(infoStructs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading JSON: %v", err), http.StatusInternalServerError)
	}
}

func refreshData(w http.ResponseWriter, r *http.Request) {
	DownloadJSON()
	serveInfo(w, r)
}

func main() {
	// TODO: Implement server
	http.HandleFunc("/api/refresh", refreshData)
	http.HandleFunc("/api/products", serveInfo)
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
