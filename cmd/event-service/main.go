package main

import (
	"event-service/config"
	"event-service/internal/db"
	"event-service/internal/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	//Читаем переменные окружения
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	// Подключение к MongoDB
	client, err := db.ConnectMongo(cfg.Mongo.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to mongo", err)
	}
	defer func() {
		if err := client.Disconnect(); err != nil {
			log.Fatal("Failed to disconnect from mongo", err)
		}
	}()

	// Получение коллекции
	database := client.GetDatabase(cfg.Mongo.DatabaseName)
	collection := database.Collection(cfg.Mongo.CollectionName)

	// Создание роутера
	r := mux.NewRouter()

	// Регистрирация обработчика
	r.HandleFunc("/api/events", handlers.GetFilterEventsHandler(collection)).Methods("POST")

	// Запуск сервера
	log.Printf("Server is running on port %s", cfg.ServerHostPort)
	if err := http.ListenAndServe(cfg.ServerHostPort, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
