package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Variants struct {
	Title string `json:"title"`
	Price string `json:"price"`
}

type Product struct {
	Title    string     `json:"title"`
	Variants []Variants `json:"variants"`
}

type Products struct {
	Product []Product `json:"products"`
}

func ReadJSON() {
	dirFiles, err := os.ReadDir("json")

	if err != nil {
		fmt.Println(err)
	}

	for _, v := range dirFiles {
		jsonFile, err := os.ReadFile(fmt.Sprintf("./json/%s", v.Name()))
		if err != nil {
			fmt.Println(err)
		}
		var Products Products
		err = json.Unmarshal(jsonFile, &Products)

		if err != nil {
			fmt.Println("Error unmarshalling the JSON: ", err)
		}

		for _, product := range Products.Product {
			fmt.Printf("Product name - %v\n", product.Title)
			for _, variant := range product.Variants {
				fmt.Printf("\tVariant - %v, Price - %v\n", variant.Title, variant.Price)
			}
		}
	}
}
