package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

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

	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD não está definido")
	}

	uri := fmt.Sprintf("mongodb+srv://kaiquemotta:%s@payments.4shch.mongodb.net/?retryWrites=true&w=majority&appName=payments", dbPassword)

	// Conecta ao MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
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
	log.Println("Registrando rotas de pagamento...")
	routes.RegisterPaymentRoutes(app, useCase)

	// Acessa a variável de ambiente PORT (Heroku define isso automaticamente)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Valor padrão para execução local
	}

	// Inicia o servidor na porta fornecida
	log.Printf("Servidor iniciado na porta %s...", port)
	log.Fatal(app.Listen(":" + port))
}
