package main

import (
	"database/sql"
	"log"
	"order-service/internal/repository"
	delivery "order-service/internal/transport/http"
	"order-service/internal/transport/http/client"
	"order-service/internal/usecase"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/order_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresOrderRepository(db)

	paymentURL := "http://localhost:8081"
	payClient := client.NewPaymentServiceClient(paymentURL)

	orderUC := usecase.NewOrderUseCase(repo, payClient)

	orderHandler := delivery.NewOrderHandler(orderUC)
	router := delivery.SetupRouter(orderHandler)

	log.Println("Order Service started on :8082")
	log.Fatal(router.Run(":8082"))
}
