package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"text/template"
)

type ProductCard struct {
	ImageLink    string
	ProductTitle string
	Price        string
	Available    bool
}

type PageData struct {
	Products   []ProductCard
	Pagination PaginationData
}

type PaginationData struct {
	CurrentPage  int
	TotalPages   int
	HasPrevious  bool
	HasNext      bool
	PreviousPage int
	NextPage     int
}

const ITEMS_PER_PAGE = 50

func renderTemplate(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	currentPage := 1
	if pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			currentPage = page
		}
	}

	var allProducts []ProductCard
	infoStructs := ReadJson()
	var imageLink string
	for _, v := range infoStructs {
		for _, product := range v.Products {
			for _, variant := range product.Variants {
				if variant.FeaturedImage.Src != "null" && len(product.Variants) > 1 {
					imageLink = variant.FeaturedImage.Src
				} else if len(product.Images) > 0 {
					imageLink = product.Images[0].Src
				}
				allProducts = append(allProducts, ProductCard{
					ImageLink:    imageLink,
					ProductTitle: product.Title,
					Price:        variant.Price,
					Available:    variant.Available,
				})
			}
		}
	}

	totalItems := len(allProducts)
	totalPages := int(math.Ceil(float64(totalItems) / float64(ITEMS_PER_PAGE)))

	if currentPage > totalPages {
		currentPage = totalPages
	}

	startIndex := (currentPage - 1) * ITEMS_PER_PAGE
	endIndex := startIndex + ITEMS_PER_PAGE
	if endIndex > totalItems {
		endIndex = totalItems
	}

	pageProducts := allProducts[startIndex:endIndex]

	pageData := PageData{
		Products: pageProducts,
		Pagination: PaginationData{
			CurrentPage:  currentPage,
			TotalPages:   totalPages,
			HasPrevious:  currentPage > 1,
			HasNext:      currentPage < totalPages,
			PreviousPage: currentPage - 1,
			NextPage:     currentPage + 1,
		},
	}

	tmpl, err := template.ParseFiles(
		"./views/index.html",
		"./views/product_card.html",
		"./views/pagination.html",
	)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "index", pageData)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
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

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", renderTemplate)
	http.HandleFunc("/api/refresh", refreshData)
	http.HandleFunc("/api/products", serveInfo)
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
