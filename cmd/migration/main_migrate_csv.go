package main

// Этот файл запускает процесс миграции данных из CSV-файла в базу данных PostgreSQL.
// Путь до CSV-файла передается через флаг командной строки -csv.
// Данные обрабатываются пакетами (batch) для эффективной загрузки.
//
// Пример запуска:
//   go run main_migrate_csv.go -csv=/path/to/ads.csv
//
// Используемые компоненты:
//   - config.SetupDatabase() — инициализация подключения к БД
//   - util.NewDatabaseFilling — загрузка и вставка данных из CSV

import (
	"SearchService/config"
	"SearchService/internal/util"
	"flag"
	"log"
)

func main() {
	csvPath := flag.String("csv", "", "Путь до CSV-файла для миграции")
	flag.Parse()

	if *csvPath == "" {
		log.Fatal("Укажите путь до CSV-файла с помощью флага -csv")
	}

	database := config.SetupDatabase()
	defer database.Close()

	dbf := util.NewDatabaseFilling(database)
	if err := dbf.FillDatabaseFromCSVAsync(*csvPath, 100); err != nil {
		log.Fatalf("ошибка миграции: %v", err)
	}

	log.Println("Миграция завершена")
	return
}
