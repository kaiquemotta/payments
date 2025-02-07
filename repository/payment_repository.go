package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"payments/domain"
)

// PaymentRepository interface define os métodos para interagir com o banco
type PaymentRepository interface {
	GetAll() ([]domain.Payment, error)
	GetByID(id string) (domain.Payment, error)
	Create(payment *domain.Payment) error
	Update(id string, payment *domain.Payment) error
	Delete(id string) error
}

// paymentRepository é a implementação concreta do PaymentRepository
type paymentRepository struct {
	db *mongo.Database
}

// NewPaymentRepository cria uma nova instância do PaymentRepository
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
	err := r.db.Collection("payments").FindOne(context.Background(), bson.M{"_id": id}).Decode(&payment)
	return payment, err
}

func (r *paymentRepository) Create(payment *domain.Payment) error {
	_, err := r.db.Collection("payments").InsertOne(context.Background(), payment)
	return err
}

func (r *paymentRepository) Update(id string, payment *domain.Payment) error {
	_, err := r.db.Collection("payments").UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": payment})
	return err
}

func (r *paymentRepository) Delete(id string) error {
	_, err := r.db.Collection("payments").DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
