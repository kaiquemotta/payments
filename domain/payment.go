package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Payment struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Amount float64            `json:"amount" bson:"amount"`
	Method string             `json:"method" bson:"method"`
	Status string             `json:"status" bson:"status"`
}
