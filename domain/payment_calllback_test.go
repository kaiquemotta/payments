package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPaymentCallbackInitialization(t *testing.T) {
	callback := PaymentCallback{
		PaymentID: "12345",
		Status:    "success",
		Amount:    150.75,
		Message:   "Payment completed successfully",
		OrderId:   "order123",
	}

	assert.Equal(t, "12345", callback.PaymentID)
	assert.Equal(t, "success", callback.Status)
	assert.Equal(t, 150.75, callback.Amount)
	assert.Equal(t, "Payment completed successfully", callback.Message)
	assert.Equal(t, "order123", callback.OrderId)
}
