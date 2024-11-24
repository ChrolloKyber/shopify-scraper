package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Variant struct {
	Title string `json:"title"`
	Price string `json:"price"`
}

type Product struct {
	Title   string    `json:"title"`
	Variant []Variant `json:"variants"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
}

func ReadJSON() {
	var vendor int
	_, err := fmt.Scan(&vendor)
	if err != nil {
		fmt.Println("Error while reading choice: ", err)
	}

	var filePath string

	switch vendor {
	case 1:
		filePath = "./json/neomacro.json"
	case 2:
		filePath = "./json/acekbd.json"
	case 3:
		filePath = "./json/genesispc.json"
	case 4:
		filePath = "./json/keebsmod.json"
	}

	jsonFile, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error while reading: ", err)
	}

	var products ProductsResponse
	err = json.Unmarshal(jsonFile, &products)

	if err != nil {
		fmt.Printf("Error unmarshalling the JSON file: %s\n", err)
		return
	}

	for _, product := range products.Products {
		fmt.Printf("Product name - %s\n", product.Title)
		for _, variant := range product.Variant {
			fmt.Printf("\tName - %s, Price - %s\n", variant.Title, variant.Price)
		}
	}
}
