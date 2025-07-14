package elasticsearch

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticSearchConfig struct {
	Addresses []string
	Username  string
	Password  string
}

func NewESClient(cfg ElasticSearchConfig) (*elasticsearch.Client, error) {
	elasticSearchConfig := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	client, err := elasticsearch.NewClient(elasticSearchConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания клиента: %w", err)
	}

	response, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения сведений о клиенте: %w", err)
	}

	defer response.Body.Close()

	return client, nil
}
