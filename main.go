package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ChrolloKryber/shopify-scraper/limiter"
	"github.com/ChrolloKryber/shopify-scraper/models"
)

const ITEMS_PER_PAGE = 51

// Reads sites.csv and returns a mapping of site names to domains
func readSites() map[string]string {
	file, err := os.Open("sites.csv")
	if err != nil {
		log.Printf("Error opening sites.csv: %v", err)
		return nil
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Error reading sites.csv: %v", err)
		return nil
	}

	siteDomainMap := make(map[string]string)
	for _, record := range records[1:] { // Skipping the header
		if len(record) >= 2 {
			siteDomainMap[record[0]] = record[1]
		}
	}
	return siteDomainMap
}

// Extracts unique filters from the product list
func getUniqueFilters(products []models.ProductCard) models.FilterData {
	tagSet := make(map[string]bool)
	vendorSet := make(map[string]bool)

	for _, product := range products {
		for _, tag := range product.Tags {
			tagSet[tag] = true
		}
		vendorSet[product.Vendor] = true
	}

	tags := sortedKeys(tagSet)
	vendors := sortedKeys(vendorSet)

	return models.FilterData{Tags: tags, Vendors: vendors}
}

// Utility function to extract sorted keys from a map
func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Filters products based on tag, vendor, and search query
func filterProducts(products []models.ProductCard, tag, vendor, search string) []models.ProductCard {
	var filtered []models.ProductCard

	for _, product := range products {
		matchesTag := tag == "" || contains(product.Tags, tag)
		matchesVendor := vendor == "" || product.Vendor == vendor
		matchesSearch := search == "" ||
			strings.Contains(strings.ToLower(product.ProductTitle), strings.ToLower(search)) ||
			strings.Contains(strings.ToLower(product.Vendor), strings.ToLower(search))

		if matchesTag && matchesVendor && matchesSearch {
			filtered = append(filtered, product)
		}
	}

	return filtered
}

// Checks if a slice contains a specific element
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Loads JSON data from files, downloading them if necessary
func loadProducts() []models.ProductCard {
	if _, err := os.Stat("json"); os.IsNotExist(err) {
		os.Mkdir("json", 0755)
	}

	jsonDir, _ := os.ReadDir("json")
	if len(jsonDir) == 0 {
		fmt.Println("No JSON data found. Downloading...")
		DownloadJSON()
		jsonDir, _ = os.ReadDir("json")
	}

	// Read site-domain mappings
	siteDomainMap := readSites()
	if siteDomainMap == nil {
		log.Println("No sites.csv data available.")
		return nil
	}

	var allProducts []models.ProductCard

	// Load JSON files
	for _, entry := range jsonDir {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		filename := entry.Name()
		siteName := strings.TrimSuffix(filename, ".json")
		domain, exists := siteDomainMap[siteName]
		if !exists {
			continue // Skip if site is not in CSV
		}

		jsonFile, err := os.ReadFile("json/" + filename)
		if err != nil {
			log.Printf("Error reading JSON file %s: %v", filename, err)
			continue
		}

		var info models.Info
		if err := json.Unmarshal(jsonFile, &info); err != nil {
			log.Printf("Error unmarshalling JSON file %s: %v", filename, err)
			continue
		}

		// Convert products into ProductCard format
		for _, product := range info.Products {
			for _, variant := range product.Variants {
				var imageLink string
				var productTitle string
				if len(product.Variants) > 0 && variant.FeaturedImage.Src != "" {
					imageLink = variant.FeaturedImage.Src
				} else if len(product.Images) > 0 {
					imageLink = product.Images[0].Src
				}

				if variant.Title != "Default Title" {
					productTitle = fmt.Sprintf("%s - %s", product.Title, variant.Title)
				} else {
					productTitle = product.Title
				}

				allProducts = append(allProducts, models.ProductCard{
					ImageLink:    imageLink,
					ProductTitle: productTitle,
					Price:        variant.Price,
					Available:    variant.Available,
					Tags:         product.Tags,
					Vendor:       product.Vendor,
					Handle:       product.Handle,
					Domain:       domain,
				})
			}
		}
	}

	return allProducts
}

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

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", limiter.PerClientRateLimiter(renderTemplate))

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
