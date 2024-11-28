package main

import "net/http"

func main() {
	ScrapeJSON()

	http.HandleFunc("/", ReadJSON)
	http.HandleFunc("/css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
	})
	http.ListenAndServe(":8000", nil)
}
