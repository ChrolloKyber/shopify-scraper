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
)

type ProductCard struct {
	ImageLink    string
	ProductTitle string
	Price        string
	Available    bool
	Tags         []string
	Vendor       string
	Handle       string
	Domain       string
}

type PageData struct {
	Products    []ProductCard
	Pagination  PaginationData
	Filters     FilterData
	SearchQuery string
}

type PaginationData struct {
	CurrentPage  int
	TotalPages   int
	HasPrevious  bool
	HasNext      bool
	PreviousPage int
	NextPage     int
}

type FilterData struct {
	Tags    []string
	Vendors []string
	Active  struct {
		Tag    string
		Vendor string
	}
}

const ITEMS_PER_PAGE = 50

func readSites() [][]string {
	file, err := os.Open("sites.csv")
	if err != nil {
		log.Printf("Error opening sites.csv: %v", err)
		return nil
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	sites, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Error reading sites.csv: %v", err)
		return nil
	}
	return sites[1:]
}

func getUniqueFilters(products []ProductCard) FilterData {
	tagMap := make(map[string]bool)
	vendorMap := make(map[string]bool)

	for _, product := range products {
		for _, tag := range product.Tags {
			tagMap[tag] = true
		}
		vendorMap[product.Vendor] = true
	}

	var tags []string
	var vendors []string
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	for vendor := range vendorMap {
		vendors = append(vendors, vendor)
	}
	sort.Strings(tags)
	sort.Strings(vendors)

	return FilterData{
		Tags:    tags,
		Vendors: vendors,
	}
}

func filterProducts(products []ProductCard, tag, vendor, search string) []ProductCard {
	var filtered []ProductCard

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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
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

	// Read sites.csv to build site-domain map
	sites := readSites()
	siteDomainMap := make(map[string]string)
	for _, site := range sites {
		siteDomainMap[site[0]] = site[1]
	}

	var allProducts []ProductCard

	// Read JSON directory
	jsonDir, err := os.ReadDir("json")
	if err != nil {
		log.Printf("Error reading JSON directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for _, entry := range jsonDir {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		siteName := strings.TrimSuffix(filename, ".json")
		domain, ok := siteDomainMap[siteName]
		if !ok {
			continue // Skip if site not found in CSV
		}

		jsonFile, err := os.ReadFile("./json/" + filename)
		if err != nil {
			log.Printf("Error reading JSON file %s: %v", filename, err)
			continue
		}

		var info Info
		if err := json.Unmarshal(jsonFile, &info); err != nil {
			log.Printf("Error unmarshalling JSON file %s: %v", filename, err)
			continue
		}

		for _, product := range info.Products {
			var imageLink string
			if len(product.Variants) > 0 && product.Variants[0].FeaturedImage.Src != "" {
				imageLink = product.Variants[0].FeaturedImage.Src
			} else if len(product.Images) > 0 {
				imageLink = product.Images[0].Src
			}

			for _, variant := range product.Variants {
				allProducts = append(allProducts, ProductCard{
					ImageLink:    imageLink,
					ProductTitle: fmt.Sprintf("%s - %s", product.Title, variant.Title),
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

	filters := getUniqueFilters(allProducts)
	filters.Active.Tag = tag
	filters.Active.Vendor = vendor

	filteredProducts := filterProducts(allProducts, tag, vendor, search)

	totalItems := len(filteredProducts)
	totalPages := int(math.Ceil(float64(totalItems) / float64(ITEMS_PER_PAGE)))

	if currentPage > totalPages {
		currentPage = totalPages
	}
	if currentPage < 1 {
		currentPage = 1
	}

	startIndex := (currentPage - 1) * ITEMS_PER_PAGE
	endIndex := startIndex + ITEMS_PER_PAGE
	if endIndex > totalItems {
		endIndex = totalItems
	}

	pageProducts := filteredProducts[startIndex:endIndex]

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
		Filters:     filters,
		SearchQuery: search,
	}

	tmpl, err := template.ParseFiles(
		"./views/index.html",
		"./views/product_card.html",
		"./views/pagination.html",
		"./views/filters.html",
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
	err := os.Mkdir("json", 0755)
	if err != nil {
		log.Printf("Error creating JSON directory: %v", err)
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", renderTemplate)
	http.HandleFunc("/api/refresh", refreshData)
	http.HandleFunc("/api/products", serveInfo)
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
