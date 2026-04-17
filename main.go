package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/redis/go-redis/v9"

	"go-day3/handlers"
	"go-day3/jobs"
	"go-day3/repository"
	"go-day3/service"
)

func main() {
	db := initDB()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	jobQueue := jobs.NewQueue(100)

	var jobWG sync.WaitGroup
	jobs.StartWorkerPool(3, jobQueue, &jobWG)

	itemRepo := repository.NewPostgresItemRepository(db)
	itemService := service.NewItemService(itemRepo, rdb, jobQueue, &jobWG)
	itemHandler := handlers.NewItemHandler(itemService)

	http.HandleFunc("/items", itemHandler.HandleItems)
	http.HandleFunc("/items/", itemHandler.HandleItemByID)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
