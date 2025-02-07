package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "payments/docs" // Importa a documentação gerada pelo Swagger
	"payments/repository"
	"payments/routes"
	"payments/usecase"
)

func main() {
	// Conecta ao MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("paymentDB")

	// Cria o repositório e o caso de uso
	repo := repository.NewPaymentRepository(db)
	useCase := usecase.NewPaymentUseCase(repo)

	// Cria o aplicativo Fiber
	app := fiber.New()

	// Configura o Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Registra as rotas de pagamento
	routes.RegisterPaymentRoutes(app, useCase)

	// Inicia o servidor
	log.Fatal(app.Listen(":8080"))
}
