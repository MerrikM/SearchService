package main

import (
	"SearchService/internal"
	"SearchService/internal/handler"
	"SearchService/internal/util"
	"context"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	DB_DRIVER_NAME       = "postgres"
	DB_CONNECTION_STRING = "postgresql://postgres:root@localhost:5432/search_service?sslmode=disable"
	SERVER_ADDRESS       = ":8080"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := internal.NewDatabaseConnection(DB_DRIVER_NAME, DB_CONNECTION_STRING)
	if err != nil {
		log.Println(err)
	}
	defer database.Close()

	databaseFilling := util.NewDatabaseFilling(database)
	databaseFillingHandler := handler.NewDatabaseFillingHandler(databaseFilling)

	router := chi.NewRouter()
	router.Route("/fill_db", func(r chi.Router) {
		r.Post("/fill_from_csv", databaseFillingHandler.FillDatabaseAsync)
	})

	server := &http.Server{
		Addr:    SERVER_ADDRESS,
		Handler: router,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Println("сервер запущен на " + SERVER_ADDRESS)
		serverErrors <- server.ListenAndServe()
	}()

	// Канал для сигналов ОС (Ctrl+C, kill и т.п.)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// Ждём сигнал или ошибку сервера
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ошибка работы сервера: %v", err)
		}
	case sig := <-signalChannel:
		log.Printf("получен сигнал %v, завершаем работу...", sig)
	}

	// Контекст с таймаутом для graceful shutdown (5 секунд)
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("ошибка при остановке сервера: %v", err)
	} else {
		log.Println("сервер успешно остановлен")
	}
}
