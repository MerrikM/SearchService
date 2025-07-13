package handler

import (
	"SearchService/internal/model"
	"SearchService/internal/ports"
	"encoding/json"
	"net/http"
	"strconv"
)

type SearchHandler struct {
	searcher ports.AdvertisementSearcher
}

func NewSearchHandler(searcher ports.AdvertisementSearcher) *SearchHandler {
	return &SearchHandler{searcher: searcher}
}

func (handler *SearchHandler) SearchInElastic(writer http.ResponseWriter, request *http.Request) {
	filters := model.SearchFilters{}

	query := request.URL.Query()

	filters.ProductName = query.Get("product_name")
	filters.Brand = query.Get("brand")
	filters.Category = query.Get("category")
	filters.InStockOnly = query.Get("in_stock_only") == "true"

	if minPrice := query.Get("min_price"); minPrice != "" {
		if value, err := strconv.ParseFloat(minPrice, 64); err != nil {
			filters.MinPrice = &value
		}
	}

	if maxPrice := query.Get("max_price"); maxPrice != "" {
		if val, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters.MaxPrice = &val
		}
	}

	results, err := handler.searcher.SearchAdvertisements(request.Context(), filters)
	if err != nil {
		http.Error(writer, "Ошибка поиска: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(results)
}
