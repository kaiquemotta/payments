package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"payments/domain"
)

type PaymentRepository interface {
	GetAll() ([]domain.Payment, error)
	GetByID(id string) (domain.Payment, error)
	Create(payment *domain.Payment) (string, error) // Atualizando para incluir a assinatura correta
	Update(id string, payment *domain.Payment) error
	Delete(id string) error
}

type paymentRepository struct {
	db *mongo.Database
}

func NewPaymentRepository(db *mongo.Database) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) GetAll() ([]domain.Payment, error) {
	var payments []domain.Payment
	cursor, err := r.db.Collection("payments").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var payment domain.Payment
		if err := cursor.Decode(&payment); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *paymentRepository) GetByID(id string) (domain.Payment, error) {
	var payment domain.Payment

	// Convertendo a string para ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return payment, fmt.Errorf("invalid ObjectID format: %v", err)
	}
	err = r.db.Collection("payments").FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&payment)
	return payment, err
}

func (r *paymentRepository) Create(payment *domain.Payment) (string, error) {
	result, err := r.db.Collection("payments").InsertOne(context.Background(), payment)
	if err != nil {
		return "", err
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert InsertedID to ObjectID")
	}
	return id.Hex(), nil // Retorna o UUID gerado
}

func (r *paymentRepository) Update(id string, payment *domain.Payment) error {
	_, err := r.db.Collection("payments").UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": payment})
	return err
}

func (r *paymentRepository) Delete(id string) error {
	_, err := r.db.Collection("payments").DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
