package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ChrolloKryber/shopify-scraper/models"
)

func ReadJson() []models.Info {
	dir, err := os.ReadDir("json")
	if err != nil {
		fmt.Println("Error reading directory: ", err)
	}

	infoStruct := []models.Info{}

	for _, v := range dir {
		jsonFile, err := os.ReadFile(fmt.Sprintf("./json/%s", v.Name()))
		if err != nil {
			fmt.Println("Error reading the file: ", err)
		}

		var Products models.Info

		err = json.Unmarshal(jsonFile, &Products)

		if err != nil {
			fmt.Println("Error unmarshalling the file: ", err)
		}
		infoStruct = append(infoStruct, Products)
	}
	return infoStruct
}
