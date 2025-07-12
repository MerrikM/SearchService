package ports

// package ports содержит интерфейсы

import "SearchService/internal/model"

type AdvertisementBatchLoader interface {
	GetAdvertisementsBatch(limit int, offset int) ([]model.Advertisement, error)
	GetAdvertisementById(id int) (model.Advertisement, error)
}
