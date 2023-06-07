package models

type Status uint

const (
	AwaitingPaymentConfirmation Status = 0
	OrderCanceled               Status = 1
	OrderCreated                Status = 2
	PaymentCanceled             Status = 3
	PaymentConfirmed            Status = 4
	PaymentRejected             Status = 5
	SentForPaymentConfirmation  Status = 6
)

func (status Status) String() string {
	switch status {
	case AwaitingPaymentConfirmation:
		return "Awaiting payment confirmation"
	case OrderCanceled:
		return "Order canceled"
	case OrderCreated:
		return "Order created"
	case PaymentCanceled:
		return "Payment canceled"
	case PaymentConfirmed:
		return "Payment confirmed"
	case PaymentRejected:
		return "Payment rejected"
	case SentForPaymentConfirmation:
		return "Sent for payment confirmation"
	}

	return "unknown"
}
