package main

import "net/http"

func main() {
	ScrapeJSON()

	http.HandleFunc("/", ReadJSON)
	http.ListenAndServe(":8000", nil)
}
