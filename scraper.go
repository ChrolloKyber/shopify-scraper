package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func ScrapeJSON() {
	URLS := make(map[string]string)
	Vendors := make(map[string]string)

	Vendors["neomacro"] = "neomacro.in"
	Vendors["acekbd"] = "acekbd.com"
	Vendors["genesispc"] = "genesispc.in"
	Vendors["keebsmod"] = "keebsmod.com"

	for vendor, domain := range Vendors {
		URLS[vendor] = fmt.Sprintf("https://%s/products.json", domain)
	}

	for vendor, url := range URLS {
		fileName := fmt.Sprintf("./json/%s.json", vendor)
		file, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching URL: ", err)
		}
		defer file.Body.Close()

		out, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Error creating file: ", err)
		}
		defer out.Close()

		if file.StatusCode != http.StatusOK {
			fmt.Println(fmt.Errorf("Bad status: %s", file.Status))
		}

		_, err = io.Copy(out, file.Body)
		if err != nil {
			fmt.Println("Error saving file", err)
		}
	}
}
