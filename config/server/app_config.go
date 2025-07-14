package server

import (
	elasticsearch2 "SearchService/config/elasticsearch"
	"SearchService/internal"
	protobuf "SearchService/proto/your/module/path/proto"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"google.golang.org/grpc"
	"log"
	"net"
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
	GrpcNetwork            string
	GrpcAddress            string
)

type gRPCServer struct {
	protobuf.UnimplementedHelloServiceServer
}

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
	GrpcNetwork = os.Getenv("GRPC_NETWORK")
	GrpcAddress = os.Getenv("GRPC_ADDRESS")
}

func SetupDatabase() *internal.Database {
	database, err := internal.NewDatabaseConnection(DbDriverName, DbConnectionString)
	if err != nil {
		log.Println(err)
	}

	return database
}

func SetupElasticSearch() *elasticsearch.Client {
	cfg := elasticsearch2.ElasticSearchConfig{
		Addresses: []string{ElasticsearchAddresses},
		Username:  ElasticsearchUsername,
		Password:  ElasticsearchPassword,
	}

	esClient, err := elasticsearch2.NewESClient(cfg)
	if err != nil {
		log.Fatalf("ошибка инициализации Elasticsearch: %s", err)
	}
	log.Println("ElasticSearch запущен на: " + ElasticsearchAddresses)

	return esClient
}

func SetupRestServer() (*http.Server, *chi.Mux) {
	router := chi.NewRouter()

	return &http.Server{
		Addr:    ServerAddress,
		Handler: router,
	}, router
}

func RunGRPCServer() {
	listener, err := net.Listen(GrpcNetwork, GrpcAddress)
	if err != nil {
		log.Fatalf("ошибка слушателя: %v", err)
	}

	server := grpc.NewServer()
	protobuf.RegisterHelloServiceServer(server, &gRPCServer{})
}
