package repository

import (
	"context"
	"payments/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository interface {
	GetAll() ([]domain.Payment, error)
	GetByID(id string) (domain.Payment, error)
	Create(payment *domain.Payment) error
	Update(id string, payment *domain.Payment) error
	Delete(id string) error
}

type paymentRepository struct {
	collection *mongo.Collection
}

func NewPaymentRepository(db *mongo.Database) PaymentRepository {
	return &paymentRepository{
		collection: db.Collection("payments"),
	}
}

func (r *paymentRepository) GetAll() ([]domain.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var payments []domain.Payment
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var payment domain.Payment
		if err := cursor.Decode(&payment); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *paymentRepository) GetByID(id string) (domain.Payment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Payment{}, err
	}

	var payment domain.Payment
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&payment)
	return payment, err
}

func (r *paymentRepository) Create(payment *domain.Payment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payment.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, payment)
	return err
}

func (r *paymentRepository) Update(id string, payment *domain.Payment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"amount": payment.Amount,
			"method": payment.Method,
			"status": payment.Status,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *paymentRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
