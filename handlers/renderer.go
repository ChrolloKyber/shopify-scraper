package main

import (
	"html/template"
	"math"
	"net/http"
	"strconv"

	"github.com/ChrolloKryber/shopify-scraper/models"
)

// Handler function to render the index page
func renderTemplate(w http.ResponseWriter, r *http.Request) {
	tag, vendor, search, currentPage := parseQueryParams(r)

	allProducts := loadProducts()
	if allProducts == nil {
		http.Error(w, "Failed to load products", http.StatusInternalServerError)
		return
	}

	filters := generateFilters(allProducts, tag, vendor)
	filteredProducts, totalPages := applyFilteringAndPagination(allProducts, tag, vendor, search, currentPage)

	pageData := preparePageData(filteredProducts, filters, search, currentPage, totalPages, tag, vendor)
	renderHTML(w, pageData)
}

func parseQueryParams(r *http.Request) (string, string, string, int) {
	tag := r.URL.Query().Get("tag")
	vendor := r.URL.Query().Get("vendor")
	search := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")

	currentPage := 1
	if pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			currentPage = page
		}
	}
	return tag, vendor, search, currentPage
}

func generateFilters(products []models.ProductCard, tag, vendor string) models.FilterData {
	filters := getUniqueFilters(products)
	filters.Active.Tag = tag
	filters.Active.Vendor = vendor
	return filters
}

func applyFilteringAndPagination(products []models.ProductCard, tag, vendor, search string, currentPage int) ([]models.ProductCard, int) {
	filteredProducts := filterProducts(products, tag, vendor, search)
	totalItems := len(filteredProducts)
	totalPages := int(math.Ceil(float64(totalItems) / float64(ITEMS_PER_PAGE)))
	if totalPages == 0 {
		totalPages = 1
	}

	if currentPage > totalPages {
		currentPage = totalPages
	}
	if currentPage < 1 {
		currentPage = 1
	}

	startIndex := (currentPage - 1) * ITEMS_PER_PAGE
	endIndex := min(startIndex+ITEMS_PER_PAGE, totalItems)

	return filteredProducts[startIndex:endIndex], totalPages
}

func preparePageData(products []models.ProductCard, filters models.FilterData, search string, currentPage, totalPages int, tag, vendor string) models.PageData {
	return models.PageData{
		Products: products,
		Pagination: models.PaginationData{
			CurrentPage:  currentPage,
			TotalPages:   totalPages,
			HasPrevious:  currentPage > 1,
			HasNext:      currentPage < totalPages,
			PreviousPage: currentPage - 1,
			NextPage:     currentPage + 1,
			Tag:          tag,
			Vendor:       vendor,
		},
		Filters:     filters,
		SearchQuery: search,
	}
}

func renderHTML(w http.ResponseWriter, pageData models.PageData) {
	tmpl, err := template.ParseFiles(
		"./views/index.html",
		"./views/product_card.html",
		"./views/pagination.html",
		"./views/filters.html",
	)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "index", pageData)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
