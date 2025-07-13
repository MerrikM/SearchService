package model

type SearchFilters struct {
	ProductName string   `json:"product_name"`
	Brand       string   `json:"brand"`
	Category    string   `json:"category"`
	MinPrice    *float64 `json:"min_price"`
	MaxPrice    *float64 `json:"max_price"`
	InStockOnly bool     `json:"in_stock_only"`
}
