package util

import (
	"SearchService/internal"
	"SearchService/internal/model"
	"encoding/csv"
	"errors"
	"fmt"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type DatabaseFilling struct {
	*internal.Database
}

type BatchTask struct {
	Advertisements []model.Advertisement
}

func NewDatabaseFilling(database *internal.Database) *DatabaseFilling {
	return &DatabaseFilling{Database: database}
}

func ReadFileCSV(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("файл не был найден: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("не удалось прочитать заголовок: %w", err)
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("ошибка чтения файла: %w", err)
		}
		fmt.Println(record)
	}

	return nil
}

// FillDatabaseFromCSVSync читает CSV-файл с данными об объявлениях и записывает их в базу данных пакетами заданного размера.
//
// Параметры:
// - filepath: путь к CSV-файлу с объявлениями.
// - batchSize: количество записей в одном пакете для вставки в базу данных.
//
// Возвращает ошибку, если возникает проблема при чтении файла, парсинге данных или сохранении в базу.

func (dbf *DatabaseFilling) FillDatabaseFromCSVSync(filepath string, batchSize int) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("файл не был найден: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("не удалось прочитать заголовок: %w", err)
	}

	var advertisements []model.Advertisement
	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("ошибка чтения CSV файла: %w", err)
		}
		advertisement, err := parseCSVRecord(record)
		if err != nil {
			return fmt.Errorf("ошибка парсинга данных: %w", err)
		}
		advertisements = append(advertisements, advertisement)

		if len(advertisements) >= batchSize {
			if err := saveBatchToDatabase(dbf.Database, advertisements); err != nil {
				return fmt.Errorf("ошибка сохранения данных в БД: %w", err)
			}
			advertisements = advertisements[:0]
		}
	}

	if len(advertisements) > 0 {
		if err := saveBatchToDatabase(dbf.Database, advertisements); err != nil {
			return fmt.Errorf("ошибка сохранения оставшихся данных в БД: %w", err)
		}
	}

	return nil
}

// FillDatabaseFromCSV читает CSV-файл с данными об объявлениях и записывает их в базу данных пакетами заданного размера.
// Для повышения производительности используется пул воркеров (горутин), каждая из которых асинхронно обрабатывает свой пакет.
//
// Параметры:
// - filepath: путь к CSV-файлу с объявлениями.
// - batchSize: количество записей в одном пакете для вставки в базу данных.
//
// Особенности:
// - CSV-файл читается последовательно, данные группируются в пакеты и отправляются в канал задач.
// - Несколько воркеров параллельно извлекают задачи из канала и сохраняют данные в базу.
// - После завершения чтения все горутины завершаются корректно.
//
// Возвращает ошибку, если возникает проблема при чтении файла, парсинге данных или сохранении в базу.

func (dbf *DatabaseFilling) FillDatabaseFromCSVAsync(filepath string, batchSize int) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("файл не был найден: %w", err)
	}
	defer file.Close()

	tasks := make(chan BatchTask, 10)
	waitGroup := &sync.WaitGroup{}
	numWorkers := 4

	for i := 0; i < numWorkers; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for task := range tasks {
				if err := saveBatchToDatabase(dbf.Database, task.Advertisements); err != nil {
					log.Printf("ошибка сохранения данных в БД: %v\n", err)
				}
			}
		}()
	}

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("не удалось прочитать заголовок: %w", err)
	}

	var advertisements []model.Advertisement
	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("ошибка чтения CSV файла: %w", err)
		}
		advertisement, err := parseCSVRecord(record)
		if err != nil {
			return fmt.Errorf("ошибка парсинга данных: %w", err)
		}
		advertisements = append(advertisements, advertisement)

		if len(advertisements) >= batchSize {
			tasks <- BatchTask{Advertisements: advertisements}
			advertisements = advertisements[:0]
		}
	}

	if len(advertisements) > 0 {
		tasks <- BatchTask{Advertisements: advertisements}
	}

	close(tasks)
	waitGroup.Wait()

	return nil
}

func saveBatchToDatabase(database *internal.Database, batch []model.Advertisement) error {
	var placeholders []string
	var args []interface{}

	for index, advertisement := range batch {
		offset := index * 11

		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7, offset+8, offset+9, offset+10, offset+11))

		args = append(args,
			advertisement.Name,
			advertisement.Description,
			advertisement.Brand,
			advertisement.Category,
			advertisement.Price,
			advertisement.Currency,
			advertisement.Stock,
			advertisement.Ean,
			advertisement.Color,
			advertisement.Size,
			advertisement.Availability)
	}
	query := fmt.Sprintf(`INSERT INTO advertisements 
        (product_name, description, brand, category, price, currency, stock, ean, color, size, availability) 
        VALUES %s`, strings.Join(placeholders, ","))

	fmt.Println(query)

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка вставки данных в БД: %v", err)
	}

	return nil
}

func parseCSVRecord(record []string) (model.Advertisement, error) {
	if len(record) < 3 {
		return model.Advertisement{}, fmt.Errorf("неверный формат строки: %v", record)
	}

	index, err := strconv.Atoi(record[0])
	if err != nil {
		return model.Advertisement{}, fmt.Errorf("неверный индекс: %v", err)
	}
	price, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		return model.Advertisement{}, fmt.Errorf("неверная цена: %v", err)
	}
	stock, err := strconv.Atoi(record[7])
	if err != nil {
		return model.Advertisement{}, fmt.Errorf("неверный stock: %v", err)
	}

	return model.Advertisement{
		Index:        index,
		Name:         record[1],
		Description:  record[2],
		Brand:        record[3],
		Category:     record[4],
		Price:        price,
		Currency:     record[6],
		Stock:        stock,
		Color:        record[8],
		Size:         record[9],
		Availability: record[10],
	}, nil
}
