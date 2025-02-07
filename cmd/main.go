package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"payments/delivery"
	"payments/repository"
	"payments/usecase"
)

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("paymentDB")

	repo := repository.NewPaymentRepository(db)
	useCase := usecase.NewPaymentUseCase(repo)

	app := fiber.New()
	delivery.NewPaymentHandler(app, useCase)

	log.Fatal(app.Listen(":8080"))
}
