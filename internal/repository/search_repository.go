package repository

import (
	"SearchService/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
)

type SearchRepository struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticRepository(client *elasticsearch.Client, index string) *SearchRepository {
	return &SearchRepository{
		client: client,
		index:  index,
	}
}

func (repo *SearchRepository) SearchAdvertisements(ctx context.Context, filters model.SearchFilters) ([]model.Advertisement, error) {
	var must []map[string]interface{}
	var filter []map[string]interface{}

	// Match по названию
	if filters.ProductName != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"product_name": filters.ProductName,
			},
		})
	}

	// Match по бренду
	if filters.Brand != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"brand": filters.Brand,
			},
		})
	}

	// Match по категории
	if filters.Category != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"category": filters.Category,
			},
		})
	}

	// Фильтр по цене
	if filters.MinPrice != nil || filters.MaxPrice != nil {
		priceRange := make(map[string]interface{})
		if filters.MinPrice != nil {
			priceRange["gte"] = *filters.MinPrice
		}
		if filters.MaxPrice != nil {
			priceRange["lte"] = *filters.MaxPrice
		}

		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"price": priceRange,
			},
		})
	}

	// Только in_stock
	if filters.InStockOnly {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"availability": "in_stock",
			},
		})
	}

	// Собираем query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   must,
				"filter": filter,
			},
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации запроса: %w", err)
	}

	// Выполнение запроса
	res, err := repo.client.Search(
		repo.client.Search.WithContext(ctx),
		repo.client.Search.WithIndex(repo.index),
		repo.client.Search.WithBody(bytes.NewReader(body)),
		repo.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("ошибка ответа от ElasticSearch: %s", res.String())
	}

	// Парсим результат
	var result struct {
		Hits struct {
			Hits []struct {
				Source model.Advertisement `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа: %w", err)
	}

	var ads []model.Advertisement
	for _, hit := range result.Hits.Hits {
		ads = append(ads, hit.Source)
	}

	return ads, nil
}
