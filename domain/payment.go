package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Payment struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrderId     string             `json:"order_id" bson:"order_id"`
	Amount      float64            `json:"amount" bson:"amount"`
	Method      string             `json:"method" bson:"method"`
	Status      string             `json:"status" bson:"status"`
	PaymentType PaymentType        `json:"payment_type" bson:"payment_type"`
}
