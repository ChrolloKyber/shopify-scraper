package main

import "fmt"

func main() {
	/* file, err := http.Get("https://neomacro.in/products.json")

	if err != nil {
		fmt.Printf("Failed to get the file: %s\n", err)
		return
	}
	defer file.Body.Close()

	data, err := io.ReadAll(file.Body)
	if err != nil {
		fmt.Printf("Failed to get the file: %s\n", err)
		return
	}

	var products ProductsResponse
	err = json.Unmarshal(data, &products)

	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %s\n", err)
		return
	}

	for _, product := range products.Products {
		fmt.Printf("Product Name - %s\n", product.Title)
		for _, variant := range product.Variant {
			fmt.Printf("\tName - %s, Price - %s\n", variant.Title, variant.Price)
		}
	} */
	// ScrapeJSON()
	fmt.Printf("Choose a vendor:\n\t(1) NeoMacro\n\t(2) AceKBD\n\t(3) GenesisPC\n\t(4) Keebsmod\n")

	ReadJSON()
}
