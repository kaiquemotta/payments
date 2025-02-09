package domain

type PaymentType string

const (
	Pix    PaymentType = "PIX"
	QRCode PaymentType = "QR_CODE"
)

func (p PaymentType) IsValid() error {
	switch p {
	case Pix, QRCode:
		return nil
	default:
		return &InvalidPaymentTypeError{Type: string(p)} // Retorna um erro customizado
	}

}
