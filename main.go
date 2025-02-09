package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "payments/docs"
	"payments/repository"
	"payments/routes"
	"payments/usecase"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD não está definido")
	}
	uri := fmt.Sprintf("mongodb+srv://kaiquemotta:%s@payments.4shch.mongodb.net/?retryWrites=true&w=majority&appName=payments", dbPassword)
	log.Printf("📝 URL de conexão com MongoDB: %s", uri)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Erro ao conectar ao MongoDB:", err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Erro ao testar a conexão com o MongoDB:", err)
	}
	log.Println("✅ Conexão com MongoDB realizada com sucesso!")
	db := client.Database("paymentDB")

	repo := repository.NewPaymentRepository(db)
	useCase := usecase.NewPaymentUseCase(repo)

	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

	log.Println("Registrando rotas de pagamento...")
	routes.RegisterPaymentRoutes(app, useCase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("📌 Servidor iniciando na porta: %s", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
