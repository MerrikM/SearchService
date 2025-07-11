package handler

import (
	"SearchService/internal/util"
	"io"
	"net/http"
	"os"
	"strconv"
)

type DatabaseFillingHandler struct {
	*util.DatabaseFilling
}

func NewDatabaseFillingHandler(filling *util.DatabaseFilling) *DatabaseFillingHandler {
	return &DatabaseFillingHandler{DatabaseFilling: filling}
}

func (handler *DatabaseFillingHandler) FillDatabaseAsync(writer http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(10 << 20)

	file, _, err := request.FormFile("file")
	if err != nil {
		http.Error(writer, `{"error": "не удалось получить файл"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	batchSizeStr := request.FormValue("batchSize")
	if batchSizeStr == "" {
		http.Error(writer, `{"error": "batchSize обязателен"}`, http.StatusBadRequest)
		return
	}

	batchSize, err := strconv.Atoi(batchSizeStr)
	if err != nil || batchSize <= 0 {
		http.Error(writer, `{"error": "batchSize должен быть числом > 0"}`, http.StatusBadRequest)
		return
	}

	// Создаём временный файл
	tempFile, err := os.CreateTemp("", "uploaded-*.csv")
	if err != nil {
		http.Error(writer, `{"error": "не удалось создать временный файл"}`, http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name()) // удаляем временный файл после использования
	defer tempFile.Close()

	// Копируем содержимое в файл
	if _, err := io.Copy(tempFile, file); err != nil {
		http.Error(writer, `{"error": "ошибка копирования файла"}`, http.StatusInternalServerError)
		return
	}

	if err := handler.FillDatabaseFromCSVAsync(tempFile.Name(), batchSize); err != nil {
		http.Error(writer, `{"error": "ошибка загрузки данных в БД"}`, http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "multipart/form-data")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(`{"статус": "успешно"}`))
}
