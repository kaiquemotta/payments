package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"payments/config"
	_ "payments/docs"
	"payments/repository"
	"payments/routes"
	"payments/usecase"
)

func main() {
	config.InitDB()
	repo := repository.NewPaymentRepository(config.MongoDB)
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
