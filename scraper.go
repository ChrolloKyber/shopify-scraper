package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func DownloadJSON() {
	file, err := os.Open("sites.csv")
	if err != nil {
		fmt.Printf("Error opening the CSV file: %s", err)
	}

	defer file.Close()
	csvRead := csv.NewReader(file)

	data, err := csvRead.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV file: %s", err)
	}

	baseURL := "https://"
	product := "/products.json"

	// TODO: implement go routines
	for k, v := range data {
		if k == 0 {
			continue
		}

		fmt.Println(k, v[0], v[1])

		file, err := http.Get(baseURL + v[1] + product)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Body.Close()

		out, err := os.Create("json/" + v[0] + ".json")
		if err != nil {
			fmt.Println(err)
		}
		defer out.Close()

		if file.StatusCode != http.StatusOK {
			fmt.Printf("Error writing: %s", file.Status)
		}

		_, err = io.Copy(out, file.Body)
		if err != nil {
			fmt.Printf("Error writing to the file: %s", err)
		}
	}
	fmt.Println("All tasks completed.")
}

func main() {
	DownloadJSON()
}
