package util

// indexer содержит логику индексации и bulk-индексации объявлений (advertisement) в Elasticsearch

import (
	"SearchService/internal/model"
	"SearchService/internal/ports"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

func indexAdvertisement(esClient *elasticsearch.Client, advertisement *model.Advertisement) error {
	body := map[string]interface{}{
		"id":           advertisement.Index,
		"product_name": advertisement.Name,
		"description":  advertisement.Description,
		"brand":        advertisement.Brand,
		"category":     advertisement.Category,
		"price":        advertisement.Price,
		"currency":     advertisement.Currency,
		"stock":        advertisement.Stock,
		"ean":          advertisement.Ean,
		"color":        advertisement.Color,
		"size":         advertisement.Size,
		"availability": advertisement.Availability,
	}

	jsonBody, _ := json.Marshal(body)

	_, err := esClient.Index(
		"advertisements", // имя индекса
		bytes.NewReader(jsonBody),
		esClient.Index.WithDocumentID(fmt.Sprint(advertisement.Index)),
		esClient.Index.WithRefresh("true"), // чтобы сразу видеть в поиске
	)
	if err != nil {
		return fmt.Errorf("ошибка индексирования: %w", err)
	}

	return nil
}

func bulkIndexAdvertisements(esClient *elasticsearch.Client, advertisements []model.Advertisement) error {
	var buffer bytes.Buffer

	for _, advertisement := range advertisements {
		// Заголовок операции bulk
		meta := map[string]map[string]string{
			"index": {
				"_index": "advertisements",
				"_id":    fmt.Sprint(advertisement.Index),
			},
		}
		metaJson, _ := json.Marshal(meta)
		buffer.Write(metaJson)
		buffer.WriteByte('\n')

		document := map[string]interface{}{
			"id":           advertisement.Index,
			"product_name": advertisement.Name,
			"description":  advertisement.Description,
			"brand":        advertisement.Brand,
			"category":     advertisement.Category,
			"price":        advertisement.Price,
			"currency":     advertisement.Currency,
			"stock":        advertisement.Stock,
			"ean":          advertisement.Ean,
			"color":        advertisement.Color,
			"size":         advertisement.Size,
			"availability": advertisement.Availability,
		}
		documentJson, _ := json.Marshal(document)
		buffer.Write(documentJson)
		buffer.WriteByte('\n')
	}

	response, err := esClient.Bulk(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return fmt.Errorf("ошибка bulk вставки: %w", err)
	}
	defer response.Body.Close()

	if response.IsError() {
		return fmt.Errorf("ошибка Bulk API: %v", response.String())
	}

	log.Println("Все объявления успешно проиндексированы")
	return nil
}

func MigrationAllAdvertisements(esClient *elasticsearch.Client, loader ports.AdvertisementBatchLoader) error {
	limit := 1000
	offset := 0

	for {
		advertisements, err := loader.GetAdvertisementsBatch(limit, offset)
		if err != nil {
			return fmt.Errorf("ошибка получения объявлений: %w", err)
		}
		if len(advertisements) == 0 {
			break
		}

		log.Printf("Загружено %d объявлений", len(advertisements))

		err = bulkIndexAdvertisements(esClient, advertisements)
		if err != nil {
			return fmt.Errorf("ошибка вставки: %v", err)
		}

		if len(advertisements) < limit {
			break
		}

		offset += limit
	}

	return nil
}

func MigrationAdvertisement(esClient *elasticsearch.Client, loader ports.AdvertisementBatchLoader, id int) error {
	advertisement, err := loader.GetAdvertisementById(id)
	if err != nil {
		return fmt.Errorf("ошибка получения объявления: %v", err)
	}

	err = indexAdvertisement(esClient, &advertisement)
	if err != nil {
		return fmt.Errorf("ошибка вставки: %v", err)
	}

	return nil
}
