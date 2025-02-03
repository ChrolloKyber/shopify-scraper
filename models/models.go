package models

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

type Info struct {
	Products []struct {
		Title       string   `json:"title"`
		Vendor      string   `json:"vendor"`
		ProductType string   `json:"product_type"`
		Tags        []string `json:"tags"`
		Handle      string   `json:"handle"`
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
