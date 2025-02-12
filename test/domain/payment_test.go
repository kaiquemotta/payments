package domain

import (
	domain2 "payments/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPaymentInitialization(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("507f191e810c19729de860ea")

	payment := domain2.Payment{
		ID:          id,
		OrderId:     "order123",
		Amount:      250.50,
		Method:      "Credit Card",
		Status:      "Pending",
		PaymentType: domain2.Pix,
	}

	assert.Equal(t, id, payment.ID)
	assert.Equal(t, "order123", payment.OrderId)
	assert.Equal(t, 250.50, payment.Amount)
	assert.Equal(t, "Credit Card", payment.Method)
	assert.Equal(t, "Pending", payment.Status)
	assert.Equal(t, domain2.Pix, payment.PaymentType)
}
