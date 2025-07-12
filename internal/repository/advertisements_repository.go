package repository

import (
	"SearchService/internal"
	"SearchService/internal/model"
	"fmt"
)

type AdvertisementRepository struct {
	Database *internal.Database
}

func NewAdvertisementRepository(database *internal.Database) *AdvertisementRepository {
	return &AdvertisementRepository{Database: database}
}

func (repo *AdvertisementRepository) Save(advertisement *model.Advertisement) (*model.Advertisement, error) {
	query := `INSERT INTO advertisements (product_name, description, brand, category, price, currency, stock, ean, color, size, availability) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
	`

	stmt, err := repo.Database.DB.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка подготовки запроса: %w", err)
	}

	err = stmt.Get(&advertisement.Index, advertisement)
	if err != nil {
		return nil, fmt.Errorf("ошибка вставки объявления в БД: %w", err)
	}

	return advertisement, nil
}

func (repo *AdvertisementRepository) GetAdvertisementsBatch(limit int, offset int) ([]model.Advertisement, error) {
	query := `
		SELECT id, product_name, brand, category, price, stock, availability
		FROM advertisements 
		ORDER BY id 
		LIMIT $1 OFFSET $2
	`

	var advertisements []model.Advertisement
	err := repo.Database.DB.Select(&advertisements, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки объявлений из БД: %w", err)
	}

	return advertisements, nil
}

func (repo *AdvertisementRepository) GetAdvertisementById(id int) (model.Advertisement, error) {
	query := `
		SELECT id, product_name, brand, category, price, stock, availability
		FROM advertisements 
		WHERE id = $1
		`

	var advertisement model.Advertisement
	err := repo.Database.DB.Get(&advertisement, query, id)
	if err != nil {
		return model.Advertisement{}, fmt.Errorf("ошибка получения объявления из БД: %w", err)
	}

	return advertisement, nil
}
