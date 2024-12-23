package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

// Параметры фильтрации и пагинации
type EventFilter struct {
	FromDate string   `json:"fromDate"`
	ToDate   string   `json:"toDate"`
	MinPrice *float64 `json:"minPrice"`
	MaxPrice *float64 `json:"maxPrice"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
}

// GetFilterEventsHandler возвращает обработчик, который работает с переданной коллекцией MongoDB
func GetFilterEventsHandler(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Парсим JSON из payload
		var filter EventFilter
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			log.Printf("Error decoding payload: %v", err)
			return
		}
		defer r.Body.Close()

		log.Printf("Received filter: %+v", filter)

		// Формируем фильтр для MongoDB
		query := bson.M{}
		dateFilter := bson.M{}
		priceFilter := bson.M{}

		// Фильтрация по дате
		if filter.FromDate != "" {
			fromDate, err := time.Parse("2006-01-02", filter.FromDate)
			if err != nil {
				http.Error(w, "Invalid fromDate format (expected YYYY-MM-DD)", http.StatusBadRequest)
				log.Printf("Error parsing fromDate: %v", err)
				return
			}
			dateFilter["$gte"] = fromDate
		}
		if filter.ToDate != "" {
			toDate, err := time.Parse("2006-01-02", filter.ToDate)
			if err != nil {
				http.Error(w, "Invalid toDate format (expected YYYY-MM-DD)", http.StatusBadRequest)
				return
			}
			dateFilter["$lte"] = toDate
		}
		if len(dateFilter) > 0 {
			query["schedules.start"] = dateFilter
		}

		// Фильтрация по цене
		if filter.MinPrice != nil {
			priceFilter["$gte"] = *filter.MinPrice
		}
		if filter.MaxPrice != nil {
			priceFilter["$lte"] = *filter.MaxPrice
		}
		if len(priceFilter) > 0 {
			query["schedules.min_price.price"] = priceFilter
		}

		log.Printf("Constructed query: %+v", query)

		// Настраиваем параметры поиска (пагинация)
		findOptions := options.Find()
		if filter.Limit > 0 {
			findOptions.SetLimit(int64(filter.Limit))
		}
		if filter.Offset > 0 {
			findOptions.SetSkip(int64(filter.Offset))
		}

		// Выполняем запрос
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := collection.Find(ctx, query, findOptions)
		if err != nil {
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		// Читаем результаты
		var results []bson.M
		if err := cursor.All(ctx, &results); err != nil {
			http.Error(w, fmt.Sprintf("Error reading results: %v", err), http.StatusInternalServerError)
			log.Printf("Error parsing cursor results: %v", err)
			return
		}

		log.Printf("Found %d results", len(results))

		// Возвращаем результаты
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
