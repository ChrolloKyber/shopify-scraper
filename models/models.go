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
