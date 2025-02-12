package domain

type PaymentCallback struct {
	PaymentID string  `json:"id"`
	Status    string  `json:"status"` // Pode ser 'success', 'failed', etc.
	Amount    float64 `json:"amount"`
	Message   string  `json:"message"`
	OrderId   string  `json:"order_id"`
}
