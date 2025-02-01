package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Info struct {
	Products []struct {
		Title       string   `json:"title"`
		Vendor      string   `json:"vendor"`
		ProductType string   `json:"product_type"`
		Tags        []string `json:"tags"`
		Variants    []struct {
			Title         string `json:"title"`
			Price         string `json:"price"`
			Available     bool   `json:"available"`
			FeaturedImage struct {
				Src string `json:"src"`
			} `json:"featured_image"`
		} `json:"variants"`
		Images []struct {
			Src string `json:"src"`
		} `json:"images"`
	} `json:"products"`
}

func ReadJson() []Info {
	dir, err := os.ReadDir("json")
	if err != nil {
		fmt.Println("Error reading directory: ", err)
	}

	infoStruct := []Info{}

	for _, v := range dir {
		jsonFile, err := os.ReadFile(fmt.Sprintf("./json/%s", v.Name()))
		if err != nil {
			fmt.Println("Error reading the file: ", err)
		}

		var Products Info

		err = json.Unmarshal(jsonFile, &Products)

		if err != nil {
			fmt.Println("Error unmarshalling the file: ", err)
		}
		infoStruct = append(infoStruct, Products)
	}
	return infoStruct
}
