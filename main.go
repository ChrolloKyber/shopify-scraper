package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type ProductCard struct {
	ImageLink    string
	ProductTitle string
	Price        string
	Available    bool
}

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

func renderTemplate(w http.ResponseWriter, r *http.Request) {
	var ProductCards []ProductCard
	infoStructs := ReadJson()
	var imagelink string
	for _, v := range infoStructs {
		for _, product := range v.Products {
			for _, variant := range product.Variants {
				if variant.FeaturedImage.Src == "" {
					imagelink = variant.FeaturedImage.Src
				} else {
					imagelink = product.Images[0].Src
				}
				ProductCards = append(ProductCards, ProductCard{
					ImageLink:    imagelink,
					ProductTitle: fmt.Sprintf("%v - %v ", product.Title, variant.Title),
					Price:        variant.Price,
					Available:    variant.Available,
				})
			}
		}
	}

	// First, parse all template files
	tmpl, err := template.ParseFiles(
		"./views/index.html",
		"./views/product_card.html",
	)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Execute the named template "index"
	err = tmpl.ExecuteTemplate(w, "index", ProductCards)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func main() {
	// TODO: Implement server
	// TODO: Implement template and page
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", renderTemplate)
	http.HandleFunc("/api/refresh", refreshData)
	http.HandleFunc("/api/products", serveInfo)
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
