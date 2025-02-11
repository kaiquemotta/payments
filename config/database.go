package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB é a instância do banco de dados MongoDB acessível globalmente
var MongoDB *mongo.Database

func InitDB() {
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
	MongoDB = client.Database("paymentDB")
}
