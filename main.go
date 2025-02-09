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
	_ "payments/docs" // Importa a documentação gerada pelo Swagger
	"payments/repository"
	"payments/routes"
	"payments/usecase"
)

func main() {

	// Carrega variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	// Recupera a senha do banco de dados do ambiente
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD não está definido")
	}

	// Construa a URI de conexão com MongoDB
	uri := fmt.Sprintf("mongodb+srv://kaiquemotta:%s@payments.4shch.mongodb.net/?retryWrites=true&w=majority&appName=payments", dbPassword)

	// Log para verificar a URL de conexão
	log.Printf("📝 URL de conexão com MongoDB: %s", uri)

	// Conecta ao MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Erro ao conectar ao MongoDB:", err)
	}

	// Testa a conexão com o MongoDB
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Erro ao testar a conexão com o MongoDB:", err)
	}

	// Log de sucesso na conexão
	log.Println("✅ Conexão com MongoDB realizada com sucesso!")

	// Acessa o banco de dados
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

	// Exibe a porta usada no log
	log.Printf("📌 Servidor iniciando na porta: %s", port)

	// Inicia o servidor na porta especificada
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
