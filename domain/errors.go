package domain

import "fmt"

type ErrorResponse struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"` // Campo opcional para indicar qual campo tem erro
}
type InvalidPaymentTypeError struct {
	Type string
}

func (e *InvalidPaymentTypeError) Error() string {
	return fmt.Sprintf("invalid payment type: %s", e.Type)
}
