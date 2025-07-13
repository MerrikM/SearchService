package main

import (
	"SearchService/config"
	"SearchService/internal/handler"
	"SearchService/internal/repository"
	"context"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database := config.SetupDatabase()
	defer database.Close()

	httpServer, router := config.SetupServer(database)

	esClient := config.SetupElasticSearch()
	repo := repository.NewElasticRepository(esClient, "advertisements")

	searchHandler := handler.NewSearchHandler(repo)
	router.Get("/search", searchHandler.SearchInElastic)

	runServer(ctx, httpServer)
}

func runServer(ctx context.Context, server *http.Server) {
	serverErrors := make(chan error, 1)
	go func() {
		log.Println("Сервер запущен на " + config.ServerAddress)
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
		log.Printf("Получен сигнал %v, завершаем работу...", sig)
	}

	// Контекст с таймаутом для graceful shutdown (5 секунд)
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("ошибка при остановке сервера: %v", err)
	} else {
		log.Println("Сервер успешно остановлен")
	}
}
