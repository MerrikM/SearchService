package config

import (
	"SearchService/internal"
	"SearchService/internal/handler"
	"SearchService/internal/util"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	DbDriverName           string
	DbConnectionString     string
	ServerAddress          string
	ElasticsearchAddresses string
	ElasticsearchUsername  string
	ElasticsearchPassword  string
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	envPath := filepath.Join(wd, "..", "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf(".env не найден по пути: %s", envPath)
	}

	DbDriverName = os.Getenv("DATABASE_DRIVER")
	DbConnectionString = os.Getenv("DATABASE_CONNECTION_URL")
	ElasticsearchAddresses = os.Getenv("ELASTICSEARCH_ADDRESSES")
	ElasticsearchUsername = os.Getenv("ELASTICSEARCH_USERNAME")
	ElasticsearchPassword = os.Getenv("ELASTICSEARCH_PASSWORD")
	DbDriverName = os.Getenv("DATABASE_DRIVER")
	DbConnectionString = os.Getenv("DATABASE_CONNECTION_URL")
	ServerAddress = os.Getenv("SERVER_ADDRESS")
}

type ElasticSearchConfig struct {
	Addresses []string
	Username  string
	Password  string
}

func newESClient(cfg ElasticSearchConfig) (*elasticsearch.Client, error) {
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

func SetupDatabase() *internal.Database {
	database, err := internal.NewDatabaseConnection(DbDriverName, DbConnectionString)
	if err != nil {
		log.Println(err)
	}

	return database
}

func SetupElasticSearch() *elasticsearch.Client {
	cfg := ElasticSearchConfig{
		Addresses: []string{ElasticsearchAddresses},
		Username:  ElasticsearchUsername,
		Password:  ElasticsearchPassword,
	}

	esClient, err := newESClient(cfg)
	if err != nil {
		log.Fatalf("ошибка инициализации Elasticsearch: %s", err)
	}
	log.Println("ElasticSearch запущен на: " + ElasticsearchAddresses)

	return esClient
}

func SetupServer(database *internal.Database) *http.Server {
	router := chi.NewRouter()
	router.Route("/search", func(r chi.Router) {
		// Контроллер для миграции данных из csv файла в БД
		r.Post("/fill_from_csv", handler.NewDatabaseFillingHandler(
			util.NewDatabaseFilling(database)).FillDatabaseAsync)
	})

	return &http.Server{
		Addr:    ServerAddress,
		Handler: router,
	}
}
