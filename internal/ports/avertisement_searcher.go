package ports

import (
	"SearchService/internal/model"
	"context"
)

type AdvertisementSearcher interface {
	SearchAdvertisements(ctx context.Context, filters model.SearchFilters) ([]model.Advertisement, error)
}
