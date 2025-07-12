package main

// Этот файл отвечает за инициализацию соединений с базой данных и Elasticsearch,
// а также за выполнение полной миграции данных из PostgreSQL в Elasticsearch.
//
// Шаги, выполняемые в этом файле:
//   1. Устанавливается соединение с PostgreSQL.
//   2. Устанавливается соединение с Elasticsearch.
//   3. Загружаются все объявления пакетами из БД.
//   4. Индексируются объявления в Elasticsearch с использованием Bulk API.
//
// Используемые компоненты:
//   - config.SetupDatabase() — инициализация подключения к БД
//	 - config.SetupElasticSearch() - инициализация подключения к Elasticsearch
//   - indexer.MigrationAllAdvertisements — загрузка и вставка данных из БД в Elasticsearch
//
// Этот процесс можно использовать для первичной инициализации или переиндексации данных.

import (
	"SearchService/config"
	"SearchService/internal/indexer"
	"SearchService/internal/repository"
	"log"
)

func main() {
	database := config.SetupDatabase()
	defer database.Close()

	esClient := config.SetupElasticSearch()
	repo := repository.NewAdvertisementRepository(database)

	err := indexer.MigrationAllAdvertisements(esClient, repo)
	if err != nil {
		log.Fatalf("ошибка миграции: %v", err)
	}

	log.Println("Миграция завершена")
	return
}
