package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type FeaturedImage struct {
	Link string `json:"src"`
}

type Variant struct {
	Title         string        `json:"title"`
	Price         string        `json:"price"`
	FeaturedImage FeaturedImage `json:"featured_image"`
}

type Product struct {
	Title   string    `json:"title"`
	Variant []Variant `json:"variants"`
	Images  []Images  `json:"images"`
}

type Images struct {
	Position int    `json:"position"`
	Link     string `json:"src"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
}

type VariantTemplate struct {
	ImageLink    string
	ProductTitle string
	VariantName  string
	Price        string
}

type TemplateData struct {
	Variants []VariantTemplate
}

func ReadJSON(w http.ResponseWriter, r *http.Request) {
	var variants []VariantTemplate
	files, err := os.ReadDir("./json/")
	if err != nil {
		fmt.Printf("Error reading directory: %s", err)
		http.Error(w, "Error reading the directory: ", http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		filePath := fmt.Sprintf("./json/%s", file.Name())

		jsonFile, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error while reading: ", err)
			http.Error(w, "Error reading the files", http.StatusInternalServerError)
			return
		}

		var products ProductsResponse
		err = json.Unmarshal(jsonFile, &products)

		if err != nil {
			fmt.Printf("Error unmarshalling the JSON file: %s\n", err)
			http.Error(w, "Error unmarshalling the JSON", http.StatusInternalServerError)
			return
		}

		/* for _, product := range products.Products {
			fmt.Printf("Product name - %s\n", product.Title)
			for _, variant := range product.Variant {
				if variant.FeaturedImage.Link != "" {
					fmt.Printf("\tName - %s, Price - %s, Link - %s\n", variant.Title, variant.Price, variant.FeaturedImage.Link)
				} else if len(product.Images) != 0 {
					fmt.Printf("\tName - %s, Price - %s, Link - %s\n", variant.Title, variant.Price, product.Images[0].Link)
				}
			}
		} */
		for _, product := range products.Products {
			for _, variant := range product.Variant {
				// Determine the image link
				var imageLink string
				if variant.FeaturedImage.Link != "" {
					imageLink = variant.FeaturedImage.Link
				} else if len(product.Images) > 0 {
					imageLink = product.Images[0].Link
				}

				// Populate the VariantTemplate struct
				variantTemplate := VariantTemplate{
					ProductTitle: product.Title,
					VariantName:  variant.Title,
					ImageLink:    imageLink,
					Price:        variant.Price,
				}

				variants = append(variants, variantTemplate)
			}
		}
	}
	tmpl, err := template.New("index").ParseFiles("./views/index.html")
	if err != nil {
		fmt.Println("Error parsing the template: ", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		Variants: variants,
	}

	w.Header().Set("Content-Type", "text/html")

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("Error executing template: ", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
